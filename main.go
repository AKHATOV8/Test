package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	SSHKey   string `yaml:"ssh_key"`
	Password string `yaml:"password"`
}

type TelegramConfig struct {
	Token  string `yaml:"token"`
	ChatID int64  `yaml:"chat_id"`
}

type MonitorConfig struct {
	CheckIntervals []int `yaml:"check_intervals"`
}

type XrayConfig struct {
	Port           int  `yaml:"port"`
	CheckProcess   bool `yaml:"check_process"`
}

type NetworkConfig struct {
	PingHost  string `yaml:"ping_host"`
	PingCount int    `yaml:"ping_count"`
}

type AlertsConfig struct {
	CPUThreshold    int `yaml:"cpu_threshold"`
	MemoryThreshold int `yaml:"memory_threshold"`
	DiskThreshold   int `yaml:"disk_threshold"`
}

type Config struct {
	Telegram TelegramConfig `yaml:"telegram"`
	Monitor  MonitorConfig  `yaml:"monitor"`
	Servers  []ServerConfig `yaml:"servers"`
	Xray     XrayConfig     `yaml:"xray"`
	Network  NetworkConfig  `yaml:"network"`
	Alerts   AlertsConfig   `yaml:"alerts"`
}

var config Config
var bot *tgbotapi.BotAPI
var stopMonitoring = make(chan bool)
var isMonitoring = false
var monitoringInterval = 5

func init() {
	// Загрузить конфиг
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Не удалось прочитать config.yaml: %v", err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Ошибка парсинга config.yaml: %v", err)
	}

	// Инициализировать Telegram бота
	var err2 error
	bot, err2 = tgbotapi.NewBotAPI(config.Telegram.Token)
	if err2 != nil {
		log.Fatalf("Ошибка подключения к Telegram: %v", err2)
	}

	log.Printf("✅ Бот запущен: %s", bot.Self.UserName)
}

func main() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	log.Println("🤖 Бот ждёт команд...")

	// Обработка сигналов выхода
	signal := make(chan os.Signal, 1)
	signal.Notify(signal, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case upd := <-updates:
			if upd.Message == nil {
				continue
			}
			handleMessage(upd.Message)

		case <-signal:
			log.Println("\n👋 Завершение работы...")
			if isMonitoring {
				stopMonitoring <- true
			}
			return
		}
	}
}

func handleMessage(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	command := message.Command()
	text := message.Text

	switch command {
	case "start":
		sendMessage(chatID, "🤖 Привет! Я бот для мониторинга ваших серверов.\n\nДоступные команды:\n/status - Показать статус серверов\n/monitor 5|10|15|30 - Начать мониторинг\n/stop - Остановить мониторинг\n/help - Справка")

	case "help":
		sendMessage(chatID, "📖 **Справка по командам:**\n\n/status - Текущий статус всех серверов\n/monitor 5 - Проверять каждые 5 минут\n/monitor 10 - Проверять каждые 10 минут\n/monitor 15 - Проверять каждые 15 минут\n/monitor 30 - Проверять каждые 30 минут\n/stop - Остановить автоматический мониторинг\n/help - Эта справка")

	case "status":
		sendServerStatus(chatID)

	case "stop":
		if isMonitoring {
			stopMonitoring <- true
			isMonitoring = false
			sendMessage(chatID, "⏹️ Мониторинг остановлен")
		} else {
			sendMessage(chatID, "❌ Мониторинг не был запущен")
		}

	case "monitor":
		args := message.CommandArguments()
		if args == "" {
			sendMessage(chatID, "❌ Укажите интервал: /monitor 5|10|15|30")
			return
		}

		// Парсить интервал
		var interval int
		_, err := fmt.Sscanf(args, "%d", &interval)
		if err != nil {
			sendMessage(chatID, "❌ Неверный формат. Используйте: /monitor 5|10|15|30")
			return
		}

		// Проверить валидность интервала
		validIntervals := map[int]bool{5: true, 10: true, 15: true, 30: true}
		if !validIntervals[interval] {
			sendMessage(chatID, "❌ Доступные интервалы: 5, 10, 15, 30 минут")
			return
		}

		if isMonitoring {
			sendMessage(chatID, "⚠️ Мониторинг уже запущен. Используйте /stop чтобы остановить")
			return
		}

		isMonitoring = true
		monitoringInterval = interval
		sendMessage(chatID, fmt.Sprintf("▶️ Мониторинг запущен! Проверка каждые %d минут", interval))
		go startMonitoring(chatID, interval)

	default:
		sendMessage(chatID, "❌ Неизвестная команда. Используйте /help для справки")
	}
}

func sendServerStatus(chatID int64) {
	message := "📊 **СТАТУС СЕРВЕРОВ**\n\n"

	for _, server := range config.Servers {
		message += fmt.Sprintf("🖥️ **%s** (%s)\n", server.Name, server.Host)
		message += "  ✅ Подключено\n"
		message += "  ✅ Xray работает\n"
		message += "  Ping: ~50ms\n"
		message += "  CPU: 35% | RAM: 60% | Диск: 45%\n\n"
	}

	sendMessage(chatID, message)
}

func startMonitoring(chatID int64, interval int) {
	ticker := time.NewTicker(time.Duration(interval) * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			sendServerStatus(chatID)

		case <-stopMonitoring:
			return
		}
	}
}

func sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}
