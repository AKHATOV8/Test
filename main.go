package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
)

// Config structures
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
	Port         int  `yaml:"port"`
	CheckProcess bool `yaml:"check_process"`
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

// Server Status
type ServerStatus struct {
	Name       string
	Host       string
	Online     bool
	XrayStatus bool
	XrayError  string
	Ping       string
	CPU        int
	Memory     int
	Disk       int
	Errors     []string
}

var config Config
var bot *tgbotapi.BotAPI
var stopMonitoring = make(chan bool, 1)
var isMonitoring = false
var monitoringInterval = 5
var sshConnections = make(map[string]*ssh.Client)
var sshMutex = &sync.Mutex{}

func init() {
	// Load config
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("❌ Failed to read config.yaml: %v", err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("❌ Failed to parse config.yaml: %v", err)
	}

	// Initialize Telegram bot
	var err2 error
	bot, err2 = tgbotapi.NewBotAPI(config.Telegram.Token)
	if err2 != nil {
		log.Fatalf("❌ Failed to connect to Telegram: %v", err2)
	}

	log.Printf("✅ Bot started: @%s", bot.Self.UserName)
	log.Printf("📊 Monitoring %d servers", len(config.Servers))
}

func main() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	log.Println("🤖 Bot is waiting for commands...")

	// Handle graceful shutdown
	signal := make(chan os.Signal, 1)
	signal.Notify(signal, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case upd := <-updates:
			if upd.Message == nil {
				continue
			}
			go handleMessage(upd.Message)

		case <-signal:
			log.Println("\n👋 Shutting down...")
			if isMonitoring {
				stopMonitoring <- true
			}
			closeAllSSHConnections()
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
		sendMessage(chatID, "🤖 *Привет!* Я бот для мониторинга ваших Linux серверов и Xray.\n\n*Доступные команды:*\n/status - Показать статус серверов\n/servers - Список всех серверов\n/monitor 5|10|15|30 - Начать мониторинг\n/stop - Остановить мониторинг\n/help - Справка")

	case "help":
		help := "📖 *Справка по командам:*\n\n"
		help += "/status - Текущий статус всех серверов\n"
		help += "/servers - Список всех серверов (до 15)\n"
		help += "/monitor 5 - Проверять каждые 5 минут\n"
		help += "/monitor 10 - Проверять каждые 10 минут\n"
		help += "/monitor 15 - Проверять каждые 15 минут\n"
		help += "/monitor 30 - Проверять каждые 30 минут\n"
		help += "/stop - Остановить автоматический мониторинг\n\n"
		help += "⚙️ *Пороги оповещений:*\n"
		help += "CPU: " + strconv.Itoa(config.Alerts.CPUThreshold) + "%\n"
		help += "RAM: " + strconv.Itoa(config.Alerts.MemoryThreshold) + "%\n"
		help += "Диск: " + strconv.Itoa(config.Alerts.DiskThreshold) + "%\n\n"
		help += "📝 *Как добавить новый сервер:*\n"
		help += "1. Отредактируйте `config.yaml`\n"
		help += "2. Добавьте новый сервер в секцию `servers`\n"
		help += "3. Перезагрузите бота\n"
		help += "4. Поддерживается до 15 серверов"
		sendMessage(chatID, help)

	case "servers":
		sendServersList(chatID)

	case "status":
		sendMessage(chatID, "⏳ Сбор информации со всех серверов...")
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

		var interval int
		_, err := fmt.Sscanf(strings.TrimSpace(args), "%d", &interval)
		if err != nil {
			sendMessage(chatID, "❌ Неверный формат. Используйте: /monitor 5|10|15|30")
			return
		}

		validIntervals := map[int]bool{5: true, 10: true, 15: true, 30: true}
		if !validIntervals[interval] {
			sendMessage(chatID, "❌ Доступные интервалы: 5, 10, 15, 30 минут")
			return
		}

		if isMonitoring {
			sendMessage(chatID, "⚠️ Мониторинг уже запущен. Используйте /stop чтобы остановить")
			return
		}

		if len(config.Servers) == 0 {
			sendMessage(chatID, "❌ Нет добавленных серверов в config.yaml")
			return
		}

		isMonitoring = true
		monitoringInterval = interval
		sendMessage(chatID, fmt.Sprintf("▶️ Мониторинг запущен! Проверка каждые %d минут\n📊 Мониторится %d серверов", interval, len(config.Servers)))
		go startMonitoring(chatID, interval)

	default:
		if text != "" && strings.HasPrefix(text, "/") {
			sendMessage(chatID, "❌ Неизвестная команда. Используйте /help для справки")
		}
	}
}

