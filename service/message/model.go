package main

import "gorm.io/gorm"

type Message struct {
	ID         int64 `gorm:"primarykey"`
	CreatedAt  int64
	UpdatedAt  int64
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	FromUserId int64
	ToUserId   int64
	Content    string
}
