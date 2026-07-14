package model

type RoleServer struct {
	RoleID   uint `gorm:"primaryKey" json:"role_id"`
	ServerID uint `gorm:"primaryKey" json:"server_id"`
}
