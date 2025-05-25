package handler

import (
	"context"
	"fmt"
	"github.com/Laniakea00/e-commerce/api-gateway/utils"
	"net/http"
	"strconv"

	userpb "github.com/Laniakea00/e-commerce/proto/user"
	"github.com/gin-gonic/gin"
)

func RegisterUser(client userpb.UserServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req userpb.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := client.RegisterUser(context.Background(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Создаем фейковый токен подтверждения (в будущем храни его в базе)
		token := "fake-verification-token" // или UUID + база
		link := fmt.Sprintf("http://localhost:8080/users/verify?email=%s&token=%s", req.Email, token)

		// Отправка письма
		go func() {
			if err := utils.SendEmailVerification(req.Email, link); err != nil {
				fmt.Println("Failed to send verification email:", err)
			}
		}()

		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "User registered. Please check your email to verify.",
			"user":    resp.User,
		})
	}
}

func AuthenticateUser(client userpb.UserServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		// gRPC вызов к user-service для проверки логина/пароля
		resp, err := client.AuthenticateUser(context.Background(), &userpb.AuthRequest{
			Email:    req.Email,
			Password: req.Password,
		})
		if err != nil || !resp.Success {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		// Генерация токена здесь, в API Gateway
		token, err := utils.GenerateJWT(int32(resp.User.Id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
			return
		}

		// Ответ клиенту
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": resp.Message,
			"user":    resp.User,
			"token":   token,
		})
	}
}

func GetUserProfile(client userpb.UserServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		resp, err := client.GetUserProfile(context.Background(), &userpb.UserID{Id: int32(id)})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}

func UpdateUserProfile(client userpb.UserServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req userpb.UpdateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := client.UpdateUserProfile(context.Background(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}

func DeleteUser(client userpb.UserServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		resp, err := client.DeleteUser(context.Background(), &userpb.UserID{Id: int32(id)})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}

func ListUsers(client userpb.UserServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := client.ListUsers(context.Background(), &userpb.Empty{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp.Users)
	}
}

func VerifyUserEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.Query("email")
		token := c.Query("token")

		if email == "" || token == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing email or token"})
			return
		}

		// Здесь могла бы быть логика проверки токена
		// А пока просто отвечаем, что email подтвержден
		c.HTML(http.StatusOK, "verified.html", gin.H{"email": email})
	}
}
