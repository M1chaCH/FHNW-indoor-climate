package main

import (
	"sensor_hub_backend/elastic"
	"sensor_hub_backend/lifecycle"
	"sensor_hub_backend/logs"
	"sensor_hub_backend/mqtt"
	"sensor_hub_backend/rest"
	"sync"
)

func main() {
	stopContext := lifecycle.Init()
	defer lifecycle.Stop()

	var wg sync.WaitGroup

	wg.Add(2)

	elastic.InitConnection()

	go func() {
		defer wg.Done()
		rest.RunGinServer()
	}()

	go func() {
		defer wg.Done()
		mqtt.RunMqttClient()
	}()

	<-stopContext.Done()
	logs.LogInfo("Shutdown signal received, stopping...")

	wg.Wait()
	logs.LogInfo("Done!")
}
