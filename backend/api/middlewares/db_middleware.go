package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func DBMiddleware(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}
