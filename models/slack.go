package models

import (
	"github.com/jinzhu/gorm"
)

type SlackChannel struct {
	gorm.Model
	SlackActionID uint
	Name          string
}
