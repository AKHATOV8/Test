#!/bin/bash

# Server Monitor Bot - Installation Script
# Скрипт для автоматической установки и запуска бота на Linux сервере

set -e

echo "🚀 Server Monitor Bot - Установка"
echo "=================================="

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go не установлен"
    echo "📥 Устанавливаю Go..."
    
    # Download and install Go
    cd /tmp
    wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
    rm go1.21.0.linux-amd64.tar.gz
    
    echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc
    source ~/.bashrc
fi

echo "✅ Go версия: $(go version)"

# Create app directory
APP_DIR="/opt/server-monitor"
mkdir -p $APP_DIR
cd $APP_DIR

echo "📂 Рабочая директория: $APP_DIR"

# Download bot from GitHub
echo "📥 Загружаю файлы бота..."
git clone https://github.com/AKHATOV8/Test.git . 2>/dev/null || git -C . pull

# Install dependencies
echo "📦 Устанавливаю зависимости..."
go mod download
go mod tidy

# Build the bot
echo "🔨 Компилирую бота..."
go build -o server-monitor main.go

# Copy config if not exists
if [ ! -f config.yaml ]; then
    cp config.example.yaml config.yaml
    echo "⚙️ Создан config.yaml - отредактируйте его!"
    echo "   Файл: $APP_DIR/config.yaml"
fi

echo ""
echo "✅ Установка завершена!"
echo ""
echo "📝 Следующие шаги:"
echo "1. Отредактируйте конфиг: nano $APP_DIR/config.yaml"
echo "2. Добавьте Telegram token и Chat ID"
echo "3. Добавьте данные серверов (до 15)"
echo ""
echo "🚀 Запуск бота:"
echo "   cd $APP_DIR && ./server-monitor"
echo ""
echo "⚙️ Для автоматического запуска создан systemd сервис"
echo "   sudo systemctl start server-monitor"
echo "   sudo systemctl enable server-monitor"
echo ""
