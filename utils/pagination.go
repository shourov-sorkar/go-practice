package utils

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaginationParams struct {
	Page  int
	Limit int
	Skip  int
	Total int64
}

func GetPaginationParams(c *gin.Context, defaultLimit int) PaginationParams {
	limit := defaultLimit
	page := 1

	if limitParam := c.Query("limit"); limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if pageParam := c.Query("page"); pageParam != "" {
		if parsedPage, err := strconv.Atoi(pageParam); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	return PaginationParams{
		Page:  page,
		Limit: limit,
		Skip:  (page - 1) * limit,
	}
}

func GetPaginatedResponse(data interface{}, params PaginationParams) gin.H {
	return gin.H{
		"data":        data,
		"total":       params.Total,
		"page":        params.Page,
		"limit":       params.Limit,
		"total_pages": int(math.Ceil(float64(params.Total) / float64(params.Limit))),
	}
}
