# 🚀 Установка Server Monitor Bot

## Требования
- Linux сервер (Ubuntu, Debian, CentOS)
- Go 1.16+ (если не установлен, скрипт установит)
- Git
- Telegram Bot Token (получить у @BotFather)
- SSH доступ к серверам для мониторинга

## Быстрая установка (Recommended)

### 1. Клонируйте репозиторий
```bash
cd /tmp
git clone https://github.com/AKHATOV8/Test.git server-monitor
cd server-monitor
```

### 2. Запустите скрипт установки
```bash
chmod +x install.sh
sudo ./install.sh
```

### 3. Отредактируйте конфиг
```bash
sudo nano /opt/server-monitor/config.yaml
```

**Обязательно заполните:**
- `telegram.token` - ваш Telegram bot token
- `telegram.chat_id` - ваш chat ID
- `servers` - данные серверов для мониторинга

### 4. Перезагрузите systemd
```bash
sudo systemctl daemon-reload
sudo systemctl start server-monitor
sudo systemctl enable server-monitor
```

## Проверка работы

```bash
# Проверить статус бота
sudo systemctl status server-monitor

# Посмотреть логи
sudo journalctl -u server-monitor -f
```

## Ручная установка (если скрипт не работает)

### 1. Установите Go
```bash
sudo apt update
sudo apt install -y golang-go git
```

### 2. Создайте директорию приложения
```bash
sudo mkdir -p /opt/server-monitor
sudo cd /opt/server-monitor
```

### 3. Клонируйте репозиторий
```bash
sudo git clone https://github.com/AKHATOV8/Test.git .
```

### 4. Скомпилируйте
```bash
sudo go mod download
sudo go mod tidy
sudo go build -o server-monitor main.go
```

### 5. Скопируйте конфиг
```bash
sudo cp config.example.yaml config.yaml
sudo nano config.yaml  # отредактируйте
```

### 6. Создайте systemd сервис
```bash
sudo cp server-monitor.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl start server-monitor
sudo systemctl enable server-monitor
```

## Конфигурация

### config.yaml структура

```yaml
telegram:
  token: "YOUR_BOT_TOKEN"        # От @BotFather
  chat_id: 123456789             # Ваш Chat ID

monitor:
  check_intervals: [5, 10, 15, 30]  # Интервалы в минутах

servers:
  - name: "Server 1"
    host: "192.168.1.100"
    port: 22
    user: "root"
    ssh_key: "/root/.ssh/id_rsa"  # SSH ключ
    # ИЛИ password: "password"     # Или пароль
    
  - name: "Server 2"
    host: "192.168.1.101"
    port: 22
    user: "root"
    ssh_key: "/root/.ssh/id_rsa"

xray:
  port: 10085
  check_process: true

network:
  ping_host: "8.8.8.8"
  ping_count: 4

alerts:
  cpu_threshold: 80      # % CPU для оповещения
  memory_threshold: 85   # % RAM для оповещения
  disk_threshold: 90     # % Диска для оповещения
```

## Команды Telegram бота

- `/start` - Приветствие и справка
- `/status` - Мгновенная проверка всех серверов
- `/servers` - Список всех серверов
- `/monitor 5` - Автоматический мониторинг каждые 5 минут
- `/monitor 10` - Каждые 10 минут
- `/monitor 15` - Каждые 15 минут
- `/monitor 30` - Каждые 30 минут
- `/stop` - Остановить автоматический мониторинг
- `/help` - Справка

## Добавление серверов

### Максимум 15 серверов

1. Отредактируйте `/opt/server-monitor/config.yaml`
2. Добавьте новый сервер в секцию `servers`
3. Перезагрузите бота:
   ```bash
   sudo systemctl restart server-monitor
   ```

## Получение Telegram Token

1. Откройте Telegram
2. Найдите @BotFather
3. Напишите `/newbot`
4. Следуйте инструкциям
5. Скопируйте полученный token

## Получение Chat ID

1. Найдите @userinfobot в Telegram
2. Напишите `/start`
3. Скопируйте ваш User ID

ИЛИ:

1. Напишите боту любое сообщение
2. Перейдите: `https://api.telegram.org/botYOUR_TOKEN/getUpdates`
3. Найдите `"chat":{"id":123456789}`

## SSH Ключи

### Если используете SSH ключи:

```bash
# На хосте бота
ssh-keygen -t rsa -N ""

# Скопируйте ключ на сервер
ssh-copy-id -i ~/.ssh/id_rsa.pub root@server_ip
```

## Проблемы и решение

### "Failed to connect to Telegram"
- Проверьте token в config.yaml
- Убедитесь, что интернет есть

### "SSH connection failed"
- Проверьте IP адрес сервера
- Проверьте доступ по SSH вручную: `ssh -i key.pem user@host`
- Убедитесь, что пользователь и ключ правильные

### "Xray process not running"
- 3x-ui может быть остановлена
- Проверьте вручную: `ssh root@server 'ps aux | grep xray'`

## Обновление бота

```bash
cd /opt/server-monitor
sudo git pull
sudo go build -o server-monitor main.go
sudo systemctl restart server-monitor
```

## Поддержка

Если у вас есть вопросы, создайте issue в репозитории:
https://github.com/AKHATOV8/Test/issues
