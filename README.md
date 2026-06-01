# 🤖 Server Monitor Bot - Telegram Bot для мониторинга Linux серверов

Телеграм-бот для мониторинга Linux серверов с 3x-ui и Xray. Отслеживайте статус серверов, CPU, RAM, диск и скорость интернета через Telegram.

## ✨ Функции

- 🖥️ **Мониторинг Xray** - проверка статуса процесса
- 📊 **Система ресурсов** - CPU, RAM, диск в реальном времени
- 🌐 **Проверка интернета** - ping и скорость соединения
- ⚠️ **Умные алерты** - оповещения при превышении порогов
- ⏰ **Гибкий мониторинг** - 5, 10, 15, 30 минут на выбор
- 📈 **До 15 серверов** - легко масштабируется
- 🔒 **SSH подключение** - безопасное соединение с серверами
- 🚀 **Автозапуск** - systemd сервис для автоматического запуска

## 🚀 Быстрая установка (3 команды)

### Способ 1: Автоматическая установка (Рекомендуется)

```bash
# Загрузить и запустить скрипт установки
curl -fsSL https://raw.githubusercontent.com/AKHATOV8/Test/main/install.sh | sudo bash
```

### Способ 2: Ручная установка

```bash
# 1. Клонировать репозиторий
git clone https://github.com/AKHATOV8/Test.git /tmp/server-monitor
cd /tmp/server-monitor

# 2. Запустить скрипт установки
chmod +x install.sh
sudo bash install.sh
```

## ⚙️ Конфигурация

После установки отредактируйте конфиг:

```bash
sudo nano /opt/server-monitor/config.yaml
```

### Пример конфигурации

```yaml
# Telegram настройки
telegram:
  token: "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"  # От @BotFather
  chat_id: 987654321                                    # От @userinfobot

# Интервалы проверки (минуты)
monitor:
  check_intervals: [5, 10, 15, 30]

# Серверы для мониторинга (макс 15)
servers:
  - name: "Server 1"
    host: "192.168.1.100"
    port: 22
    user: "root"
    ssh_key: "/root/.ssh/id_rsa"  # Путь к приватному ключу
    # password: ""                  # ИЛИ используйте пароль

  - name: "Server 2"
    host: "192.168.1.101"
    port: 22
    user: "root"
    ssh_key: "/root/.ssh/id_rsa"

  - name: "Server 3"
    host: "192.168.1.102"
    port: 22
    user: "root"
    ssh_key: "/root/.ssh/id_rsa"

  - name: "Server 4"
    host: "192.168.1.103"
    port: 22
    user: "root"
    ssh_key: "/root/.ssh/id_rsa"

# Xray настройки
xray:
  port: 10085
  check_process: true

# Сетевые настройки
network:
  ping_host: "8.8.8.8"
  ping_count: 4

# Пороги оповещений
alerts:
  cpu_threshold: 80      # % CPU
  memory_threshold: 85   # % RAM
  disk_threshold: 90     # % диска
```

## 🎮 Команды Telegram бота

| Команда | Описание |
|---------|---------|
| `/start` | Приветствие и справка |
| `/status` | Мгновенная проверка статуса всех серверов |
| `/servers` | Список всех добавленных серверов |
| `/monitor 5` | Автоматический мониторинг каждые 5 минут |
| `/monitor 10` | Автоматический мониторинг каждые 10 минут |
| `/monitor 15` | Автоматический мониторинг каждые 15 минут |
| `/monitor 30` | Автоматический мониторинг каждые 30 минут |
| `/stop` | Остановить автоматический мониторинг |
| `/help` | Справка по командам |

## 📋 Получение необходимых данных

### 1. Получить Telegram Bot Token

1. Откройте Telegram
2. Найдите **@BotFather**
3. Напишите `/newbot`
4. Следуйте инструкциям
5. Скопируйте полученный token

### 2. Получить Chat ID

**Способ 1:**
1. Найдите **@userinfobot** в Telegram
2. Напишите `/start`
3. Скопируйте ваш User ID

**Способ 2:**
1. Напишите боту любое сообщение
2. Откройте: `https://api.telegram.org/bot<YOUR_TOKEN>/getUpdates`
3. Найдите `"chat":{"id":123456789}`

### 3. Подготовить SSH доступ

#### Способ 1: SSH ключи (рекомендуется)

```bash
# На машине, где установлен бот:
ssh-keygen -t rsa -N "" -f ~/.ssh/id_rsa

# Скопируйте ключ на сервер:
ssh-copy-id -i ~/.ssh/id_rsa.pub root@server_ip
```

