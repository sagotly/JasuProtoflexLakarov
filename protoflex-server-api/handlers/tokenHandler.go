package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"protoflex-server-api/controllers"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// TokenHandler manages the generation and validation of tokens.
type TokenHandler struct {
	mu     sync.RWMutex
	tokens map[string]storedToken
	ttl    time.Duration

	WgController *controllers.WgController
}

type storedToken struct {
	token  string
	expiry time.Time
	ip     string // Stores the IP address of the client
}

// NewTokenHandler creates a new instance of TokenHandler.
func NewTokenHandler(ttl time.Duration /*wgController *controllers.WgController*/) *TokenHandler {
	return &TokenHandler{
		tokens: make(map[string]storedToken),
		ttl:    ttl,
		// WgController: wgController,
	}
}

// generateToken creates a new random token.
func generateToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateTokenHandler generates a token and stores it for the client's IP.
func (th *TokenHandler) GenerateTokenHandler(c *gin.Context) {
	clientIP := c.ClientIP()

	// Generate a new token
	token, err := generateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Store the token with expiration and client IP
	th.mu.Lock()
	th.tokens[token] = storedToken{
		token:  token,
		expiry: time.Now().Add(th.ttl),
		ip:     clientIP,
	}
	th.mu.Unlock()

	// Print the token to the console in bold and red
	fmt.Printf("\n\n\033[1;31mGenerated token: %s for IP: %s\033[0m\n\n\n", token, clientIP)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

// ValidateTokenHandler validates the token and ensures it originates from the same IP.
func (th *TokenHandler) ValidateTokenHandler(c *gin.Context) {
	clientIP := c.ClientIP()
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
		return
	}

	th.mu.RLock()
	stored, exists := th.tokens[token]
	th.mu.RUnlock()

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	if time.Now().After(stored.expiry) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
		return
	}

	if stored.ip != clientIP {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token does not match client IP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Token is valid",
	})

	// serverIp, err := utils.GetServerEndpoint()
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get server IP"})
	// 	return
	// }
	// config, err := th.WgController.AddPeerToWGServer(serverIp)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	// c.JSON(http.StatusOK, gin.H{
	// 	"config": config,
	// })
}

// // TestAddPeerHandler validates the token and ensures it originates from the same IP.
// func (th *TokenHandler) TestAddPeerHandler(c *gin.Context) {
// 	serverIp, err := utils.GetServerEndpoint()
// 	config, err := th.WgController.AddPeerToWGServer(serverIp)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"config": config,
// 	})
// }
