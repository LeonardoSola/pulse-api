package main

import (
	"fmt"

	"github.com/LeonardoSola/wagon-api/adapters/bootstrap"
	"github.com/LeonardoSola/wagon-api/adapters/cmd/http"
	"github.com/LeonardoSola/wagon-api/adapters/cmd/whastapp"
	"github.com/LeonardoSola/wagon-api/adapters/database/postgres"
	"github.com/LeonardoSola/wagon-api/internal/config"

	stdhttp "net/http"
)

func main() {
	err := config.Load()
	if err != nil {
		panic("❌ Config: " + err.Error())
	}
	fmt.Println("✅ Config")

	err = postgres.NewConnection()
	if err != nil {
		panic("❌ Database: " + err.Error())
	}
	fmt.Println("✅ Database")
	defer postgres.Close()

	err = whastapp.InitContainer()
	if err != nil {
		panic("❌ WhatsApp: " + err.Error())
	}
	fmt.Println("✅ WhatsApp")

	bootstrap.Init()

	router := http.Init()
	if err != nil {
		panic("❌ HTTP handlers: " + err.Error())
	}
	fmt.Println("✅ HTTP | Port" + config.Config.Server.Port)

	if err := stdhttp.ListenAndServe(config.Config.Server.Port, router); err != nil {
		panic("❌ HTTP server: " + err.Error())
	}
}
