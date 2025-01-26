package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

const TableNameAppCode = "app_code"

// AppCode mapped from table <app_code>
type AppCode struct {
	ID          int64          `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Code        string         `gorm:"column:code;not null" json:"code"`
	Description string         `gorm:"column:description;not null" json:"description"`
	TimeZone    string         `gorm:"column:time_zone;not null" json:"time_zone"`
	CreatedAt   time.Time      `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName AppCode's table name
func (*AppCode) TableName() string {
	return TableNameAppCode
}

// Validate 校验应用代码
func (a *AppCode) Validate() error {
	if a.Code == "" {
		return errors.New("code is required")
	}
	if a.TimeZone == "" {
		return errors.New("time_zone is required")
	}
	return nil
}
