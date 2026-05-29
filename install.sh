#!/bin/bash

################################################################################
# Server Monitor Bot - Автоматический скрипт установки
# Поддерживает: Ubuntu, Debian, CentOS, Rocky Linux
################################################################################

set -e

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Функции для вывода
print_header() {
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}ℹ️  $1${NC}"
}

print_header "🚀 Server Monitor Bot - Установка"

# Проверка прав администратора
if [[ $EUID -ne 0 ]]; then
    print_error "Этот скрипт должен запускаться от root"
    echo "Попробуйте: sudo bash install.sh"
    exit 1
fi

# Определение типа системы
if [ -f /etc/os-release ]; then
    . /etc/os-release
    OS=$ID
    VER=$VERSION_ID
else
    print_error "Не удалось определить тип ОС"
    exit 1
fi

print_info "Обнаружена ОС: $OS $VER"

# Установка зависимостей
print_header "📦 Установка зависимостей"

if [[ "$OS" == "ubuntu" || "$OS" == "debian" ]]; then
    print_info "Обновление пакетов..."
    apt-get update
    apt-get install -y curl wget git build-essential
    
    if ! command -v go &> /dev/null; then
        print_info "Установка Go..."
        wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz -O /tmp/go.tar.gz
        rm -rf /usr/local/go
        tar -C /usr/local -xzf /tmp/go.tar.gz
        rm /tmp/go.tar.gz
        echo "export PATH=\$PATH:/usr/local/go/bin" >> /etc/profile
        source /etc/profile
    fi
    
elif [[ "$OS" == "centos" || "$OS" == "rhel" || "$OS" == "rocky" ]]; then
    print_info "Обновление пакетов..."
    yum groupinstall -y "Development Tools"
    yum install -y curl wget git
    
    if ! command -v go &> /dev/null; then
        print_info "Установка Go..."
        wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz -O /tmp/go.tar.gz
        rm -rf /usr/local/go
        tar -C /usr/local -xzf /tmp/go.tar.gz
        rm /tmp/go.tar.gz
        echo "export PATH=\$PATH:/usr/local/go/bin" >> /etc/profile
        source /etc/profile
    fi
else
    print_error "Неподдерживаемая ОС: $OS"
    exit 1
fi

# Проверка Go
if ! command -v go &> /dev/null; then
    print_error "Go не был установлен правильно"
    exit 1
fi

print_success "Go установлен: $(go version)"

# Создание директории приложения
print_header "📂 Подготовка директории"

APP_DIR="/opt/server-monitor"
if [ -d "$APP_DIR" ]; then
    print_info "Директория $APP_DIR уже существует"
    read -p "Перезаписать? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        rm -rf "$APP_DIR"
    else
        print_info "Использую существующую директорию"
    fi
fi

mkdir -p "$APP_DIR"
cd "$APP_DIR"
print_success "Директория: $APP_DIR"

# Загрузка файлов с GitHub
print_header "📥 Загрузка файлов с GitHub"

if [ -d ".git" ]; then
    print_info "Репозиторий уже клонирован, обновляю..."
    git pull origin main 2>/dev/null || true
else
    print_info "Клонирую репозиторий..."
    git clone https://github.com/AKHATOV8/Test.git . 2>/dev/null || {
        print_error "Не удалось клонировать репозиторий"
        print_info "Скачиваю файлы напрямую..."
        wget -q https://raw.githubusercontent.com/AKHATOV8/Test/main/main.go -O main.go
        wget -q https://raw.githubusercontent.com/AKHATOV8/Test/main/go.mod -O go.mod
        wget -q https://raw.githubusercontent.com/AKHATOV8/Test/main/config.example.yaml -O config.example.yaml
    }
fi

print_success "Файлы загружены"

# Установка зависимостей Go
print_header "🔨 Компиляция бота"

print_info "Загрузка Go модулей..."
/usr/local/go/bin/go mod download
/usr/local/go/bin/go mod tidy

print_info "Компиляция..."
/usr/local/go/bin/go build -o server-monitor main.go 2>&1

if [ ! -f "server-monitor" ]; then
    print_error "Компиляция не удалась"
    exit 1
fi

print_success "Бот скомпилирован: $(ls -lh server-monitor | awk '{print $5}')"

# Создание конфига
print_header "⚙️  Конфигурация"

if [ ! -f "config.yaml" ]; then
    if [ -f "config.example.yaml" ]; then
        cp config.example.yaml config.yaml
        print_success "Создан config.yaml"
    else
        print_error "config.example.yaml не найден"
        exit 1
    fi
else
    print_info "config.yaml уже существует"
fi

# Установка systemd сервиса
print_header "🔧 Установка systemd сервиса"

cat > /etc/systemd/system/server-monitor.service << 'EOF'
[Unit]
Description=Server Monitor Telegram Bot
After=network.target
Wants=network-online.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/server-monitor
ExecStart=/opt/server-monitor/server-monitor
Restart=on-failure
RestartSec=10
StandardOutput=journal
StandardError=journal
SyslogIdentifier=server-monitor

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
print_success "Сервис зарегистрирован"

# Установка прав доступа
chmod +x /opt/server-monitor/server-monitor
chmod 600 /opt/server-monitor/config.yaml

# Итоговая информация
print_header "✅ Установка завершена!"

echo ""
echo -e "${BLUE}📝 СЛЕДУЮЩИЕ ШАГИ:${NC}"
echo ""
echo "1️⃣  Отредактируйте конфигурацию:"
echo -e "   ${YELLOW}nano /opt/server-monitor/config.yaml${NC}"
echo ""
echo "2️⃣  Необходимо указать:"
echo "   • telegram.token - получите от @BotFather в Telegram"
echo "   • telegram.chat_id - получите от @userinfobot в Telegram"
echo "   • servers - данные ваших серверов (до 15)"
echo ""
echo "3️⃣  Запустите бота:"
echo -e "   ${YELLOW}sudo systemctl start server-monitor${NC}"
echo ""
echo "4️⃣  Разрешите автозапуск:"
echo -e "   ${YELLOW}sudo systemctl enable server-monitor${NC}"
echo ""
echo -e "${BLUE}📊 Проверка статуса:${NC}"
echo -e "   ${YELLOW}sudo systemctl status server-monitor${NC}"
echo ""
echo -e "${BLUE}📋 Просмотр логов:${NC}"
echo -e "   ${YELLOW}sudo journalctl -u server-monitor -f${NC}"
echo ""
echo -e "${BLUE}🔗 Ссылки:${NC}"
echo "   • Репозиторий: https://github.com/AKHATOV8/Test"
echo "   • Получить token: https://t.me/BotFather"
echo "   • Узнать Chat ID: https://t.me/userinfobot"
echo ""
