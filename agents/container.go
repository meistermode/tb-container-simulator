package agents

import (
	"container-simulator/models"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/json"
	"os"
)

type ContainerAgent struct {
	lootTables map[int]models.LootTable
	rng        *RNGAgent
}

func NewContainerAgent(rng *RNGAgent) *ContainerAgent {
	ca := &ContainerAgent{
		lootTables: make(map[int]models.LootTable),
		rng:        rng,
	}
	ca.loadLootTables()
	return ca
}

func (c *ContainerAgent) loadLootTables() {
	encrypted, err := os.ReadFile("data/loot_tables.dat")
	if err != nil {
		panic(err)
	}
	
	data, err := c.decryptLootTables(encrypted)
	if err != nil {
		panic(err)
	}
	
	var tables map[string]models.LootTable
	if err := json.Unmarshal(data, &tables); err != nil {
		panic(err)
	}
	
	c.lootTables[1] = tables["l1"]
	c.lootTables[2] = tables["l2"]
	c.lootTables[3] = tables["l3"]
}

func (c *ContainerAgent) decryptLootTables(data []byte) ([]byte, error) {
	key := sha256.Sum256([]byte("loot-tables-encryption-key-v1"))
	
	block, err := aes.NewCipher(key[:])
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

func (c *ContainerAgent) OpenContainer(level int) []models.Reward {
	table := c.lootTables[level]
	rewards := []models.Reward{}
	
	for _, entry := range table.Guaranteed {
		amount := entry.Amount
		if entry.Type == models.RewardX4XP && level == 2 {
			if c.rng.Float64() < 0.5 {
				amount = 1
			} else {
				amount = 2
			}
		} else if entry.Type == models.RewardX5XP && level == 3 {
			roll := c.rng.Float64()
			if roll < 0.3333 {
				amount = 1
			} else if roll < 0.6666 {
				amount = 2
			} else {
				amount = 3
			}
		} else if entry.Type == models.RewardCert2300 && level == 2 {
			if c.rng.Float64() < 0.5 {
				amount = 1
			} else {
				amount = 2
			}
		}
		rewards = append(rewards, models.Reward{Type: entry.Type, Amount: amount})
	}
	
	for _, entry := range table.Random {
		if c.rng.Float64()*100 < entry.Chance {
			rewards = append(rewards, models.Reward{Type: entry.Type, Amount: entry.Amount})
		}
	}
	
	return rewards
}
