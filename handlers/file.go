package handlers

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"ft-backend/common/config"
	"ft-backend/database"
	"ft-backend/models"
	"ft-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UploadFile 文件上传
func UploadFile(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "未授权"})
		return
	}

	// 获取配置
	cfg := c.MustGet("config").(*config.Config)

	// 限制文件大小
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, cfg.File.MaxFileSize)

	// 获取上传文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "获取文件失败", "error": err.Error()})
		return
	}
	defer file.Close()

	// 验证文件格式
	if !utils.ValidateFileExtension(header.Filename, cfg.File.AllowedFormats) {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "文件格式不允许"})
		return
	}

	// 计算文件哈希
	fileHash, err := utils.CalculateFileHash(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "计算文件哈希失败"})
		return
	}

	// 生成唯一文件名
	uniqueFilename := utils.GenerateUniqueFilename(header.Filename)

	// 保存文件
	filePath, fileSize, err := utils.SaveUploadedFile(file, header, cfg.File.UploadDir, uniqueFilename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "保存文件失败", "error": err.Error()})
		return
	}

	// 保存文件信息到数据库
	newFile := models.File{
		UserID:       userID.(uint),
		Filename:     uniqueFilename,
		OriginalName: header.Filename,
		Size:         fileSize,
		Path:         filePath,
		MimeType:     header.Header.Get("Content-Type"),
		Extension:    utils.GetFileExtension(header.Filename),
		Hash:         fileHash,
		Status:       "available",
		Visibility:   "private",
	}

	if err := database.DB.Create(&newFile).Error; err != nil {
		// 删除已上传的文件
		utils.DeleteFile(filePath)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "保存文件元数据失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code": 201,
		"msg":  "文件上传成功",
		"data": gin.H{
			"file": newFile,
		},
	})
}

// DownloadFile 文件下载
func DownloadFile(c *gin.Context) {
	// 获取文件ID
	fileIDStr := c.Param("file_id")
	fileID, err := strconv.ParseUint(fileIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "无效的文件ID"})
		return
	}

	// 获取文件信息
	var file models.File
	if err := database.DB.Where("id = ? AND status = ?", fileID, "available").First(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "文件不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "数据库错误"})
		}
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(file.Path); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "服务器上文件不存在"})
		return
	}

	// 设置响应头
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+file.OriginalName)
	c.Header("Content-Type", file.MimeType)
	c.Header("Content-Length", strconv.FormatInt(file.Size, 10))

	// 发送文件
	c.File(file.Path)

	// 更新下载次数
	go func() {
		database.DB.Model(&file).UpdateColumn("download_count", gorm.Expr("download_count + ?", 1))
	}()
}

// ListFiles 文件列表
func ListFiles(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "未授权"})
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	offset := (page - 1) * pageSize

	// 查询文件列表
	var files []models.File
	var total int64

	db := database.DB.Model(&models.File{}).Where("user_id = ? AND deleted_at IS NULL", userID)

	// 应用筛选条件
	if status := c.Query("status"); status != "" {
		db = db.Where("status = ?", status)
	}

	if visibility := c.Query("visibility"); visibility != "" {
		db = db.Where("visibility = ?", visibility)
	}

	// 计算总数
	db.Count(&total)

	// 获取分页数据
	if err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&files).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "获取文件列表失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取文件列表成功",
		"data": gin.H{
			"list":  files,
			"total": total,
		},
	})
}

// GetFileInfo 文件信息
func GetFileInfo(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "未授权"})
		return
	}

	// 获取文件ID
	fileIDStr := c.Param("file_id")
	fileID, err := strconv.ParseUint(fileIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "无效的文件ID"})
		return
	}

	// 获取文件信息
	var file models.File
	if err := database.DB.Where("id = ? AND user_id = ? AND deleted_at IS NULL", fileID, userID).First(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "文件不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "数据库错误"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取文件信息成功",
		"data": file,
	})
}

// DeleteFile 删除文件
func DeleteFile(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "未授权"})
		return
	}

	// 获取文件ID
	fileIDStr := c.Param("file_id")
	fileID, err := strconv.ParseUint(fileIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "无效的文件ID"})
		return
	}

	// 获取文件信息
	var file models.File
	if err := database.DB.Where("id = ? AND user_id = ? AND deleted_at IS NULL", fileID, userID).First(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "文件不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "数据库错误"})
		}
		return
	}

	// 软删除文件
	if err := database.DB.Delete(&file).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "删除文件失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "文件删除成功",
	})
}

// ShareFile 分享文件
func ShareFile(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "未授权"})
		return
	}

	// 获取文件ID
	fileIDStr := c.Param("file_id")
	fileID, err := strconv.ParseUint(fileIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "无效的文件ID"})
		return
	}

	// 获取文件信息
	var file models.File
	if err := database.DB.Where("id = ? AND user_id = ? AND deleted_at IS NULL", fileID, userID).First(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "文件不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "数据库错误"})
		}
		return
	}

	// 生成分享密钥
	shareKey := utils.GenerateUniqueFilename(file.Filename)

	// 创建分享记录
	share := models.Share{
		FileID:    file.ID,
		ShareKey:  shareKey,
		ExpiresAt: database.DB.NowFunc().AddDate(0, 0, 7), // 7天过期
	}

	if err := database.DB.Create(&share).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "创建分享链接失败", "error": err.Error()})
		return
	}

	// 构建分享链接
	shareURL := strings.Join([]string{
		c.Request.Host, "/api/files/download/", shareKey,
	}, "")

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "文件分享成功",
		"data": gin.H{
			"share_key":  shareKey,
			"share_url":  shareURL,
			"expires_at": share.ExpiresAt,
		},
	})
}

// GetSharedFiles 共享文件列表
func GetSharedFiles(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "未授权"})
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	offset := (page - 1) * pageSize

	// 查询分享列表
	var shares []models.Share
	var total int64

	db := database.DB.Model(&models.Share{}).Preload("File").Joins("JOIN files ON shares.file_id = files.id").Where("files.user_id = ? AND shares.expires_at > ?", userID, database.DB.NowFunc())

	// 计算总数
	db.Count(&total)

	// 获取分页数据
	if err := db.Order("shares.created_at DESC").Offset(offset).Limit(pageSize).Find(&shares).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "获取共享文件列表失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取共享文件列表成功",
		"data": gin.H{
			"list":  shares,
			"total": total,
		},
	})
}