#### Способ 2: Пароль

В конфиге просто используйте `password` вместо `ssh_key`:

```yaml
servers:
  - name: "Server 1"
    host: "192.168.1.100"
    port: 22
    user: "root"
    password: "your_password"
```

## 🔧 Управление ботом

### Запуск

```bash
# Запустить бота
sudo systemctl start server-monitor

# Разрешить автозапуск при перезагрузке
sudo systemctl enable server-monitor
```

### Проверка статуса

```bash
# Статус сервиса
sudo systemctl status server-monitor

# Просмотр логов в реальном времени
sudo journalctl -u server-monitor -f

# Просмотр последних 50 строк логов
sudo journalctl -u server-monitor -n 50
```

### Остановка

```bash
sudo systemctl stop server-monitor
```

### Перезагрузка

```bash
sudo systemctl restart server-monitor
```

## ➕ Добавление нового сервера

1. Отредактируйте конфиг:
```bash
sudo nano /opt/server-monitor/config.yaml
```

2. Добавьте новый сервер в секцию `servers`:
```yaml
- name: "Server 5"
  host: "192.168.1.105"
  port: 22
  user: "root"
  ssh_key: "/root/.ssh/id_rsa"
```

3. Перезагрузите бота:
```bash
sudo systemctl restart server-monitor
```

⚠️ **Максимум 15 серверов**

## 🐛 Решение проблем

### "Failed to connect to Telegram"
- Проверьте token в config.yaml
- Убедитесь, что интернет есть
- Проверьте firewall

### "SSH connection failed"
- Проверьте IP адрес сервера
- Проверьте SSH доступ вручную: `ssh -i key.pem user@host`
- Убедитесь, что правильно указаны user и port
- Проверьте права доступа на ключ: `chmod 600 ~/.ssh/id_rsa`

### "Xray process not running"
- 3x-ui может быть остановлена
- Проверьте вручную: `ssh root@server 'ps aux | grep xray'`
- Перезагрузите 3x-ui на сервере

### Просмотр ошибок
```bash
sudo journalctl -u server-monitor -f --lines=100
```

## 📦 Требования

- **Linux** сервер (Ubuntu, Debian, CentOS, Rocky Linux)
- **Go 1.16+** (устанавливается автоматически)
- **Git** (устанавливается автоматически)
- **SSH доступ** к мониторимым серверам
- **Telegram Bot Token** (от @BotFather)

## 🔄 Обновление бота

```bash
cd /opt/server-monitor
sudo git pull
sudo go build -o server-monitor main.go
sudo systemctl restart server-monitor
```

## 📝 Логирование

Все логи сохраняются в systemd журнал:

```bash
# Просмотр логов
sudo journalctl -u server-monitor

# Последние 100 строк
sudo journalctl -u server-monitor -n 100

# В реальном времени
sudo journalctl -u server-monitor -f

# За последний час
sudo journalctl -u server-monitor --since "1 hour ago"
```

## 🗂️ Структура файлов

```
/opt/server-monitor/
├── server-monitor        # Скомпилированный бот
├── config.yaml          # Конфигурация (заполняется вручную)
├── config.example.yaml  # Пример конфигурации
├── main.go             # Исходный код
├── go.mod              # Зависимости
└── README.md           # Документация
```

## 💡 Советы

1. **SSH ключи безопаснее** чем пароли - используйте их
2. **Тестируйте SSH подключение** перед настройкой бота
3. **Проверяйте логи** при любых проблемах
4. **Используйте мониторинг каждые 30 минут** для экономии ресурсов
5. **Добавьте бота в избранное** в Telegram для быстрого доступа

## 📚 Дополнительные ссылки

- [Telegram Bot API](https://core.telegram.org/bots/api)
- [Go Programming Language](https://golang.org/)
- [SSH на Linux](https://linux.die.net/man/1/ssh)
- [Systemd](https://www.freedesktop.org/wiki/Software/systemd/)

## 📞 Поддержка

Если у вас есть вопросы или проблемы:
- Создайте Issue на GitHub: [AKHATOV8/Test/issues](https://github.com/AKHATOV8/Test/issues)
- Проверьте документацию: [INSTALL.md](INSTALL.md)

## 📄 Лицензия

MIT License - смотрите LICENSE файл для деталей

---

**Последнее обновление:** 2026-06-01  
**Версия:** 1.0.0  
**Автор:** AKHATOV8
