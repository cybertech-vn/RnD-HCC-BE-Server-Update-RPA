package server

import "github.com/gofiber/fiber/v2"

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Hàm để trả về response dưới dạng JSON
func ResponseReturn(ctx *fiber.Ctx, success bool, message string, data interface{}, statusCode int) error {
	response := Response{
		Success: success,
		Message: message,
		Data:    data,
	}

	return ctx.Status(statusCode).JSON(response)
}
