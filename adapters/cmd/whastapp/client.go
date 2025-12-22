package whastapp

import (
	"database/sql"
	"fmt"
	"log"

	"go.mau.fi/whatsmeow"
	waLog "go.mau.fi/whatsmeow/util/log"
)

var client *whatsmeow.Client

func Init(sqlInstance *sql.DB) error {
	logger := waLog.Stdout("Client", "ERROR", true)
	deviceStore, err := getStorage(sqlInstance)
	if err != nil {
		return fmt.Errorf("failed to get device store: %w", err)
	}

	client = whatsmeow.NewClient(deviceStore, logger)

	if client.Store.ID == nil {
		log.Println("🔐 No existing session found, starting QR login...")
		err = Connect()
	} else {
		log.Println("🔗 Attempting to connect with existing session...")
		err = client.Connect()
	}

	if err != nil {
		log.Printf("❌ Failed to connect: %v", err)
		return err
	}

	log.Println("✅ WhatsApp client initialized successfully")
	return nil
}

func GetClient() *whatsmeow.Client {
	return client
}
