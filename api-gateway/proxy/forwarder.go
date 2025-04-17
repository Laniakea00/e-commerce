package proxy

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Forward(c *gin.Context, target string) {
	client := &http.Client{}
	req, err := http.NewRequest(c.Request.Method, target, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	for k, v := range c.Request.Header {
		req.Header[k] = v
	}

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to reach service"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}
