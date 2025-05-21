package handler

import (
	"context"
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

		c.JSON(http.StatusCreated, resp)
	}
}

func AuthenticateUser(client userpb.UserServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req userpb.AuthRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := client.AuthenticateUser(context.Background(), &req)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, resp)
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
