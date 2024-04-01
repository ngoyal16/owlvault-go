package main

import (
	"log"

	"github.com/ngoyal16/owlvault/config"
	"github.com/ngoyal16/owlvault/encrypt"
	"github.com/ngoyal16/owlvault/storage"
	"github.com/ngoyal16/owlvault/vault"
)

func main() {
	// Read configurations
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("Failed to read configuration: %v", err)
	}

	// Initialize storage based on configuration
	dbStorage, err := storage.NewStorage(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// Initialize encryptor based on configuration
	encryptor, err := encrypt.NewEncryptor(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize encryptor: %v", err)
	}

	// Initialize OwlVault with the chosen storage implementation
	owlVault := vault.NewOwlVault(dbStorage, encryptor)

	// Example usage: Store and retrieve data
	key := "example_key"
	err = owlVault.Store(key, map[string]interface{}{
		"k1": "val1",
	})
	if err != nil {
		log.Fatalf("Failed to store data: %v", err)
	}

	storedValue, err := owlVault.RetrieveVersion(key, 1)
	if err != nil {
		log.Fatalf("Failed to retrieve data: %v", err)
	}
	log.Printf("Retrieved value: %s", storedValue)

	storedValue, err = owlVault.RetrieveLatestVersion(key)
	if err != nil {
		log.Fatalf("Failed to retrieve data: %v", err)
	}
	log.Printf("Retrieved value: %s", storedValue)
}