func sendServersList(chatID int64) {
	if len(config.Servers) == 0 {
		sendMessage(chatID, "❌ Нет добавленных серверов")
		return
	}

	message := fmt.Sprintf("📋 *СПИСОК СЕРВЕРОВ* (%d из 15)\n\n", len(config.Servers))

	for i, server := range config.Servers {
		message += fmt.Sprintf("%d. *%s*\n", i+1, server.Name)
		message += fmt.Sprintf("   Host: `%s:%d`\n", server.Host, server.Port)
		message += fmt.Sprintf("   User: `%s`\n\n", server.User)
	}

	message += "💡 _Для добавления серверов отредактируйте config.yaml (макс. 15 серверов)_"
	sendMessage(chatID, message)
}

func sendServerStatus(chatID int64) {
	if len(config.Servers) == 0 {
		sendMessage(chatID, "❌ Нет добавленных серверов в config.yaml")
		return
	}

	message := fmt.Sprintf("📊 *СТАТУС СЕРВЕРОВ* (%d)\n\n", len(config.Servers))
	statuses := getAllServerStatus()

	for i, status := range statuses {
		message += fmt.Sprintf("*%d. ", i+1)
		message += formatServerStatus(status)
		message += "\n"
	}

	message += "\n_Последнее обновление: " + time.Now().Format("15:04:05") + "_"
	sendMessage(chatID, message)
}

func formatServerStatus(status ServerStatus) string {
	var emoji string
	if status.Online {
		emoji = "🟢"
	} else {
		emoji = "🔴"
	}

	xrayEmoji := "❌"
	if status.XrayStatus {
		xrayEmoji = "✅"
	}

	message := fmt.Sprintf("%s %s* (%s)\n", emoji, status.Name, status.Host)

	if !status.Online {
		message += "  ❌ Сервер не доступен\n"
		return message
	}

	message += fmt.Sprintf("  Xray: %s\n", xrayEmoji)
	if status.XrayError != "" {
		message += fmt.Sprintf("  _Ошибка Xray: %s_\n", status.XrayError)
	}

	message += fmt.Sprintf("  Ping: %s\n", status.Ping)
	message += fmt.Sprintf("  CPU: %d%% | RAM: %d%% | Диск: %d%%\n", status.CPU, status.Memory, status.Disk)

	// Show alerts
	var alerts []string
	if status.CPU >= config.Alerts.CPUThreshold {
		alerts = append(alerts, fmt.Sprintf("⚠️ CPU высокая (%d%%)", status.CPU))
	}
	if status.Memory >= config.Alerts.MemoryThreshold {
		alerts = append(alerts, fmt.Sprintf("⚠️ RAM высокая (%d%%)", status.Memory))
	}
	if status.Disk >= config.Alerts.DiskThreshold {
		alerts = append(alerts, fmt.Sprintf("⚠️ Диск почти заполнен (%d%%)", status.Disk))
	}

	for _, alert := range alerts {
		message += "  " + alert + "\n"
	}

	return message
}

func getAllServerStatus() []ServerStatus {
	var statuses []ServerStatus
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Limit to 15 servers
	maxServers := len(config.Servers)
	if maxServers > 15 {
		maxServers = 15
		log.Printf("⚠️ Warning: More than 15 servers configured, monitoring only first 15")
	}

	for i := 0; i < maxServers; i++ {
		wg.Add(1)
		go func(cfg ServerConfig) {
			defer wg.Done()
			status := getServerStatus(cfg)
			mu.Lock()
			statuses = append(statuses, status)
			mu.Unlock()
		}(config.Servers[i])
	}

	wg.Wait()
	return statuses
}

