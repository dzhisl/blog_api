package models

import (
	"time"

	"example.com/m/internal/types"
)

type User struct {
	ID           int       `db:"id" json:"id"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Username     string    `db:"username" json:"username"`
	FirstName    *string   `db:"first_name" json:"first_name,omitempty"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UserRole     string    `db:"user_role" json:"user_role,omitempty"`
	UserStatus   string    `db:"user_status" json:"user_status,omitempty"`
}

type Blog struct {
	ID         int       `db:"id"`
	UserID     int       `db:"user_id"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	Title      string    `db:"title"`
	Content    string    `db:"content"`
	Tags       []string  `db:"tags"`
	BlogStatus string    `db:"blog_status"`
}

type Comment struct {
	ID                 int       `db:"id"`
	BlogID             int       `db:"blog_id"`
	UserID             int       `db:"user_id"`
	CreatedAt          time.Time `db:"created_at"`
	Content            string    `db:"content"`
	RepliedToCommentID *int      `db:"replied_to_comment_id"` // nullable
	CommentStatus      string    `db:"comment_status"`
}

type Subscription struct {
	FollowerID  int       `db:"follower_id"`
	FollowingID int       `db:"following_id"`
	CreatedAt   time.Time `db:"created_at"`
}

type Reaction struct {
	UserID    int       `db:"user_id"`
	CommentID int       `db:"comment_id"`
	Reaction  string    `db:"reaction"`
	CreatedAt time.Time `db:"created_at"`
}

type JwtToken struct {
	UserID    int       `db:"user_id"`
	TokenHash string    `db:"token_hash"`
	Active    bool      `db:"active"`
	CreatedAt time.Time `db:"created_at"`
	ExpiresAt time.Time `db:"expires_at"`
}

func NewUserObject(passwordHash, username, firstname string, role types.Role, status types.Status) *User {
	return &User{
		PasswordHash: passwordHash,
		Username:     username,
		FirstName:    &firstname,
		UserRole:     string(role),
		UserStatus:   string(status),
	}
}
