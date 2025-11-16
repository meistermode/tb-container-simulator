package models

type RewardType string

const (
	RewardGold           RewardType = "gold"
	RewardCredits        RewardType = "credits"
	RewardFreeXP         RewardType = "free_xp"
	RewardCert200        RewardType = "cert_200"
	RewardCert2300       RewardType = "cert_2300"
	RewardCert28000      RewardType = "cert_28000"
	RewardMystCert       RewardType = "myst_cert"
	RewardX3XP           RewardType = "x3_xp"
	RewardX4XP           RewardType = "x4_xp"
	RewardX5XP           RewardType = "x5_xp"
	RewardTalismanL2     RewardType = "talisman_l2"
	RewardTalismanL3     RewardType = "talisman_l3"
	RewardTankoin        RewardType = "tankoin"
	RewardAvatarSilver   RewardType = "avatar_silver"
	RewardAvatarGold     RewardType = "avatar_gold"
	RewardAvatarGoldAnim RewardType = "avatar_gold_anim"
	RewardProfileBG      RewardType = "profile_bg"
	RewardRiddle         RewardType = "riddle"
	RewardContainerL2    RewardType = "container_l2"
	RewardContainerL3    RewardType = "container_l3"
)

type Reward struct {
	Type   RewardType
	Amount int
}

type LootEntry struct {
	Type   RewardType `json:"type"`
	Amount int        `json:"amount"`
	Chance float64    `json:"chance"`
}

type LootTable struct {
	Guaranteed []LootEntry `json:"guaranteed"`
	Random     []LootEntry `json:"random"`
}

type Inventory struct {
	Gold           int
	Credits        int
	FreeXP         int
	Cert200        int
	Cert2300       int
	Cert28000      int
	MystCert       int
	X3XP           int
	X4XP           int
	X5XP           int
	TalismanL2     int
	TalismanL3     int
	Tankoin        int
	AvatarSilver   int
	AvatarGold     int
	AvatarGoldAnim int
	ProfileBG      int
	Riddle         int
	ContainerL1    int
	ContainerL2    int
	ContainerL3    int
}

type PityCounters struct {
	L1 int
	L2 int
	L3 int
}
