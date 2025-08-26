package main

import (
	"context"

	"example.com/m/internal/api/router"
	"example.com/m/internal/storage"
	"example.com/m/pkg/config"
	"example.com/m/pkg/logger"
)

func InitApp() {
	config.InitConfig()
	logger.InitLogger()
	storage.InitDb(context.Background())
}

func main() {
	InitApp()
	logger.Info(context.Background(), "app is setup")
	r := router.InitRouter()
	logger.Info(context.Background(), "running API")
	r.Run(":8080")
}
