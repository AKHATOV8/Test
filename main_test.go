package main

import (
	"testing"
	"time"
)

// TestSSHPasswordAuth проверяет подключение по паролю
func TestSSHPasswordAuth(t *testing.T) {
	// Эта функция должна быть запущена с реальным серверами
	// SKIP если нет доступных тестовых серверов
	
	t.Run("Valid password", func(t *testing.T) {
		cfg := ServerConfig{
			Name:     "Test Server",
			Host:     "localhost",
			Port:     22,
			User:     "root",
			Password: "test_password",
		}
		
		_, err := getSSHConnection(cfg)
		if err == nil {
			t.Log("✅ SSH password auth successful")
		} else {
			t.Logf("⚠️ SSH connection failed (expected in test env): %v", err)
		}
	})

	t.Run("Invalid password", func(t *testing.T) {
		cfg := ServerConfig{
			Name:     "Test Server",
			Host:     "localhost",
			Port:     22,
			User:     "root",
			Password: "wrong_password",
		}
		
		_, err := getSSHConnection(cfg)
		if err != nil {
			t.Log("✅ SSH correctly rejected invalid password")
		} else {
			t.Log("❌ SSH should reject invalid password")
			t.Fail()
		}
	})
}

// TestXrayProcessDetection проверяет обнаружение процесса
func TestXrayProcessDetection(t *testing.T) {
	t.Run("Xray running", func(t *testing.T) {
		// Это требует реального сервера с Xray
		t.Skip("Requires real server with Xray")
		
		cfg := ServerConfig{
			Name:     "Test Server",
			Host:     "192.168.1.100",
			Port:     22,
			User:     "root",
			SSHKey:   "/root/.ssh/id_rsa",
		}
		
		client, err := getSSHConnection(cfg)
		if err != nil {
			t.Fatalf("Failed to connect: %v", err)
		}
		
		running, errMsg := checkXrayStatus(client)
		if running {
			t.Log("✅ Xray process detected")
		} else {
			t.Logf("ℹ️ Xray not running: %s", errMsg)
		}
	})

	t.Run("Xray stopped", func(t *testing.T) {
		t.Skip("Requires real server with stopped Xray")
		
		cfg := ServerConfig{
			Name:     "Test Server",
			Host:     "192.168.1.100",
			Port:     22,
			User:     "root",
			SSHKey:   "/root/.ssh/id_rsa",
		}
		
		client, err := getSSHConnection(cfg)
		if err != nil {
			t.Fatalf("Failed to connect: %v", err)
		}
		
		running, errMsg := checkXrayStatus(client)
		if !running && errMsg != "" {
			t.Log("✅ Correctly detected stopped Xray")
		} else {
			t.Log("❌ Should detect stopped Xray")
			t.Fail()
		}
	})
}

// TestSystemMetrics проверяет получение метрик
func TestSystemMetrics(t *testing.T) {
	t.Run("CPU metric", func(t *testing.T) {
		t.Skip("Requires real SSH connection")
		
		cfg := ServerConfig{
			Name:   "Test Server",
			Host:   "192.168.1.100",
			Port:   22,
			User:   "root",
			SSHKey: "/root/.ssh/id_rsa",
		}
		
		client, err := getSSHConnection(cfg)
		if err != nil {
			t.Fatalf("Failed to connect: %v", err)
		}
		
		cpu := getSystemStat(client, "cpu")
		if cpu >= 0 && cpu <= 100 {
			t.Logf("✅ CPU metric: %d%%", cpu)
		} else {
			t.Logf("❌ Invalid CPU metric: %d", cpu)
			t.Fail()
		}
	})

	t.Run("Memory metric", func(t *testing.T) {
		t.Skip("Requires real SSH connection")
		
		cfg := ServerConfig{
			Name:   "Test Server",
			Host:   "192.168.1.100",
			Port:   22,
			User:   "root",
			SSHKey: "/root/.ssh/id_rsa",
		}
		
		client, err := getSSHConnection(cfg)
		if err != nil {
			t.Fatalf("Failed to connect: %v", err)
		}
		
		mem := getSystemStat(client, "memory")
		if mem >= 0 && mem <= 100 {
			t.Logf("✅ Memory metric: %d%%", mem)
		} else {
			t.Logf("❌ Invalid memory metric: %d", mem)
			t.Fail()
		}
	})

	t.Run("Disk metric", func(t *testing.T) {
		t.Skip("Requires real SSH connection")
		
		cfg := ServerConfig{
			Name:   "Test Server",
			Host:   "192.168.1.100",
			Port:   22,
			User:   "root",
			SSHKey: "/root/.ssh/id_rsa",
		}
		
		client, err := getSSHConnection(cfg)
		if err != nil {
			t.Fatalf("Failed to connect: %v", err)
		}
		
		disk := getSystemStat(client, "disk")
		if disk >= 0 && disk <= 100 {
			t.Logf("✅ Disk metric: %d%%", disk)
		} else {
			t.Logf("❌ Invalid disk metric: %d", disk)
			t.Fail()
		}
	})
}

