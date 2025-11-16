package agents

import (
	"bufio"
	"container-simulator/models"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

type SimulationAgent struct {
	player *PlayerAgent
	market *MarketAgent
	save   *SaveAgent
	reader *bufio.Reader
}

func NewSimulationAgent(player *PlayerAgent, market *MarketAgent, save *SaveAgent) *SimulationAgent {
	return &SimulationAgent{
		player: player,
		market: market,
		save:   save,
		reader: bufio.NewReader(os.Stdin),
	}
}

func (s *SimulationAgent) Run() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	
	go func() {
		<-sigChan
		s.save.SaveGame(s.player, s.market)
		fmt.Println("\nИгра сохранена. Выход...")
		os.Exit(0)
	}()
	
	if err := s.save.LoadGame(s.player, s.market); err == nil {
		s.clearScreen()
		fmt.Println("Сохранение загружено!")
		
		if s.market.CheckAndAutoRefresh() {
			fmt.Println("Магазин автоматически обновлён! (новый день)")
		}
		
		fmt.Print("\nНажмите Enter для продолжения...")
		s.readInput()
	} else {
		s.clearScreen()
		fmt.Println("Добро пожаловать в симулятор контейнеров!")
		s.market.RefreshMarket()
		fmt.Print("\nНажмите Enter для продолжения...")
		s.readInput()
	}
	
	for {
		s.clearScreen()
		s.showMainMenu()
		choice := s.readInput()
		
		switch choice {
		case "1":
			s.showMarket()
			s.save.SaveGame(s.player, s.market)
		case "2":
			s.showInventory()
		case "3":
			s.openContainers()
			s.save.SaveGame(s.player, s.market)
		case "4":
			s.clearScreen()
			s.save.SaveGame(s.player, s.market)
			fmt.Println("Игра сохранена. Выход...")
			return
		default:
			fmt.Println("Неверный выбор. Попробуйте снова.")
			fmt.Print("Нажмите Enter...")
			s.readInput()
		}
	}
}

func (s *SimulationAgent) clearScreen() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func (s *SimulationAgent) showMainMenu() {
	fmt.Println("=== ГЛАВНОЕ МЕНЮ ===")
	fmt.Printf("Баланс: %d золота\n\n", s.player.GetBalance())
	fmt.Println("1. Магазин")
	fmt.Println("2. Хранилище")
	fmt.Println("3. Открыть контейнеры")
	fmt.Println("4. Выход")
	fmt.Print("\nВаш выбор: ")
}

func (s *SimulationAgent) showMarket() {
	for {
		s.clearScreen()
		fmt.Println("=== МАГАЗИН ===")
		fmt.Printf("Баланс: %d золота\n\n", s.player.GetBalance())
		
		packs := s.market.GetAllPacks()
		for i, pack := range packs {
			if pack.Stock < 0 {
				fmt.Printf("%d. Купить %s — %d\n", i+1, pack.Name, pack.Price)
			} else {
				fmt.Printf("%d. Купить %s — %d (осталось: %d)\n", i+1, pack.Name, pack.Price, pack.Stock)
			}
		}
		fmt.Println("0. Назад в меню")
		fmt.Print("\nВыберите действие: ")
		
		choice := s.readInput()
		
		switch choice {
		case "1":
			s.buyPack(Pack10)
			fmt.Print("\nНажмите Enter...")
			s.readInput()
		case "2":
			s.buyPack(Pack50)
			fmt.Print("\nНажмите Enter...")
			s.readInput()
		case "3":
			s.buyPack(Pack25)
			fmt.Print("\nНажмите Enter...")
			s.readInput()
		case "4":
			s.buyPack(Pack17)
			fmt.Print("\nНажмите Enter...")
			s.readInput()
		case "5":
			s.buyPack(Pack13)
			fmt.Print("\nНажмите Enter...")
			s.readInput()
		case "6":
			s.buyPack(Pack9)
			fmt.Print("\nНажмите Enter...")
			s.readInput()
		case "7":
			s.buyPack(Pack5)
			fmt.Print("\nНажмите Enter...")
			s.readInput()
		case "8":
			s.buyPack(Pack2)
			fmt.Print("\nНажмите Enter...")
			s.readInput()
		case "0":
			return
		default:
			fmt.Println("Неверный выбор.")
			fmt.Print("Нажмите Enter...")
			s.readInput()
		}
	}
}

func (s *SimulationAgent) buyPack(packType PackType) {
	pack := s.market.GetPack(packType)
	
	if pack.Stock == 0 {
		fmt.Println("Этот набор закончился!")
		return
	}
	
	if !s.player.SpendGold(pack.Price) {
		fmt.Println("Недостаточно золота!")
		return
	}
	
	containers, _ := s.market.BuyPack(packType)
	s.player.AddContainers(1, containers)
	fmt.Printf("Вы купили %s! Получено %d контейнеров 1 уровня.\n", pack.Name, containers)
}

