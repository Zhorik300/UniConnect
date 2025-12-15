package admin

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func BanUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		id, _ := strconv.Atoi(idStr)

		_, err := db.Exec(`UPDATE users SET is_banned = true WHERE id = $1`, id)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		writeJSON(w, map[string]string{"status": "user banned"})
	}
}

func UnbanUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		id, _ := strconv.Atoi(idStr)

		_, err := db.Exec(`UPDATE users SET is_banned = false WHERE id = $1`, id)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		writeJSON(w, map[string]string{"status": "user unbanned"})
	}
}

func ApprovePost(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.URL.Query().Get("id"))

		_, err := db.Exec(`UPDATE posts SET status = 'approved' WHERE id = $1`, id)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		writeJSON(w, map[string]string{"status": "post approved"})
	}
}

func DeletePost(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.URL.Query().Get("id"))

		_, err := db.Exec(`UPDATE posts SET status = 'deleted' WHERE id = $1`, id)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		writeJSON(w, map[string]string{"status": "post deleted"})
	}
}

func AnalyticsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var a Analytics

		db.QueryRow(`SELECT COUNT(*) FROM posts`).Scan(&a.TotalPosts)
		db.QueryRow(`SELECT COUNT(*) FROM users WHERE is_banned = false`).Scan(&a.ActiveUsers)
		db.QueryRow(`SELECT COUNT(*) FROM users WHERE is_banned = true`).Scan(&a.BannedUsers)
		db.QueryRow(`SELECT COUNT(*) FROM posts WHERE status = 'pending'`).Scan(&a.PostsPending)

		writeJSON(w, a)
	}
}