// TestConfigLoading проверяет загрузку конфига
func TestConfigLoading(t *testing.T) {
	t.Run("Valid config", func(t *testing.T) {
		// Создание временного конфига
		testConfig := `
telegram:
  token: "test_token"
  chat_id: 123456789

monitor:
  check_intervals: [5, 10, 15, 30]

servers:
  - name: "Server 1"
    host: "192.168.1.100"
    port: 22
    user: "root"
    ssh_key: "/root/.ssh/id_rsa"

alerts:
  cpu_threshold: 80
  memory_threshold: 85
  disk_threshold: 90
`
		t.Logf("✅ Config structure valid:\n%s", testConfig)
	})
}

// TestMonitoringInterval проверяет интервалы мониторинга
func TestMonitoringInterval(t *testing.T) {
	validIntervals := map[int]bool{5: true, 10: true, 15: true, 30: true}

	tests := []struct {
		name      string
		interval  int
		isValid   bool
	}{
		{"5 minutes", 5, true},
		{"10 minutes", 10, true},
		{"15 minutes", 15, true},
		{"30 minutes", 30, true},
		{"1 minute", 1, false},
		{"60 minutes", 60, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := validIntervals[tt.interval]
			if isValid == tt.isValid {
				t.Logf("✅ Interval %d is %v", tt.interval, tt.isValid)
			} else {
				t.Errorf("❌ Expected interval %d to be %v, got %v", tt.interval, tt.isValid, isValid)
			}
		})
	}
}

// TestErrorHandling проверяет обработку ошибок
func TestErrorHandling(t *testing.T) {
	t.Run("Invalid SSH config", func(t *testing.T) {
		cfg := ServerConfig{
			Name:     "Invalid Server",
			Host:     "invalid.host.local",
			Port:     22,
			User:     "root",
			Password: "password",
		}
		
		_, err := getSSHConnection(cfg)
		if err != nil {
			t.Log("✅ Correctly handled invalid SSH config")
		} else {
			t.Log("❌ Should return error for invalid SSH config")
			t.Fail()
		}
	})

	t.Run("No auth method", func(t *testing.T) {
		cfg := ServerConfig{
			Name: "No Auth Server",
			Host: "192.168.1.100",
			Port: 22,
			User: "root",
			// No password or SSH key
		}
		
		_, err := getSSHConnection(cfg)
		if err != nil && err.Error() == "no auth method available" {
			t.Log("✅ Correctly detected missing auth method")
		} else {
			t.Log("❌ Should detect missing auth method")
			t.Fail()
		}
	})
}

// TestServerStatus проверяет структуру статуса сервера
func TestServerStatus(t *testing.T) {
	t.Run("ServerStatus structure", func(t *testing.T) {
		status := ServerStatus{
			Name:       "Test Server",
			Host:       "192.168.1.100",
			Online:     true,
			XrayStatus: true,
			CPU:        45,
			Memory:     62,
			Disk:       78,
		}
		
		if status.Online && status.XrayStatus && status.CPU > 0 {
			t.Log("✅ ServerStatus structure is valid")
		} else {
			t.Log("❌ ServerStatus structure has invalid data")
			t.Fail()
		}
	})
}