func (s *SimulationAgent) showInventory() {
	s.clearScreen()
	inv := s.player.GetInventory()
	
	fmt.Println("=== ХРАНИЛИЩЕ ===")
	
	hasContainers := inv.ContainerL1 > 0 || inv.ContainerL2 > 0 || inv.ContainerL3 > 0
	if hasContainers {
		fmt.Println("Контейнеры:")
		if inv.ContainerL1 > 0 {
			fmt.Printf("- 1 уровня: %d\n", inv.ContainerL1)
		}
		if inv.ContainerL2 > 0 {
			fmt.Printf("- 2 уровня: %d\n", inv.ContainerL2)
		}
		if inv.ContainerL3 > 0 {
			fmt.Printf("- 3 уровня: %d\n", inv.ContainerL3)
		}
	}
	
	hasResources := inv.Credits > 0 || inv.FreeXP > 0 || 
		inv.Cert200 > 0 || inv.Cert2300 > 0 || inv.Cert28000 > 0 || inv.MystCert > 0 ||
		inv.X3XP > 0 || inv.X4XP > 0 || inv.X5XP > 0 ||
		inv.TalismanL2 > 0 || inv.TalismanL3 > 0 || inv.Tankoin > 0 ||
		inv.AvatarSilver > 0 || inv.AvatarGold > 0 || inv.AvatarGoldAnim > 0 ||
		inv.ProfileBG > 0 || inv.Riddle > 0
	
	if hasResources {
		fmt.Println("\nРесурсы:")
		if inv.Credits > 0 {
			fmt.Printf("- Кредиты: %d\n", inv.Credits)
		}
		if inv.FreeXP > 0 {
			fmt.Printf("- Свободный опыт: %d\n", inv.FreeXP)
		}
		if inv.Cert200 > 0 {
			fmt.Printf("- Сертификаты (200): %d\n", inv.Cert200)
		}
		if inv.Cert2300 > 0 {
			fmt.Printf("- Сертификаты (2300): %d\n", inv.Cert2300)
		}
		if inv.Cert28000 > 0 {
			fmt.Printf("- Сертификаты (28000): %d\n", inv.Cert28000)
		}
		if inv.MystCert > 0 {
			fmt.Printf("- Мистические сертификаты: %d\n", inv.MystCert)
		}
		if inv.X3XP > 0 {
			fmt.Printf("- x3 опыта: %d\n", inv.X3XP)
		}
		if inv.X4XP > 0 {
			fmt.Printf("- x4 опыта: %d\n", inv.X4XP)
		}
		if inv.X5XP > 0 {
			fmt.Printf("- x5 опыта: %d\n", inv.X5XP)
		}
		if inv.TalismanL2 > 0 {
			fmt.Printf("- Талисманы L2: %d\n", inv.TalismanL2)
		}
		if inv.TalismanL3 > 0 {
			fmt.Printf("- Талисманы L3: %d\n", inv.TalismanL3)
		}
		if inv.Tankoin > 0 {
			fmt.Printf("- Танкоины: %d\n", inv.Tankoin)
		}
		if inv.AvatarSilver > 0 {
			fmt.Printf("- Аватары (В серебре): %d\n", inv.AvatarSilver)
		}
		if inv.AvatarGold > 0 {
			fmt.Printf("- Аватары (В золоте): %d\n", inv.AvatarGold)
		}
		if inv.AvatarGoldAnim > 0 {
			fmt.Printf("- Аватары (анимированные): %d\n", inv.AvatarGoldAnim)
		}
		if inv.ProfileBG > 0 {
			fmt.Printf("- Фоны профиля: %d\n", inv.ProfileBG)
		}
		if inv.Riddle > 0 {
			fmt.Printf("- Загадки: %d\n", inv.Riddle)
		}
	}
	
	if !hasContainers && !hasResources {
		fmt.Println("\nХранилище пусто.")
	}
	
	fmt.Print("\n1. Назад\nВаш выбор: ")
	s.readInput()
}

func (s *SimulationAgent) openContainers() {
	for {
		s.clearScreen()
		inv := s.player.GetInventory()
		fmt.Println("=== ОТКРЫТИЕ КОНТЕЙНЕРОВ ===")
		fmt.Printf("1 уровня: %d | 2 уровня: %d | 3 уровня: %d\n\n", inv.ContainerL1, inv.ContainerL2, inv.ContainerL3)
		fmt.Println("1. Открыть 1 уровень")
		fmt.Println("2. Открыть 2 уровень")
		fmt.Println("3. Открыть 3 уровень")
		fmt.Println("4. Открыть всё")
		fmt.Println("5. Назад")
		fmt.Print("\nВыберите действие: ")
		
		choice := s.readInput()
		
		switch choice {
		case "1":
			s.openContainerLoop(1)
		case "2":
			s.openContainerLoop(2)
		case "3":
			s.openContainerLoop(3)
		case "4":
			s.openAllContainers()
			fmt.Print("\nНажмите Enter...")
			s.readInput()
		case "5":
			return
		default:
			fmt.Println("Неверный выбор.")
			fmt.Print("Нажмите Enter...")
			s.readInput()
		}
	}
}

