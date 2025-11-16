package agents

import (
	"container-simulator/models"
)

type PlayerAgent struct {
	balance        int
	Pity           models.PityCounters
	rewardAgent    *RewardAgent
	containerAgent *ContainerAgent
}

func NewPlayerAgent(startBalance int, rewardAgent *RewardAgent, containerAgent *ContainerAgent) *PlayerAgent {
	return &PlayerAgent{
		balance:        startBalance,
		Pity:           models.PityCounters{},
		rewardAgent:    rewardAgent,
		containerAgent: containerAgent,
	}
}

func (p *PlayerAgent) GetBalance() int {
	return p.balance
}

func (p *PlayerAgent) AddGold(amount int) {
	p.balance += amount
}

func (p *PlayerAgent) SpendGold(amount int) bool {
	if p.balance < amount {
		return false
	}
	p.balance -= amount
	return true
}

func (p *PlayerAgent) AddContainers(level, count int) {
	inv := p.rewardAgent.GetInventory()
	switch level {
	case 1:
		inv.ContainerL1 += count
	case 2:
		inv.ContainerL2 += count
	case 3:
		inv.ContainerL3 += count
	}
}

func (p *PlayerAgent) OpenContainer(level int) ([]models.Reward, bool) {
	inv := p.rewardAgent.GetInventory()
	
	switch level {
	case 1:
		if inv.ContainerL1 <= 0 {
			return nil, false
		}
		inv.ContainerL1--
	case 2:
		if inv.ContainerL2 <= 0 {
			return nil, false
		}
		inv.ContainerL2--
	case 3:
		if inv.ContainerL3 <= 0 {
			return nil, false
		}
		inv.ContainerL3--
	}
	
	rewards := p.containerAgent.OpenContainer(level)
	
	for _, reward := range rewards {
		if reward.Type == models.RewardGold {
			p.balance += reward.Amount
		} else {
			p.rewardAgent.AddReward(reward)
		}
	}
	
	pityTriggered := false
	switch level {
	case 1:
		p.Pity.L1++
		if p.Pity.L1 >= 10 {
			p.Pity.L1 = 0
			inv.ContainerL2++
			inv.TalismanL2 = 0
			pityTriggered = true
			rewards = append(rewards, models.Reward{Type: models.RewardContainerL2, Amount: 1})
		}
	case 2:
		p.Pity.L2++
		if p.Pity.L2 >= 10 {
			p.Pity.L2 = 0
			inv.ContainerL3++
			inv.TalismanL2 = 0
			pityTriggered = true
			rewards = append(rewards, models.Reward{Type: models.RewardContainerL3, Amount: 1})
		}
	case 3:
		p.Pity.L3++
		if p.Pity.L3 >= 10 {
			p.Pity.L3 = 0
			inv.ContainerL3++
			inv.TalismanL3 = 0
			pityTriggered = true
			rewards = append(rewards, models.Reward{Type: models.RewardContainerL3, Amount: 1})
		}
	}
	
	return rewards, pityTriggered
}

func (p *PlayerAgent) GetInventory() *models.Inventory {
	return p.rewardAgent.GetInventory()
}
