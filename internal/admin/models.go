package admin

import "time"

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	IsBanned  bool      `json:"is_banned"`
	CreatedAt time.Time `json:"created_at"`
}

type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type Analytics struct {
	TotalPosts   int `json:"total_posts"`
	ActiveUsers  int `json:"active_users"`
	BannedUsers  int `json:"banned_users"`
	PostsPending int `json:"posts_pending"`
}
