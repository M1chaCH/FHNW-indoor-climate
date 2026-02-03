package rest

import (
	"context"
	"errors"
	"net/http"
	"sensor_hub_backend/lifecycle"
	"sensor_hub_backend/logs"
	"sensor_hub_backend/rest/templates"
	"time"

	"github.com/gin-gonic/gin"
)

func RunGinServer() {
	router := gin.Default()

	templates.InitTemplates(router)
	RegisterIndexRoutes(router)

	devicesGroup := router.Group("/api/v1/devices")
	RegisterDevicesRoutes(devicesGroup)

	deviceConfigGroup := router.Group("/api/v1/devices/config")
	RegisterDeviceConfigRoutes(deviceConfigGroup)

	sensorGroup := router.Group("/api/v1/sensor")
	RegisterSensorRoutes(sensorGroup)

	addr := "0.0.0.0:8080"
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		logs.LogInfo("Starting gin server on %s\n", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logs.LogErr("Error starting gin server", err)
		}
	}()

	<-lifecycle.GetStopContext().Done()
	logs.LogInfo("Stopping gin server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logs.LogErr("Error shutting down gin server", err)
	}

	logs.LogInfo("Gin server stopped")
}
