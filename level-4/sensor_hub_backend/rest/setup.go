package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sensor_hub_backend/lifecycle"
	"time"

	"github.com/gin-gonic/gin"
)

func RunGinServer() {
	router := gin.Default()
	router.GET("/api/v1/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	devicesGroup := router.Group("/api/v1/devices")
	RegisterAuthorizedDevicesRoutes(devicesGroup)

	addr := "0.0.0.0:8080"
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		fmt.Printf("Starting gin server on %s\n", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Error starting gin server: %s\n", err)
		}
	}()

	<-lifecycle.GetStopContext().Done()
	fmt.Println("Stopping gin server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("Error shutting down gin server: %s\n", err)
	}

	fmt.Println("Gin server stopped")
}
