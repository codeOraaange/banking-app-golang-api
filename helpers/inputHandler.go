package helpers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func HandleDBContext(ctx *gin.Context) (*pgxpool.Pool, error) {
	DB, ok := ctx.MustGet("DB").(*pgxpool.Pool)
	if !ok {
		return nil, fmt.Errorf("failed to get DB from context")
	}
	return DB, nil
}

func HandleErrorResponse(ctx *gin.Context, statusCode int, message string) {
	ctx.JSON(statusCode, gin.H{"message": message})
}
