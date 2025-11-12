package post

type Post struct {
	ID        int64  `db:"id" json:"id"`
	UserID    int64  `db:"user_id" json:"user_id"`
	Title     string `db:"title" json:"title"`
	Content   string `db:"content" json:"content"`
	CreatedAt string `db:"created_at" json:"created_at"`
}
