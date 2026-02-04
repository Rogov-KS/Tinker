package database

import (
	"time"
)

type User struct {
	ID        string    `gorm:"type:text;column:id;primaryKey" json:"id"`
	Username  string    `gorm:"type:text;column:username;uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"type:text;column:password;not null" json:"-"` // Исключаем из JSON
	Avatar    *string   `gorm:"type:text;column:avatar" json:"avatar,omitempty"`
	Bio       *string   `gorm:"type:text;column:bio" json:"bio,omitempty"`
	CreatedAt time.Time `gorm:"column:createdAt;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt;autoUpdateTime" json:"updatedAt"`

	// Relations
	Posts     []Post    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Comments  []Comment `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Likes     []Like    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Following []Follow  `gorm:"foreignKey:FollowerID;constraint:OnDelete:CASCADE" json:"-"`
	Followers []Follow  `gorm:"foreignKey:FollowingID;constraint:OnDelete:CASCADE" json:"-"`
}

type Post struct {
	ID        string    `gorm:"type:text;column:id;primaryKey" json:"id"`
	UserID    string    `gorm:"type:text;column:userId;not null;index" json:"userId"`
	Content   string    `gorm:"type:text;column:content;not null" json:"content"`
	CreatedAt time.Time `gorm:"column:createdAt;autoCreateTime;index" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt;autoUpdateTime" json:"updatedAt"`

	// Relations
	User     User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Comments []Comment `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"-"`
	Likes    []Like    `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"-"`
}

type Comment struct {
	ID        string    `gorm:"type:text;column:id;primaryKey" json:"id"`
	UserID    string    `gorm:"type:text;column:userId;not null;index" json:"userId"`
	PostID    string    `gorm:"type:text;column:postId;not null;index" json:"postId"`
	Content   string    `gorm:"type:text;column:content;not null" json:"content"`
	CreatedAt time.Time `gorm:"column:createdAt;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt;autoUpdateTime" json:"updatedAt"`

	// Relations
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Post Post `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"-"`
}

type Like struct {
	ID        string    `gorm:"type:text;column:id;primaryKey" json:"id"`
	UserID    string    `gorm:"type:text;column:userId;not null;index" json:"userId"`
	PostID    string    `gorm:"type:text;column:postId;not null;index" json:"postId"`
	CreatedAt time.Time `gorm:"column:createdAt;autoCreateTime" json:"createdAt"`

	// Relations
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Post Post `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"-"`
}

type Follow struct {
	ID          string    `gorm:"type:text;column:id;primaryKey" json:"id"`
	FollowerID  string    `gorm:"type:text;column:followerId;not null;index" json:"followerId"`
	FollowingID string    `gorm:"type:text;column:followingId;not null;index" json:"followingId"`
	CreatedAt   time.Time `gorm:"column:createdAt;autoCreateTime" json:"createdAt"`

	// Relations
	Follower  User `gorm:"foreignKey:FollowerID;constraint:OnDelete:CASCADE" json:"-"`
	Following User `gorm:"foreignKey:FollowingID;constraint:OnDelete:CASCADE" json:"-"`
}

// TableName для соответствия именам таблиц из Prisma (с заглавной буквы)
func (User) TableName() string {
	return "User"
}

func (Post) TableName() string {
	return "Post"
}

func (Comment) TableName() string {
	return "Comment"
}

func (Like) TableName() string {
	return "Like"
}

func (Follow) TableName() string {
	return "Follow"
}
