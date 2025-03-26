### Почему такой README
- **Описание**: Кратко объясняет, что делает проект.
- **Инструкции**: Подробно описывает установку и запуск (с Docker и без).
- **Использование**: Показывает, как взаимодействовать с ботом и проверять просмотры.
- **Диагностика**: Помогает решить типичные проблемы, с которыми вы столкнулись (например, `connection refused` или `no such file`).
# View Booster


View Booster — это сервис для накрутки просмотров веб-страниц с использованием Telegram-бота для управления и трекера для подсчёта реальных просмотров.

## Структура проекта
Проект состоит из трёх микросервисов:
1. **Backend** — API для запуска накрутки просмотров с использованием headless Chromium.
2. **Bot** — Telegram-бот для взаимодействия с пользователем и отправки запросов на backend.
3. **Counter** — Сервис для подсчёта реальных просмотров страниц.

## Требования
- **Docker** и **Docker Compose** для запуска сервисов.
- Токен Telegram-бота (получите у [@BotFather](https://t.me/BotFather)).
- Go 1.23+ (если вы хотите собирать без Docker).

## Установка

1. **Клонируйте репозиторий**:
   ```bash
   git clone https://github.com/yourusername/view-booster.git
   cd view-booster
   ```

2. **Создайте файл `.env`**:
   В корневой директории создайте файл `.env` и добавьте токен Telegram-бота:
   ```bash
   TG_BOT_TOKEN=your_telegram_bot_token_here
   ```

3. **Проверьте структуру**:
   Убедитесь, что директории содержат нужные файлы:
   ```
   view-booster/
   ├── backend/
   │   ├── Dockerfile
   │   └── main.go
   ├── bot/
   │   ├── Dockerfile
   │   └── main.go
   ├── counter/
   │   ├── Dockerfile
   │   └── main.go
   ├── docker-compose.yml
   └── .env
   ```

## Запуск

### Через Docker Compose
1. **Соберите и запустите сервисы**:
   ```bash
   docker-compose up -d --build
   ```
   - `-d`: Запускает в фоновом режиме.
   - `--build`: Пересобирает образы.

2. **Проверьте статус**:
   ```bash
   docker ps
   ```
   Вы должны увидеть три контейнера: `view-booster_backend`, `view-booster_bot`, `view-booster_counter`.

3. **Просмотрите логи**:
   ```bash
   docker-compose logs backend
   docker-compose logs bot
   docker-compose logs counter
   ```

### Отдельно (например, counter)
1. Перейдите в директорию сервиса:
   ```bash
   cd counter
   ```
2. Соберите образ:
   ```bash
   docker build -t counter-app .
   ```
3. Запустите контейнер:
   ```bash
   docker run -d -p 5001:5001 --name counter-container counter-app
   ```

## Использование
1. **Найдите бота в Telegram**:
   - Используйте имя бота, указанное при создании токена (например, `@YourViewBoosterBot`).

2. **Отправьте команду**:
   - Начните с `/start` для получения инструкций.
   - Отправьте URL и количество просмотров в формате:
     ```
     https://example.com 1000
     ```
   - Бот ответит: `✅ Накрутка запущена!`.

3. **Проверка реальных просмотров**:
   - Отправьте GET-запрос на `counter`:
     ```bash
     curl "http://localhost:5001/real-views?url=https://example.com"
     ```
     Ответ:
     ```json
     {"url":"https://example.com","views":0}
     ```

## Остановка
- Остановить все сервисы:
  ```bash
  docker-compose down
  ```
- Удалить образы (опционально):
  ```bash
  docker rmi view-booster_backend view-booster_bot view-booster_counter
  ```

## Устранение неполадок
- **Ошибка "connection refused"**:
  - Убедитесь, что `backend` запущен (`docker ps`).
  - Проверьте, что `bot` использует `http://backend:5000/boost` вместо `localhost`.

- **Ошибка "no such file or directory"**:
  - Проверьте наличие `main.go` в каждой директории:
    ```bash
    ls -la ./backend ./bot ./counter
    ```
  - Пересоберите без кэша:
    ```bash
    docker-compose build --no-cache
    ```

- Логи для диагностики:
  ```bash
  docker-compose logs <service_name>
  ```

## Зависимости
- **Backend**: `github.com/chromedp/chromedp`, `github.com/gofiber/fiber/v2`.
- **Bot**: `github.com/go-telegram-bot-api/telegram-bot-api/v5`.
- **Counter**: Только стандартная библиотека Go.

## Примечания
- Сервисы используют уникальные `WORKDIR` (`/backend`, `/bot`, `/counter`) для изоляции.
- Порты по умолчанию: `5000` (backend), `5001` (counter).

## Лицензия
MIT License (или укажите свою).

---
