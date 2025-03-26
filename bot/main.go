package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const apiURL = "http://backend:8080/boost" // URL API сервиса накрутки

type BoostRequest struct {
	URL   string `json:"url"`
	Views int    `json:"views"`
}

func sendBoostRequest(url string, views int) (string, error) {
	data := BoostRequest{URL: url, Views: views}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Ошибка: %s", resp.Status)
	}

	return "✅ Накрутка запущена!", nil
}

func main() {
	botToken := os.Getenv("TG_BOT_TOKEN") // Получаем токен из переменной окружения
	if botToken == "" {
		log.Fatal("❌ Укажите токен бота в переменной TG_BOT_TOKEN")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = false
	log.Printf("🤖 Бот %s запущен!", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		text := update.Message.Text
		chatID := update.Message.Chat.ID

		if text == "/start" {
			msg := tgbotapi.NewMessage(chatID, "👋 Привет! Отправь ссылку и количество просмотров в формате:\n\n`https://example.com 1000`")
			msg.ParseMode = "Markdown"
			bot.Send(msg)
			continue
		}

		// Разбираем сообщение
		parts := strings.Fields(text)
		if len(parts) != 2 {
			msg := tgbotapi.NewMessage(chatID, "⚠️ Неверный формат! Отправь в формате:\n\n`https://example.com 1000`")
			msg.ParseMode = "Markdown"
			bot.Send(msg)
			continue
		}

		url := parts[0]
		views, err := strconv.Atoi(parts[1])
		if err != nil || views <= 0 {
			msg := tgbotapi.NewMessage(chatID, "⚠️ Введи корректное число просмотров!")
			bot.Send(msg)
			continue
		}

		// Отправляем запрос в API
		response, err := sendBoostRequest(url, views)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("❌ Ошибка: %s", err.Error()))
			bot.Send(msg)
			continue
		}

		msg := tgbotapi.NewMessage(chatID, response)
		bot.Send(msg)
	}
}
