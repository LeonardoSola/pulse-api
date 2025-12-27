package main

import (
	"fmt"

	"github.com/LeonardoSola/wagon-api/adapters/bootstrap"
	"github.com/LeonardoSola/wagon-api/adapters/cmd/http"
	"github.com/LeonardoSola/wagon-api/adapters/cmd/whastapp"
	"github.com/LeonardoSola/wagon-api/adapters/database/postgres"
	"github.com/LeonardoSola/wagon-api/internal/config"
	"github.com/LeonardoSola/wagon-api/internal/logger"

	stdhttp "net/http"
)

func init() {
	err := config.Load()
	if err != nil {
		panic("❌ Config: " + err.Error())
	}

	err = logger.Init()
	if err != nil {
		panic("❌ Logger: " + err.Error())
	}
	logger.Log.Info("✅ Logger")
}

func main() {
	err := postgres.NewConnection()
	if err != nil {
		panic("❌ Database: " + err.Error())
	}
	logger.Log.Info("✅ Database")
	defer postgres.Close()

	err = whastapp.InitContainer()
	if err != nil {
		panic("❌ WhatsApp: " + err.Error())
	}
	logger.Log.Info("✅ WhatsApp")

	bootstrap.Init()

	router := http.Init()
	if err != nil {
		panic("❌ HTTP handlers: " + err.Error())
	}
	logger.Log.Info("✅ HTTP", "port", config.Config.Server.Port)

	if err := stdhttp.ListenAndServe(fmt.Sprintf(":%d", config.Config.Server.Port), router); err != nil {
		panic("❌ HTTP server: " + err.Error())
	}
}
