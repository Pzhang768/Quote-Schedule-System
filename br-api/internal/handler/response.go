package handler

import "github.com/gin-gonic/gin"

type Response[T any] struct {
	Data T `json:"data"`
}

type PagedResponse[T any] struct {
	Data       T   `json:"data"`
	Total      int `json:"total"`
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalPages int `json:"total_pages"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func Success[T any](c *gin.Context, status int, data T) {
	c.JSON(status, Response[T]{Data: data})
}

func SuccessPaged[T any](c *gin.Context, status int, data T, total, page, pageSize int) {
	totalPages := (total + pageSize - 1) / pageSize
	c.JSON(status, PagedResponse[T]{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

func Fail(c *gin.Context, status int, msg string) {
	c.JSON(status, ErrorResponse{Error: msg})
}
