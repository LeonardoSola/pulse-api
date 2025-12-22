package main

import "github.com/LeonardoSola/wagon-api/adapters/cmd/whastapp"

func main() {
	err := whastapp.Init()
	if err != nil {
		log.Fatalf("Failed to initialize WhatsApp client: %v", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	logger.Infof("Shutting down...")
	client.Disconnect()
}
