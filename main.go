package main

import (
	"container-simulator/agents"
)

func main() {
	rng := agents.NewRNGAgent()
	rewardAgent := agents.NewRewardAgent()
	containerAgent := agents.NewContainerAgent(rng)
	playerAgent := agents.NewPlayerAgent(50000, rewardAgent, containerAgent)
	marketAgent := agents.NewMarketAgent()
	saveAgent := agents.NewSaveAgent()
	
	simulation := agents.NewSimulationAgent(playerAgent, marketAgent, saveAgent)
	simulation.Run()
}
