package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	defaultGroqAPIKey = os.Getenv("GROQ_API_KEY")
	proxyURL          = os.Getenv("PROXY_URL")
)

func main() {
	r := gin.Default()
	r.POST("/v1/audio/transcriptions", transcribeAudio)
	r.Run(":8000")
}

func transcribeAudio(c *gin.Context) {
	// Extract API Key from Authorization header
	openaiAPIKey := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")

	// Get file from form data
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}
	defer file.Close()

	// Get model from form data
	model := c.PostForm("model")
	if model == "" {
		model = "whisper-large-v3"
	}

	// Use openaiAPIKey if provided, otherwise use default
	groqAPIKey := openaiAPIKey
	if groqAPIKey == "" {
		groqAPIKey = defaultGroqAPIKey
	}

	// Prepare the request to Groq
	groqURL := "https://api.groq.com/openai/v1/audio/transcriptions"
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add file to the request
	part, err := writer.CreateFormFile("file", header.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_, err = io.Copy(part, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Add other form fields
	writer.WriteField("model", model)
	writer.WriteField("temperature", "0")
	writer.WriteField("language", "zh")
	writer.WriteField("response_format", "json")

	err = writer.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create the request
	req, err := http.NewRequest("POST", groqURL, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", groqAPIKey))

	// Create a client with proxy if set
	client := &http.Client{}
	if proxyURL != "" {
		proxyURLParsed, err := url.Parse(proxyURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid proxy URL"})
			return
		}
		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURLParsed)}
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"error": fmt.Sprintf("Groq API returned status code %d", resp.StatusCode)})
		return
	}

	// Read and parse the response
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse Groq API response"})
		return
	}

	// Return the Groq response
	c.JSON(http.StatusOK, result)
}
