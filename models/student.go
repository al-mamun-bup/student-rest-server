package models

import "gorm.io/gorm"

type Student struct {
	gorm.Model        // Includes ID, CreatedAt, UpdatedAt, DeletedAt
	Name       string `json:"name" gorm:"type:varchar(100);not null"`
	Age        int    `json:"age" gorm:"not null"`
	Grade      string `json:"grade" gorm:"type:varchar(20);not null"`
}

func (Student) TableName() string {
	return "students"
}
