# 🧪 Server Monitor Bot - Руководство тестирования

Полное руководство по тестированию всех функций бота.

## 📋 Содержание

1. [Быстрое тестирование](#быстрое-тестирование)
2. [Unit тесты](#unit-тесты)
3. [Тестирование SSH](#тестирование-ssh)
4. [Тестирование Telegram](#тестирование-telegram)
5. [Тестирование мониторинга](#тестирование-мониторинга)
6. [Чек-лист перед продакшеном](#чек-лист-перед-продакшеном)

## ⚡ Быстрое тестирование

### Способ 1: Используя скрипт (рекомендуется)

```bash
# Быстрый тест (конфиг, синтаксис, зависимости)
sudo bash test.sh
# Выберите пункт 1
```

### Способ 2: Вручную

```bash
# 1. Проверка файлов
ls -la main.go config.yaml config.example.yaml

# 2. Проверка Go синтаксиса
go fmt ./...
go vet ./...

# 3. Проверка конфига
grep -n "telegram:" config.yaml
grep -n "servers:" config.yaml

# 4. Проверка зависимостей
go mod verify
go mod tidy
```

## 🧬 Unit тесты

### Запуск всех тестов

```bash
# Стандартный запуск
go test -v

# С покрытием кода
go test -v -cover

# С подробным выводом
go test -v -cover -run=Test
```

### Запуск конкретного теста

```bash
# Только SSH тесты
go test -v -run TestSSH

# Только мониторинг
go test -v -run TestMonitoring

# Только Xray
go test -v -run TestXray

# Только ошибки
go test -v -run TestError
```

### Что тестируется

| Функция | Статус | Описание |
|---------|--------|---------|
| SSH подключение с ключом | ✓ | Базовая проверка структуры |
| SSH подключение с паролем | ✓ | Базовая проверка структуры |
| Обнаружение Xray | ✓ | Требует реального сервера |
| Получение CPU | ✓ | Требует реального сервера |
| Получение RAM | ✓ | Требует реального сервера |
| Получение диска | ✓ | Требует реального сервера |
| Конфигурация | ✓ | Проверка загрузки |
| Интервалы мониторинга | ✓ | Валидация |
| Обработка ошибок | ✓ | Базовая проверка |
| Алерты | ✓ | Проверка порогов |

## 🔐 Тестирование SSH

### Интерактивный тест

```bash
sudo bash test.sh
# Выберите пункт 4 (SSH)
```

### Ручное тестирование SSH ключа

```bash
# 1. Генерация ключа (если нет)
ssh-keygen -t rsa -N "" -f ~/.ssh/id_rsa

# 2. Копирование на сервер
ssh-copy-id -i ~/.ssh/id_rsa.pub root@192.168.1.100

# 3. Тест подключения
ssh -i ~/.ssh/id_rsa root@192.168.1.100 "echo SSH Key Works"

# 4. Проверка команд, которые нужны боту
ssh -i ~/.ssh/id_rsa root@192.168.1.100 "ps aux | grep xray"
ssh -i ~/.ssh/id_rsa root@192.168.1.100 "top -bn1 | head -5"
ssh -i ~/.ssh/id_rsa root@192.168.1.100 "free | grep Mem"
ssh -i ~/.ssh/id_rsa root@192.168.1.100 "df /"
```

### Ручное тестирование SSH пароля

```bash
# Устанавливаем sshpass (если нужно)
sudo apt-get install -y sshpass

# 1. Тест подключения с паролем
sshpass -p 'your_password' ssh -o StrictHostKeyChecking=no root@192.168.1.100 "echo SSH Password Works"

# 2. Проверка команд
sshpass -p 'your_password' ssh root@192.168.1.100 "ps aux | grep xray"
```

### Диагностика SSH проблем

```bash
# Проверка SSH конфига
ssh -v root@192.168.1.100

# Проверка прав на ключ
ls -la ~/.ssh/id_rsa  # должно быть 600

# Проверка на сервере
ssh root@192.168.1.100 "cat ~/.ssh/authorized_keys"

# Проверка SSH сервиса на сервере
ssh root@192.168.1.100 "sudo systemctl status ssh"
```

## 📱 Тестирование Telegram

### Получение Bot Token

```bash
# 1. Открыть Telegram
# 2. Найти @BotFather
# 3. Команда /newbot
# 4. Следовать инструкциям
# 5. Скопировать token вида: 123456:ABC-DEF...
```

### Получение Chat ID

```bash
# Способ 1: Через @userinfobot
# 1. Найти @userinfobot в Telegram
# 2. Команда /start
# 3. Скопировать User ID

# Способ 2: Через API
# 1. Написать боту сообщение
# 2. Открыть URL (замените TOKEN):
curl "https://api.telegram.org/botTOKEN/getUpdates"
# 3. Найти "chat":{"id":ваш_ID}
```

### Интерактивный тест

```bash
sudo bash test.sh
# Выберите пункт 5 (Telegram)
# Введите token и chat_id
```

### Ручное тестирование

```bash
# Замените TOKEN и CHAT_ID на ваши
export BOT_TOKEN="123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
export CHAT_ID="987654321"

# 1. Отправить простое сообщение
curl -X POST "https://api.telegram.org/bot$BOT_TOKEN/sendMessage" \
  -d chat_id="$CHAT_ID" \
  -d text="Test message"

# 2. Отправить сообщение с Markdown
curl -X POST "https://api.telegram.org/bot$BOT_TOKEN/sendMessage" \
  -d chat_id="$CHAT_ID" \
  -d parse_mode="Markdown" \
  -d text="*Bold* _Italic_ \`Code\`"

# 3. Проверить, получено ли сообщение
curl "https://api.telegram.org/bot$BOT_TOKEN/getUpdates"
```

## 🔍 Тестирование мониторинга

### 1. Проверка интервалов

```bash
# Проверить, что интервалы работают правильно
go test -v -run TestMonitoringInterval
```

### 2. Проверка получения метрик

Вручную на сервере:

```bash
# CPU
top -bn1 | grep "Cpu(s)" | awk '{print int($2)}'

# RAM
free | grep Mem | awk '{print int($3/$2 * 100)}'

# Диск
df / | tail -1 | awk '{print int($5)}'
```

### 3. Проверка Xray

```bash
# Проверить, запущен ли Xray
ps aux | grep -i xray | grep -v grep

# Если не запущен:
ps aux | grep -i xray | grep -v grep | wc -l  # должно быть 0
```

### 4. Проверка ping

```bash
# Локально
nc -zv 8.8.8.8 53

# Через SSH
ssh root@server "nc -zv 8.8.8.8 53"
```

## 🔧 Полный тест перед продакшеном

### Чек-лист

- [ ] **Конфигурация**
  - [ ] Заполнен telegram.token
  - [ ] Заполнен telegram.chat_id
  - [ ] Добавлены все серверы (макс 15)
  - [ ] Проверены SSH данные каждого сервера
  - [ ] Проверены пороги алертов

- [ ] **SSH подключение**
  - [ ] Тест подключения к каждому серверу
  - [ ] Проверка SSH ключей или паролей
  - [ ] Проверка прав доступа на ключи (600)
  - [ ] Проверка доступных команд (ps, top, free, df)

- [ ] **Xray**
  - [ ] Xray запущен на каждом сервере
  - [ ] Команда `ps aux | grep xray` работает
  - [ ] Порт 10085 открыт (если требуется)

- [ ] **Telegram**
  - [ ] Bot token валиден
  - [ ] Chat ID правильный
  - [ ] Бот может отправлять сообщения
  - [ ] Сообщения форматируются правильно

- [ ] **Мониторинг**
  - [ ] Интервалы работают (5, 10, 15, 30 минут)
  - [ ] Метрики обновляются корректно
  - [ ] Параллельные подключения работают
  - [ ] Алерты срабатывают при превышении порогов

- [ ] **Системные сервисы**
  - [ ] Systemd сервис работает
  - [ ] Автозапуск включен
  - [ ] Логи сохраняются
  - [ ] Бот восстанавливается после перезагрузки

### Автоматический чек

```bash
# Запустить все тесты
sudo bash test.sh
# Выберите пункт 6
```

## 📊 Результаты тестирования

### Ожидаемый результат

```
✅ Быстрое тестирование завершено
✅ Unit тесты пройдены
✅ SSH подключение работает
✅ Telegram интеграция работает
✅ Мониторинг функционирует
✅ Все тесты завершены!
```

### Интерпретация результатов

| Результат | Значение |
|-----------|----------|
| ✅ (зелено) | Тест пройден успешно |
| ❌ (красно) | Тест не пройден, требуется исправление |
| ⊘ (желтый) | Тест пропущен (требуется реальный сервер) |
| ⓘ (синий) | Информационное сообщение |

## 🐛 Решение проблем при тестировании

### SSH подключение не удалось

```bash
# 1. Проверить IP адрес
ping 192.168.1.100

# 2. Проверить SSH доступ
ssh -vvv -i ~/.ssh/id_rsa root@192.168.1.100

# 3. Проверить права на ключ
chmod 600 ~/.ssh/id_rsa

# 4. Проверить на сервере
ssh root@server "sudo systemctl status ssh"
```

### Telegram сообщения не отправляются

```bash
# 1. Проверить token
curl "https://api.telegram.org/botTOKEN/getMe"

# 2. Проверить chat_id
curl "https://api.telegram.org/botTOKEN/getUpdates"

# 3. Проверить интернет соединение
curl -I https://api.telegram.org
```

### Метрики не получаются

```bash
# 1. Проверить доступность команд
ssh root@server "which top free df ps"

# 2. Запустить команды вручную
ssh root@server "top -bn1 | head -5"
ssh root@server "free | grep Mem"
ssh root@server "df /"
```

## 📝 Запись результатов

### Сохранить результаты тестов

```bash
# Запустить тесты и сохранить результаты
go test -v -cover 2>&1 | tee test_results.txt

# Просмотреть результаты
cat test_results.txt
```

## 🚀 Продакшн

После успешного прохождения всех тестов:

```bash
# 1. Скопировать конфиг на сервер
sudo cp config.yaml /opt/server-monitor/

# 2. Запустить бота
sudo systemctl start server-monitor

# 3. Проверить статус
sudo systemctl status server-monitor

# 4. Просмотреть логи
sudo journalctl -u server-monitor -f
```

---

**Версия:** 1.0.0  
**Последнее обновление:** 2026-06-01
