package main

import (
	"api-status-check/internal/api"
	"api-status-check/internal/service"
	"log"
)

func main() {
	log.Println("Starting API Status Check System...")

	checkService := service.NewCheckService()
	scheduler := service.NewScheduler(checkService)
	scheduler.Start()

	handler := api.NewHandler()
	router := api.SetupRouter(handler)

	log.Println("Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
