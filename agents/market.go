package agents

import (
	"time"
)

type PackType int

const (
	Pack10 PackType = iota
	Pack50
	Pack25
	Pack17
	Pack13
	Pack9
	Pack5
	Pack2
)

type Pack struct {
	Name       string
	Price      int
	Containers int
	MaxStock   int
	Stock      int
}

type MarketAgent struct {
	packs          map[PackType]*Pack
	lastRefreshDay string
}

func NewMarketAgent() *MarketAgent {
	m := &MarketAgent{
		packs:          make(map[PackType]*Pack),
		lastRefreshDay: "",
	}
	m.initPacks()
	return m
}

func (m *MarketAgent) initPacks() {
	m.packs[Pack10] = &Pack{
		Name:       "Мистический контейнер 1 уровня x10",
		Price:      10000,
		Containers: 10,
		MaxStock:   -1,
		Stock:      -1,
	}
	m.packs[Pack50] = &Pack{
		Name:       "Мистический контейнер 1 уровня x50",
		Price:      22500,
		Containers: 50,
		MaxStock:   2,
		Stock:      2,
	}
	m.packs[Pack25] = &Pack{
		Name:       "Мистический контейнер 1 уровня x25",
		Price:      12500,
		Containers: 25,
		MaxStock:   2,
		Stock:      2,
	}
	m.packs[Pack17] = &Pack{
		Name:       "Мистический контейнер 1 уровня x17",
		Price:      6000,
		Containers: 17,
		MaxStock:   1,
		Stock:      1,
	}
	m.packs[Pack13] = &Pack{
		Name:       "Мистический контейнер 1 уровня x13",
		Price:      4500,
		Containers: 13,
		MaxStock:   1,
		Stock:      1,
	}
	m.packs[Pack9] = &Pack{
		Name:       "Мистический контейнер 1 уровня x9",
		Price:      3000,
		Containers: 9,
		MaxStock:   1,
		Stock:      1,
	}
	m.packs[Pack5] = &Pack{
		Name:       "Мистический контейнер 1 уровня x5",
		Price:      1500,
		Containers: 5,
		MaxStock:   1,
		Stock:      1,
	}
	m.packs[Pack2] = &Pack{
		Name:       "Мистический контейнер 1 уровня x2",
		Price:      500,
		Containers: 2,
		MaxStock:   1,
		Stock:      1,
	}
}

func (m *MarketAgent) GetPack(packType PackType) *Pack {
	return m.packs[packType]
}

func (m *MarketAgent) BuyPack(packType PackType) (int, bool) {
	pack := m.packs[packType]
	if pack.Stock == 0 {
		return 0, false
	}
	if pack.Stock > 0 {
		pack.Stock--
	}
	return pack.Containers, true
}

func (m *MarketAgent) RefreshMarket() {
	for _, pack := range m.packs {
		pack.Stock = pack.MaxStock
	}
	m.lastRefreshDay = getCurrentRefreshDay()
}

func (m *MarketAgent) CheckAndAutoRefresh() bool {
	currentDay := getCurrentRefreshDay()
	if m.lastRefreshDay != currentDay {
		m.RefreshMarket()
		return true
	}
	return false
}

func (m *MarketAgent) GetLastRefreshDay() string {
	return m.lastRefreshDay
}

func (m *MarketAgent) SetLastRefreshDay(day string) {
	m.lastRefreshDay = day
}

func getCurrentRefreshDay() string {
	msk := time.FixedZone("MSK", 3*60*60)
	now := time.Now().In(msk)
	
	if now.Hour() < 12 {
		now = now.AddDate(0, 0, -1)
	}
	
	return now.Format("2006-01-02")
}

func (m *MarketAgent) GetAllPacks() []*Pack {
	return []*Pack{
		m.packs[Pack10],
		m.packs[Pack50],
		m.packs[Pack25],
		m.packs[Pack17],
		m.packs[Pack13],
		m.packs[Pack9],
		m.packs[Pack5],
		m.packs[Pack2],
	}
}
