# Telegram Bot - Linux Server Monitor

Телеграм-бот для мониторинга Linux серверов с 3x-ui и Xray.

## Функции
- 📊 Мониторинг статуса Xray
- 🌐 Проверка пинга (скорость интернета)
- 💻 Контроль ресурсов (CPU, RAM, диск)
- 📋 Логирование ошибок
- ⏰ Выборочный мониторинг (5, 10, 15, 30 минут)
- 4️⃣ Поддержка 4+ серверов

## Установка

### Требования
- Go 1.16+
- Telegram Bot Token
- SSH доступ к серверам

### Шаги
1. Клонируйте репозиторий
2. Скопируйте `config.example.yaml` в `config.yaml`
3. Заполните параметры серверов и Telegram token
4. Запустите: `go run main.go`

## Конфигурация

Отредактируйте `config.yaml`:

```yaml
telegram:
  token: "YOUR_BOT_TOKEN_HERE"
  chat_id: 123456789

monitor:
  check_intervals: [5, 10, 15, 30]  # минуты

servers:
  - name: "server-1"
    host: "192.168.1.1"
    port: 22
    user: "root"
    ssh_key: "/path/to/key"
    
  - name: "server-2"
    host: "192.168.1.2"
    port: 22
    user: "root"
    ssh_key: "/path/to/key"
```

## Команды бота

- `/status` - Статус всех серверов
- `/monitor 5` - Мониторить каждые 5 минут
- `/monitor 10` - Мониторить каждые 10 минут
- `/monitor 15` - Мониторить каждые 15 минут
- `/monitor 30` - Мониторить каждые 30 минут
- `/stop` - Остановить мониторинг
- `/help` - Справка

## Лицензия
MIT
