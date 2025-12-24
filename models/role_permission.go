package models

type RolePermission struct {
	RoleID       string    `gorm:"primaryKey;size:20" json:"roleId"`
	PermissionID uint      `gorm:"primaryKey" json:"permissionId"`
}
