package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"shop/constants"
	"shop/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var creds models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Bad request", http.StatusBadGateway)
		return
	}

	var storedHash string
	var userID int
	err := db.QueryRow("SELECT id, password_hash FROM users WHERE nick=$1", creds.Nick).Scan(&userID, &storedHash)

	if err == sql.ErrNoRows {
		http.Error(w, "Invalid nick or password", http.StatusUnauthorized)
		return
	}

	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(creds.Password)); err != nil {
		http.Error(w, "Invalid nick or password", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &models.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(constants.JwtKey)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func Protected(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		return constants.JwtKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	w.Write([]byte("Добро пожаловать, user_id = " + string(rune(claims.UserID))))
}

// Проверка: авторизован ли пользователь (вызывается с фронта при загрузке страницы)
func Me(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Извлекаем токен
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
		tokenStr = tokenStr[7:]
	}

	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return constants.JwtKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var user string

	err = db.QueryRow("SELECT nick FROM users WHERE id = $1", claims.UserID).Scan(&user)

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"authenticated": true,
		"user_id":       claims.UserID,
		"nick":          user,
	})
}
