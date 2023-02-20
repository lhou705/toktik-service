package main

import "gorm.io/gorm"

type User struct {
	ID              int64 `gorm:"primarykey"`
	CreatedAt       int64
	UpdatedAt       int64
	DeletedAt       gorm.DeletedAt `gorm:"index"`
	Name            string         `gorm:"unique"`
	Password        string
	FollowCount     int64  `gorm:"default:0;check:follow_count >=0"`
	FollowerCount   int64  `gorm:"default:0;check:follower_count >=0"`
	Avatar          string `gorm:"default:https://cdn.lhou.ltd/avatar/avatar.jpeg"`
	BackgroundImage string `gorm:"default:https://cdn.lhou.ltd/background/background.jpeg"`
	Signature       string `gorm:"default:'这是一个签名'"`
	TotalFavorited  int64  `gorm:"default:0;check:total_favorited >=0"`
	WorkCount       int64  `gorm:"default:0;check:work_count >=0"`
	FavoriteCount   int64  `gorm:"default:0;check:favorite_count >=0"`
}

type Follow struct {
	ID         int64 `gorm:"primarykey"`
	CreatedAt  int64
	UpdatedAt  int64
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	FollowId   int64
	FollowerId int64
	IsFollow   bool `gorm:"default:false"`
	IsMutual   bool `gorm:"default:false"`
}

type Message struct {
	ID         int64 `gorm:"primarykey"`
	CreatedAt  int64
	UpdatedAt  int64
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	FromUserId int64
	ToUserId   int64
	Content    string
}