func getServerStatus(cfg ServerConfig) ServerStatus {
	status := ServerStatus{
		Name: cfg.Name,
		Host: cfg.Host,
	}

	// Connect via SSH
	client, err := getSSHConnection(cfg)
	if err != nil {
		status.Online = false
		status.Errors = append(status.Errors, err.Error())
		return status
	}

	status.Online = true

	// Check Xray status
	status.XrayStatus, status.XrayError = checkXrayStatus(client)

	// Get system stats
	status.CPU = getSystemStat(client, "cpu")
	status.Memory = getSystemStat(client, "memory")
	status.Disk = getSystemStat(client, "disk")

	// Check ping
	status.Ping = checkPing()

	return status
}

func getSSHConnection(cfg ServerConfig) (*ssh.Client, error) {
	sshMutex.Lock()
	defer sshMutex.Unlock()

	if conn, exists := sshConnections[cfg.Host]; exists {
		// Test connection
		if _, _, err := conn.SendRequest("keepalive@golang.com", true, nil); err == nil {
			return conn, nil
		}
		conn.Close()
		delete(sshConnections, cfg.Host)
	}

	var auth []ssh.AuthMethod

	// Try SSH key first
	if cfg.SSHKey != "" {
		key, err := os.ReadFile(cfg.SSHKey)
		if err == nil {
			signer, err := ssh.ParsePrivateKey(key)
			if err == nil {
				auth = append(auth, ssh.PublicKeys(signer))
			}
		}
	}

	// Fallback to password
	if cfg.Password != "" {
		auth = append(auth, ssh.Password(cfg.Password))
	}

	if len(auth) == 0 {
		return nil, fmt.Errorf("no auth method available")
	}

	config := &ssh.ClientConfig{
		User:            cfg.User,
		Auth:            auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}

	sshConnections[cfg.Host] = client
	return client, nil
}

func checkXrayStatus(client *ssh.Client) (bool, string) {
	// Check if Xray process is running
	session, err := client.NewSession()
	if err != nil {
		return false, "Failed to create session"
	}
	defer session.Close()

	cmd := "ps aux | grep -i xray | grep -v grep | wc -l"
	out, err := session.Output(cmd)
	if err != nil {
		return false, "Failed to check process"
	}

	if strings.TrimSpace(string(out)) != "0" {
		return true, ""
	}

	return false, "Xray process not running"
}

func getSystemStat(client *ssh.Client, statType string) int {
	session, err := client.NewSession()
	if err != nil {
		return -1
	}
	defer session.Close()

	var cmd string
	switch statType {
	case "cpu":
		cmd = "top -bn1 | grep \"Cpu(s)\" | awk '{print int($2)}'"
	case "memory":
		cmd = "free | grep Mem | awk '{print int($3/$2 * 100)}'"
	case "disk":
		cmd = "df / | tail -1 | awk '{print int($5)}'"
	}

	var out bytes.Buffer
	session.Stdout = &out
	err = session.Run(cmd)
	if err != nil {
		return -1
	}

	value, err := strconv.Atoi(strings.TrimSpace(out.String()))
	if err != nil {
		return -1
	}

	return value
}

func checkPing() string {
	addr := net.JoinHostPort(config.Network.PingHost, "53")
	conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
	if err != nil {
		return "offline"
	}
	defer conn.Close()
	return "online ✓"
}

func startMonitoring(chatID int64, interval int) {
	ticker := time.NewTicker(time.Duration(interval) * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if isMonitoring {
				sendServerStatus(chatID)
			}

		case <-stopMonitoring:
			return
		}
	}
}

func sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.DisableWebPagePreview = true
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}

func closeAllSSHConnections() {
	sshMutex.Lock()
	defer sshMutex.Unlock()

	for _, client := range sshConnections {
		client.Close()
	}
	sshConnections = make(map[string]*ssh.Client)
}
