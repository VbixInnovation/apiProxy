package main

import (
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/gin-contrib/cors"
	"time"
)

type RequestBody struct {
	Endpoint string                 `json:"endpoint"`
	Body     map[string]interface{} `json:"body"`
}

func apiProxy(c *gin.Context) {
	var reqBody RequestBody

	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if reqBody.Endpoint == "" || reqBody.Body == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing endpoint or body in request"})
		return
	}

	client := resty.New()
	resp, err := client.R().
		SetBody(reqBody.Body).
		Post(reqBody.Endpoint)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request", "details": err.Error()})
		return
	}

	c.Data(resp.StatusCode(), resp.Header().Get("Content-Type"), resp.Body())
}

func main() {
	router := gin.Default()

	// การตั้งค่า CORS ด้วย wildcard domain
	router.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return strings.HasSuffix(origin, ".kissflow.com")
		},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.POST("/apiProxy", apiProxy)
	router.Run(":8080")
}
