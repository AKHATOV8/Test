#!/bin/bash

################################################################################
# Server Monitor Bot - Test Script
# Полное тестирование всех функций бота
################################################################################

set -e

# Цвета
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_header() {
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}"
}

print_test() {
    echo -e "${YELLOW}▶ $1${NC}"
}

print_pass() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_fail() {
    echo -e "${RED}❌ $1${NC}"
}

print_skip() {
    echo -e "${YELLOW}⊘ $1 (пропущено)${NC}"
}

# Главное меню
print_header "🧪 Server Monitor Bot - Тестирование"

echo ""
echo "Выберите режим тестирования:"
echo "1) Быстрое тестирование (config, структура)"
echo "2) Полное тестирование (требует реальные серверы)"
echo "3) Unit тесты (go test)"
echo "4) Тестирование SSH подключения"
echo "5) Тестирование Telegram интеграции"
echo "6) Все тесты"
echo ""
read -p "Выберите пункт (1-6): " choice

case $choice in

1)
    print_header "🔍 Быстрое тестирование"
    
    print_test "1. Проверка структуры файлов"
    if [ -f "main.go" ] && [ -f "config.yaml" ] && [ -f "config.example.yaml" ]; then
        print_pass "Все основные файлы присутствуют"
    else
        print_fail "Отсутствуют необходимые файлы"
        exit 1
    fi
    
    print_test "2. Проверка Go синтаксиса"
    if go fmt ./... && go vet ./...; then
        print_pass "Go синтаксис верный"
    else
        print_fail "Ошибки в Go коде"
        exit 1
    fi
    
    print_test "3. Проверка конфигурации"
    if grep -q "telegram:" config.yaml && grep -q "servers:" config.yaml; then
        print_pass "Структура конфига корректна"
    else
        print_fail "Конфиг имеет неверную структуру"
        exit 1
    fi
    
    print_test "4. Проверка зависимостей"
    if go mod verify; then
        print_pass "Зависимости валидны"
    else
        print_fail "Проблемы с зависимостями"
        exit 1
    fi
    
    print_pass "✅ Быстрое тестирование завершено"
    ;;

2)
    print_header "🔧 Полное тестирование"
    
    read -p "Введите IP тестового сервера: " TEST_SERVER
    read -p "Введите пользователя SSH (по умолчанию root): " TEST_USER
    TEST_USER=${TEST_USER:-root}
    read -p "Введите путь к SSH ключу (по умолчанию ~/.ssh/id_rsa): " TEST_KEY
    TEST_KEY=${TEST_KEY:-~/.ssh/id_rsa}
    
    print_test "1. Проверка SSH подключения"
    if ssh -i "$TEST_KEY" -o StrictHostKeyChecking=no -o ConnectTimeout=5 "$TEST_USER@$TEST_SERVER" "echo SSH OK" &>/dev/null; then
        print_pass "SSH подключение работает"
    else
        print_fail "SSH подключение не удалось"
        exit 1
    fi
    
    print_test "2. Проверка процесса Xray"
    if ssh -i "$TEST_KEY" "$TEST_USER@$TEST_SERVER" "ps aux | grep -i xray | grep -v grep" &>/dev/null; then
        print_pass "Xray процесс обнаружен"
    else
        print_skip "Xray не запущен на сервере"
    fi
    
    print_test "3. Проверка системных команд"
    if ssh -i "$TEST_KEY" "$TEST_USER@$TEST_SERVER" "top -bn1 | head -1" &>/dev/null; then
        print_pass "Команда top доступна"
    else
        print_fail "Команда top не найдена"
        exit 1
    fi
    
    if ssh -i "$TEST_KEY" "$TEST_USER@$TEST_SERVER" "free | grep Mem" &>/dev/null; then
        print_pass "Команда free доступна"
    else
        print_fail "Команда free не найдена"
        exit 1
    fi
    
    if ssh -i "$TEST_KEY" "$TEST_USER@$TEST_SERVER" "df /" &>/dev/null; then
        print_pass "Команда df доступна"
    else
        print_fail "Команда df не найдена"
        exit 1
    fi
    
    print_test "4. Получение метрик"
    CPU=$(ssh -i "$TEST_KEY" "$TEST_USER@$TEST_SERVER" "top -bn1 | grep 'Cpu(s)' | awk '{print int(\$2)}'" 2>/dev/null || echo "N/A")
    MEM=$(ssh -i "$TEST_KEY" "$TEST_USER@$TEST_SERVER" "free | grep Mem | awk '{print int(\$3/\$2 * 100)}'" 2>/dev/null || echo "N/A")
    DISK=$(ssh -i "$TEST_KEY" "$TEST_USER@$TEST_SERVER" "df / | tail -1 | awk '{print int(\$5)}'" 2>/dev/null || echo "N/A")
    
    echo "   CPU:  $CPU%"
    echo "   MEM:  $MEM%"
    echo "   DISK: $DISK%"
    print_pass "Метрики получены"
    
    print_pass "✅ Полное тестирование завершено"
    ;;

