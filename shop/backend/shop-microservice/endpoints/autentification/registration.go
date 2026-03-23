package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"shop/models"

	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var creds models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	var id int
	err = db.QueryRow(`INSERT INTO users (nick, password_hash)
					  VALUES ($1, $2)
	`, creds.Nick, hashedPassword).Scan(&id)

	if err != nil {
		http.Error(w, "Nick already exists", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"user_id": id})
}
