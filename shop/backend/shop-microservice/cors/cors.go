// cors/cors.go
package cors

import "net/http"

func WithCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Разрешаем любой origin (в продакшене поставь свой домен!)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// ← ВОТ ГЛАВНОЕ! Разрешаем Authorization и Content-Type
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Для preflight (OPTIONS) запросов
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
