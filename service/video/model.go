package main

import "gorm.io/gorm"

type Video struct {
	ID            int64 `gorm:"primarykey"`
	CreatedAt     int64
	UpdatedAt     int64
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	AuthorId      int64
	PlayUrl       string
	CoverUrl      string
	FavoriteCount int64 `gorm:"default:0;check:favorite_count >= 0"`
	CommentCount  int64 `gorm:"default:0;check:comment_count >= 0"`
	Title         string
}

type Favorite struct {
	ID         int64 `gorm:"primarykey"`
	CreatedAt  int64
	UpdatedAt  int64
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	UserId     int64
	VideoId    int64
	IsFavorite bool `gorm:"default:false"`
}

type Comment struct {
	ID        int64 `gorm:"primarykey"`
	CreatedAt int64
	UpdatedAt int64
	DeletedAt gorm.DeletedAt `gorm:"index"`
	UserId    int64
	VideoId   int64
	Content   string
}
type User struct {
	ID              int64 `gorm:"primarykey"`
	CreatedAt       int64
	UpdatedAt       int64
	DeletedAt       gorm.DeletedAt `gorm:"index"`
	Name            string         `gorm:"unique"`
	Password        string
	FollowCount     int64 `gorm:"default:0;check:follow_count >=0"`
	FollowerCount   int64 `gorm:"default:0;check:follower_count >=0"`
	Avatar          string
	BackgroundImage string
	Signature       string `gorm:"default:'这是一个签名'"`
	TotalFavorited  int64  `gorm:"default:0;check:total_favorited >=0"`
	WorkCount       int64  `gorm:"default:0;check:work_count >=0"`
	FavoriteCount   int64  `gorm:"default:0;check:favorite_count >=0"`
}