func (s *SimulationAgent) openContainerLoop(level int) {
	for {
		inv := s.player.GetInventory()
		var count int
		switch level {
		case 1:
			count = inv.ContainerL1
		case 2:
			count = inv.ContainerL2
		case 3:
			count = inv.ContainerL3
		}
		
		if count == 0 {
			s.clearScreen()
			fmt.Printf("У вас нет контейнеров %d уровня!\n", level)
			fmt.Print("\nНажмите Enter...")
			s.readInput()
			return
		}
		
		s.clearScreen()
		fmt.Println("=== ПОЗДРАВЛЯЕМ! ===\n")
		
		rewards, pityTriggered := s.player.OpenContainer(level)
		if rewards == nil {
			return
		}
		
		fmt.Println("Выпало:")
		s.displayRewards(rewards, pityTriggered)
		
		// Обновляем количество после открытия
		inv = s.player.GetInventory()
		switch level {
		case 1:
			count = inv.ContainerL1
		case 2:
			count = inv.ContainerL2
		case 3:
			count = inv.ContainerL3
		}
		
		fmt.Printf("\nОставшиеся контейнеры %d уровня: %d\n", level, count)
		
		if count > 0 {
			fmt.Println("\n1. Открыть")
			fmt.Println("2. Назад")
			fmt.Print("\nВыберите действие: ")
			
			choice := s.readInput()
			if choice != "1" {
				return
			}
		} else {
			fmt.Println("\nКонтейнеры закончились!")
			fmt.Print("\nНажмите Enter...")
			s.readInput()
			return
		}
	}
}

func (s *SimulationAgent) openAllContainers() {
	inv := s.player.GetInventory()
	total := inv.ContainerL1 + inv.ContainerL2 + inv.ContainerL3
	
	if total == 0 {
		fmt.Println("У вас нет контейнеров!")
		return
	}
	
	opened := 0
	
	for inv.ContainerL1 > 0 {
		s.player.OpenContainer(1)
		opened++
	}
	for inv.ContainerL2 > 0 {
		s.player.OpenContainer(2)
		opened++
	}
	for inv.ContainerL3 > 0 {
		s.player.OpenContainer(3)
		opened++
	}
	
	fmt.Printf("\nОткрыто контейнеров: %d\n", opened)
	fmt.Println("Все награды добавлены в хранилище.")
}

func (s *SimulationAgent) displayRewards(rewards []models.Reward, pityTriggered bool) {
	rewardMap := make(map[models.RewardType]int)
	pityRewards := make(map[models.RewardType]bool)
	
	for _, reward := range rewards {
		isPity := pityTriggered && (reward.Type == models.RewardContainerL2 || reward.Type == models.RewardContainerL3)
		rewardMap[reward.Type] += reward.Amount
		if isPity {
			pityRewards[reward.Type] = true
		}
	}
	
	for rewardType, amount := range rewardMap {
		name := s.getRewardName(rewardType)
		hasPity := pityRewards[rewardType]
		
		if amount > 1 {
			if hasPity {
				fmt.Printf("- %s ×%d (включая гарант)\n", name, amount)
			} else {
				fmt.Printf("- %s ×%d\n", name, amount)
			}
		} else {
			if hasPity {
				fmt.Printf("- %s (гарант)\n", name)
			} else {
				fmt.Printf("- %s\n", name)
			}
		}
	}
}

func (s *SimulationAgent) getRewardName(rewardType models.RewardType) string {
	names := map[models.RewardType]string{
		models.RewardGold:           "Золото",
		models.RewardCredits:        "Кредиты",
		models.RewardFreeXP:         "Свободный опыт",
		models.RewardCert200:        "Сертификат на 200 опыта",
		models.RewardCert2300:       "Сертификат на 2300 опыта",
		models.RewardCert28000:      "Сертификат на 28000 опыта",
		models.RewardMystCert:       "Мистический сертификат",
		models.RewardX3XP:           "×3 опыта",
		models.RewardX4XP:           "×4 опыта",
		models.RewardX5XP:           "×5 опыта",
		models.RewardTalismanL2:     "Талисман контейнера II",
		models.RewardTalismanL3:     "Талисман контейнера III",
		models.RewardTankoin:        "Танкоин",
		models.RewardAvatarSilver:   "Аватар «В серебре»",
		models.RewardAvatarGold:     "Аватар «В золоте»",
		models.RewardAvatarGoldAnim: "Анимированный аватар «В золоте»",
		models.RewardProfileBG:      "Фон профиля «Из сумрака»",
		models.RewardRiddle:         "Загадка",
		models.RewardContainerL2:    "Контейнер II уровня",
		models.RewardContainerL3:    "Контейнер III уровня",
	}
	return names[rewardType]
}

func (s *SimulationAgent) readInput() string {
	input, _ := s.reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func (s *SimulationAgent) readInt() int {
	input := s.readInput()
	val, _ := strconv.Atoi(input)
	return val
}
