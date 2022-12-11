package gpt

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"encoding/json"
	"net/http"
)

const (
	url   = "https://api.openai.com/v1/completions"
	model = "text-davinci-002"
)

type Client interface {
	SendRequest(requestText string) (string, error)
}

type gptClient struct {
	gptApiKey string
}

func NewClient(apiKey string) Client {
	return &gptClient{apiKey}
}

// RequestBody is the structure of the JSON object that we will
// send as the request body in our POST request
type RequestBody struct {
	Prompt      string  `json:"prompt"`
	Model       string  `json:"model"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float32 `json:"temperature"`
}

type ResponseBody struct {
	CompletionTokens string   `json:"id"`
	Choices          []Choice `json:"choices"`
}

type Choice struct {
	Text string `json:"text"`
}

func (c *gptClient) SendRequest(requestText string) (string, error) {
	// Create a new RequestBody struct with some sample data
	body := RequestBody{
		Prompt:      requestText,
		Model:       model,
		MaxTokens:   1024,
		Temperature: 0.5,
	}

	// Marshal the RequestBody struct into a JSON object
	jsonValue, err := json.Marshal(body)
	if err != nil {
		log.Println(fmt.Sprintf("Error Marshaling JSON request: %s", err))
		return "", err
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Println(fmt.Sprintf("Error creating HTTP request: %s", err))
		return "", err

	}

	// Set the "Content-Type" header to "application/json" so that the
	// server knows how to interpret the request body
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.gptApiKey))

	// Send the HTTP request and get the response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(fmt.Sprintf("Error sending HTTP request: %s", err))
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body and log it
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(fmt.Sprintf("Error reading HTTP response body: %s", err))
		return "", err
	}
	log.Println(string(responseBody))

	response := ResponseBody{}
	if err := json.Unmarshal(responseBody, &response); err != nil {
		log.Println(fmt.Sprintf("Error Unmarshal response body: %s", err))
		return "", err
	}

	if len(response.Choices) == 0 {
		return "", nil
	} else {
		var resultText string
		for _, c := range response.Choices {
			resultText += c.Text
			resultText += "\n"
		}
		return resultText, nil
	}

}
