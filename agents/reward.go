package agents

import (
	"container-simulator/models"
)

type RewardAgent struct {
	inventory *models.Inventory
}

func NewRewardAgent() *RewardAgent {
	return &RewardAgent{
		inventory: &models.Inventory{},
	}
}

func (r *RewardAgent) AddReward(reward models.Reward) {
	switch reward.Type {
	case models.RewardGold:
		r.inventory.Gold += reward.Amount
	case models.RewardCredits:
		r.inventory.Credits += reward.Amount
	case models.RewardFreeXP:
		r.inventory.FreeXP += reward.Amount
	case models.RewardCert200:
		r.inventory.Cert200 += reward.Amount
	case models.RewardCert2300:
		r.inventory.Cert2300 += reward.Amount
	case models.RewardCert28000:
		r.inventory.Cert28000 += reward.Amount
	case models.RewardMystCert:
		r.inventory.MystCert += reward.Amount
	case models.RewardX3XP:
		r.inventory.X3XP += reward.Amount
	case models.RewardX4XP:
		r.inventory.X4XP += reward.Amount
	case models.RewardX5XP:
		r.inventory.X5XP += reward.Amount
	case models.RewardTalismanL2:
		r.inventory.TalismanL2 += reward.Amount
	case models.RewardTalismanL3:
		r.inventory.TalismanL3 += reward.Amount
	case models.RewardTankoin:
		r.inventory.Tankoin += reward.Amount
	case models.RewardAvatarSilver:
		r.inventory.AvatarSilver += reward.Amount
	case models.RewardAvatarGold:
		r.inventory.AvatarGold += reward.Amount
	case models.RewardAvatarGoldAnim:
		r.inventory.AvatarGoldAnim += reward.Amount
	case models.RewardProfileBG:
		r.inventory.ProfileBG += reward.Amount
	case models.RewardRiddle:
		r.inventory.Riddle += reward.Amount
	case models.RewardContainerL2:
		r.inventory.ContainerL2 += reward.Amount
	case models.RewardContainerL3:
		r.inventory.ContainerL3 += reward.Amount
	}
}

func (r *RewardAgent) GetInventory() *models.Inventory {
	return r.inventory
}
