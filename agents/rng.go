package agents

import (
	"math/rand"
	"time"
)

type RNGAgent struct {
	rng *rand.Rand
}

func NewRNGAgent() *RNGAgent {
	return &RNGAgent{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (r *RNGAgent) Float64() float64 {
	return r.rng.Float64()
}

func (r *RNGAgent) Intn(n int) int {
	return r.rng.Intn(n)
}
