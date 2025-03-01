package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type ChatController struct {
	DB *gorm.DB
	upgrader websocket.Upgrader
	clients map[*websocket.Conn]bool
	mutex sync.Mutex
}

type Message struct {
	ID uint `json:"id" gorm:"primaryKey"`
	UserID uint `json:"user_id"`
	Content string `json:"content"`
	Reply string `json:"reply"`
}

func NewChatController(db *gorm.DB) *ChatController {
	return &ChatController{
		DB: db,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		clients: make(map[*websocket.Conn]bool),
	}
}

func (cc *ChatController) HandleWebSocket(c *gin.Context) {
	conn, err := cc.upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	defer conn.Close()

	cc.mutex.Lock()
	cc.clients[conn] = true
	cc.mutex.Unlock()

	userID, exists := c.Get("userId")

	if !exists {
		log.Printf("User ID not found in context")
	}

	for {
		_, message, err := conn.ReadMessage() // membaca pesan dari koneksi socket
		if err != nil {
			cc.mutex.Lock()
			delete(cc.clients, conn)
			cc.mutex.Unlock()
			break
		}

		aiResponse, err := cc.getAIResponse(string(message))

		if err != nil {
			log.Printf("error")
			continue
		}

		response := Message {
			UserID: userID.(uint),
			Content: string(message),
			Reply: aiResponse,
		}

		if err := cc.DB.Create(&response).Error; err != nil {
			log.Printf("error saving message: %v", err)
			continue
		}

		responseJSON, _ := json.Marshal(response)
		if err := conn.WriteMessage(websocket.TextMessage, responseJSON); err != nil {
			log.Printf("error sending response: %v", err)
			return
		}
	}
}

func (cc *ChatController) getAIResponse(message string) (string, error) {
	client := &http.Client{}
	
	prompt := fmt.Sprintf(`[INST] Kamu adalah asisten AI yang sangat membantu. 
Berikan jawaban yang singkat, jelas, dan dalam Bahasa Indonesia untuk pertanyaan berikut:

%s

Jawab dengan format yang mudah dibaca dan dipahami. [/INST]`, message)
	
	payload := map[string]interface{}{
		"inputs": prompt,
		"parameters": map[string]interface{}{
			"max_length": 300,
			"temperature": 0.7,
			"top_p": 0.95,
			"repetition_penalty": 1.15,
			"return_full_text": false,
		},
	}
	
	modelURL := "https://api-inference.huggingface.co/models/nvidia/Llama-3.1-Nemotron-70B-Instruct-HF"
	
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", modelURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", err
	}

	apiKey := os.Getenv("HUGGINGFACE_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("HUGGINGFACE_API_KEY not set in environment")
	}
	
	req.Header.Add("Authorization", "Bearer "+apiKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	log.Printf("Raw response: %s", string(body))


	// response := Message {
	// 	UserID: userID.(uint),
	// 	Content: string(message),
	// 	Reply: aiResponse,
	// }

	var result []map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		// cc.DB.Create(&response)
		return "", fmt.Errorf("error parsing response: %v", err)
	}

	// if err := cc.DB.Create(&response).Error; err != nil {
	// 	log.Printf("error saving message: %v", err)
	// 	continue
	// }

	if len(result) == 0 {
		return "", fmt.Errorf("empty response from API")
	}

	if generatedText, ok := result[0]["generated_text"].(string); ok {
		return generatedText, nil
	}

	return "", fmt.Errorf("unexpected response format")
} 