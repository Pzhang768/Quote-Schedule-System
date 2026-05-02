package handler

import "github.com/gin-gonic/gin"

type Response[T any] struct {
	Data T `json:"data"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func Success[T any](c *gin.Context, status int, data T) {
	c.JSON(status, Response[T]{Data: data})
}

func Fail(c *gin.Context, status int, msg string) {
	c.JSON(status, ErrorResponse{Error: msg})
}
