package posts

import "time"

// Post структура для поста
type Post struct {
    ID        int       `db:"id" json:"id"`
    Title     string    `db:"title" json:"title"`
    Content   string    `db:"content" json:"content"`
    AuthorID  int       `db:"author_id" json:"author_id"`
    CreatedAt time.Time `db:"created_at" json:"created_at"`
    UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