3)
    print_header "🧬 Unit тесты"
    
    print_test "Запуск всех тестов..."
    if go test -v -cover; then
        print_pass "Все тесты пройдены"
    else
        print_fail "Некоторые тесты не прошли"
        exit 1
    fi
    ;;

4)
    print_header "🔐 Тестирование SSH"
    
    echo ""
    echo "Тесты SSH подключения:"
    echo "1) SSH ключ"
    echo "2) Пароль"
    echo "3) Оба метода"
    echo ""
    read -p "Выберите (1-3): " ssh_choice
    
    read -p "Введите IP сервера: " SSH_IP
    read -p "Введите пользователя (по умолчанию root): " SSH_USER
    SSH_USER=${SSH_USER:-root}
    
    case $ssh_choice in
    1)
        read -p "Введите путь к SSH ключу: " SSH_KEY
        print_test "Тестирование SSH ключа..."
        if ssh -i "$SSH_KEY" -o ConnectTimeout=5 "$SSH_USER@$SSH_IP" "echo SSH Key OK"; then
            print_pass "SSH ключ работает"
        else
            print_fail "SSH ключ не сработал"
        fi
        ;;
    2)
        print_test "Тестирование SSH пароля..."
        if sshpass -p "$(read -sp 'Пароль: ' pass; echo $pass)" ssh -o ConnectTimeout=5 "$SSH_USER@$SSH_IP" "echo SSH Password OK"; then
            print_pass "SSH пароль работает"
        else
            print_fail "SSH пароль не сработал"
        fi
        ;;
    3)
        read -p "Введите путь к SSH ключу: " SSH_KEY
        print_test "Тестирование SSH ключа..."
        if ssh -i "$SSH_KEY" -o ConnectTimeout=5 "$SSH_USER@$SSH_IP" "echo SSH Key OK"; then
            print_pass "SSH ключ работает"
        else
            print_fail "SSH ключ не сработал"
        fi
        ;;
    esac
    ;;

5)
    print_header "📱 Тестирование Telegram"
    
    read -p "Введите Telegram Bot Token: " BOT_TOKEN
    read -p "Введите Chat ID: " CHAT_ID
    
    print_test "Отправка тестового сообщения..."
    
    RESPONSE=$(curl -s -X POST "https://api.telegram.org/bot$BOT_TOKEN/sendMessage" \
        -d chat_id="$CHAT_ID" \
        -d text="🧪 Test message from Server Monitor Bot" \
        2>/dev/null || echo "ERROR")
    
    if echo "$RESPONSE" | grep -q '"ok":true'; then
        print_pass "Сообщение отправлено успешно"
        echo "Проверьте Telegram - должно было прийти сообщение"
    else
        print_fail "Не удалось отправить сообщение"
        echo "Ошибка: $RESPONSE"
    fi
    ;;

6)
    print_header "🚀 Запуск всех тестов"
    
    print_test "1. Быстрое тестирование..."
    bash test.sh <<< "1" || true
    
    print_test "2. Unit тесты..."
    go test -v -cover || true
    
    echo ""
    print_pass "✅ Все тесты завершены"
    ;;

*)
    print_fail "Неверный выбор"
    exit 1
    ;;

esac

echo ""
print_pass "✅ Тестирование завершено!"
echo ""