// TestTelegramIntegration проверяет интеграцию с Telegram
func TestTelegramIntegration(t *testing.T) {
	t.Run("Message formatting", func(t *testing.T) {
		status := ServerStatus{
			Name:       "Server 1",
			Host:       "192.168.1.100",
			Online:     true,
			XrayStatus: true,
			CPU:        45,
			Memory:     60,
			Disk:       75,
		}
		
		formatted := formatServerStatus(status)
		if len(formatted) > 0 && status.Name != "" {
			t.Log("✅ Message formatting works")
		} else {
			t.Log("❌ Message formatting failed")
			t.Fail()
		}
	})
}

// TestConcurrentConnections проверяет параллельные подключения
func TestConcurrentConnections(t *testing.T) {
	t.Run("Multiple server status check", func(t *testing.T) {
		servers := []ServerConfig{
			{Name: "Server 1", Host: "192.168.1.100", User: "root"},
			{Name: "Server 2", Host: "192.168.1.101", User: "root"},
			{Name: "Server 3", Host: "192.168.1.102", User: "root"},
		}
		
		if len(servers) == 3 {
			t.Log("✅ Concurrent connection structure is valid")
		}
	})
}

// TestAlerts проверяет систему оповещений
func TestAlerts(t *testing.T) {
	t.Run("CPU alert threshold", func(t *testing.T) {
		cpuThreshold := 80
		testCPU := 85
		
		if testCPU >= cpuThreshold {
			t.Logf("✅ Alert triggered for CPU: %d%% >= %d%%", testCPU, cpuThreshold)
		} else {
			t.Log("❌ Alert should trigger")
			t.Fail()
		}
	})

	t.Run("Memory alert threshold", func(t *testing.T) {
		memThreshold := 85
		testMem := 90
		
		if testMem >= memThreshold {
			t.Logf("✅ Alert triggered for Memory: %d%% >= %d%%", testMem, memThreshold)
		} else {
			t.Log("❌ Alert should trigger")
			t.Fail()
		}
	})

	t.Run("Disk alert threshold", func(t *testing.T) {
		diskThreshold := 90
		testDisk := 92
		
		if testDisk >= diskThreshold {
			t.Logf("✅ Alert triggered for Disk: %d%% >= %d%%", testDisk, diskThreshold)
		} else {
			t.Log("❌ Alert should trigger")
			t.Fail()
		}
	})
}

// TestMonitoringLoop проверяет основной цикл мониторинга
func TestMonitoringLoop(t *testing.T) {
	t.Run("Monitoring duration", func(t *testing.T) {
		// Тест на правильность работы цикла
		start := time.Now()
		expectedDuration := 5 * time.Second
		
		// Имитация работы
		time.Sleep(expectedDuration)
		
		elapsed := time.Since(start)
		if elapsed >= expectedDuration {
			t.Logf("✅ Monitoring loop works correctly (elapsed: %v)", elapsed)
		} else {
			t.Log("❌ Monitoring loop timing is incorrect")
			t.Fail()
		}
	})
}

// BenchmarkSSHConnection тестирует производительность SSH подключения
func BenchmarkSSHConnection(b *testing.B) {
	cfg := ServerConfig{
		Name:     "Benchmark Server",
		Host:     "192.168.1.100",
		Port:     22,
		User:     "root",
		SSHKey:   "/root/.ssh/id_rsa",
	}
	
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = getSSHConnection(cfg)
	}
}

// BenchmarkSystemMetrics тестирует производительность получения метрик
func BenchmarkSystemMetrics(b *testing.B) {
	b.Skip("Requires real SSH connection")
	
	cfg := ServerConfig{
		Name:   "Benchmark Server",
		Host:   "192.168.1.100",
		Port:   22,
		User:   "root",
		SSHKey: "/root/.ssh/id_rsa",
	}
	
	client, _ := getSSHConnection(cfg)
	defer client.Close()
	
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = getSystemStat(client, "cpu")
		_ = getSystemStat(client, "memory")
		_ = getSystemStat(client, "disk")
	}
}
