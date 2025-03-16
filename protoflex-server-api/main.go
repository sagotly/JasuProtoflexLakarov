package main

import (
	"fmt"
	"protoflex-server-api/handlers"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
)

var port int
var prod string

func main() {
	pflag.IntVarP(&port, "port", "p", 8080, "Port for the server to listen on")
	pflag.Parse()
	// WgInterfaceFlag := flag.String("i", "wg0", "Wireguard interface to use")
	// AllowedPeerIpFlag := flag.String("a", "10.0.0.2/32", "Allowed peer ip for the wireguard interface")
	// wgController := controllers.NewWGController(*WgInterfaceFlag, *AllowedPeerIpFlag)

	tokenHandler := handlers.NewTokenHandler(5 * time.Minute)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	// Enable CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	// Route to generate a token
	r.GET("/generate", tokenHandler.GenerateTokenHandler)

	// Route to validate a token
	r.GET("/validate", tokenHandler.ValidateTokenHandler)
	// Route to validate a token
	// r.GET("/testAddPeer", tokenHandler.TestAddPeerHandler)

	// Start the server on the specified port
	addr := fmt.Sprintf(":%d", port)
	fmt.Print("\nServer started on port ", port, "\n\n")
	err := r.Run(addr)
	if err != nil {
		fmt.Println(err)
	}
}
