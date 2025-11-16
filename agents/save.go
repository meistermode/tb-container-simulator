package agents

import (
	"container-simulator/models"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"io"
	"os"
)

type SaveData struct {
	Balance        int                 `json:"b"`
	Inventory      models.Inventory    `json:"i"`
	Pity           models.PityCounters `json:"p"`
	LastRefreshDay string              `json:"r"`
}

type SaveAgent struct {
	key []byte
}

func NewSaveAgent() *SaveAgent {
	machineID := getMachineID()
	hash := sha256.Sum256([]byte(machineID + "container-sim-secret-key-v1"))
	return &SaveAgent{
		key: hash[:],
	}
}

func getMachineID() string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "default-machine"
	}
	return hostname
}

func (s *SaveAgent) SaveGame(player *PlayerAgent, market *MarketAgent) error {
	data := SaveData{
		Balance:        player.GetBalance(),
		Inventory:      *player.GetInventory(),
		Pity:           player.Pity,
		LastRefreshDay: market.GetLastRefreshDay(),
	}
	
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	
	encrypted, err := s.encrypt(jsonData)
	if err != nil {
		return err
	}
	
	return os.WriteFile("save.dat", encrypted, 0644)
}

func (s *SaveAgent) LoadGame(player *PlayerAgent, market *MarketAgent) error {
	encrypted, err := os.ReadFile("save.dat")
	if err != nil {
		return err
	}
	
	decrypted, err := s.decrypt(encrypted)
	if err != nil {
		return err
	}
	
	var data SaveData
	if err := json.Unmarshal(decrypted, &data); err != nil {
		return err
	}
	
	player.balance = data.Balance
	*player.GetInventory() = data.Inventory
	player.Pity = data.Pity
	market.SetLastRefreshDay(data.LastRefreshDay)
	
	return nil
}

func (s *SaveAgent) encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return nil, err
	}
	
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	
	return gcm.Seal(nonce, nonce, data, nil), nil
}

func (s *SaveAgent) decrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return nil, err
	}
	
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, err
	}
	
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
