package orders

import (
	"database/sql"
	"net/http"

	"github.com/redis/go-redis/v9"
)

func GetOrderOneDate(w http.ResponseWriter, r *http.Request, db *sql.DB, rdb *redis.Client) {

}
