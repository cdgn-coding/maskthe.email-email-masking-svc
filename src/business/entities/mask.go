package entities

import (
	"gorm.io/gorm"
)

type EmailMask struct {
	gorm.Model
	Alias  string `json:"address" gorm:"primaryKey"`
	Target string `json:"target"`
}
