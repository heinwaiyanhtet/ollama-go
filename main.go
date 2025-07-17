package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type AskRequest struct {
	Prompt string `json:"prompt"`
}

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}


type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}


func main(){

	r := gin.Default()


	r.Use(cors.New(cors.Config{
    	AllowOrigins:     []string{"*"}, 
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))
	

	r.POST("/api/ask",func(c *gin.Context) {

		var req AskRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		ollamaReq := OllamaRequest{
			Model:  "llama3.2",
			Prompt: req.Prompt,
		}

		ollamaBody, _ := json.Marshal(ollamaReq)

		err := godotenv.Load()
		
		if err != nil {
			panic("Error loading .env file")
		}
		modelEndpoint := os.Getenv("MODEL_ENDPOINT")

		resp, err := http.Post(modelEndpoint+"/api/generate", "application/json", bytes.NewBuffer(ollamaBody));

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to contact Ollama"})
			return
		}

		defer resp.Body.Close()

		var buf bytes.Buffer

		dec := json.NewDecoder(resp.Body)
		for {
			var or OllamaResponse
			if err := dec.Decode(&or); err == io.EOF {
				break
			} else if err != nil {
				break
			}
			buf.WriteString(or.Response)
			if or.Done {
				break
			}
		}

		c.JSON(http.StatusOK, gin.H{"response": buf.String()})

	})

    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })



    r.Run(":8080")
}