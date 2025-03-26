package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/gofiber/fiber/v2"
)

// Список User-Agent'ов для рандомизации
var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 15_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.2 Mobile/15E148 Safari/604.1",
}

// Прокси (если есть)
var proxies = []string{
	"socks5://127.0.0.1:9050",                         // Пример: локальный Tor прокси
	"http://username:password@proxy.example.com:8080", // Пример с авторизацией
}

// Максимальное количество потоков
const maxConcurrent = 5

// Функция для накрутки просмотров
func increaseViews(url string, views int) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrent) // Ограничение потоков

	for i := 0; i < views; i++ {
		wg.Add(1)
		sem <- struct{}{} // Лимитируем потоки

		go func(i int) {
			defer wg.Done()
			defer func() { <-sem }() // Освобождаем слот

			// Рандомный User-Agent
			userAgent := userAgents[rand.Intn(len(userAgents))]

			// Рандомный прокси (если есть)
			proxy := ""
			if len(proxies) > 0 {
				proxy = proxies[rand.Intn(len(proxies))]
			}

			// Запускаем Chrome с настройками
			ctx, cancel := chromedp.NewExecAllocator(context.Background(),
				chromedp.NoFirstRun,
				chromedp.NoDefaultBrowserCheck,
				chromedp.Flag("headless", true),
				chromedp.Flag("disable-gpu", true),
				chromedp.Flag("proxy-server", proxy), // Прокси
			)
			defer cancel()

			browserCtx, cancel := chromedp.NewContext(ctx)
			defer cancel()

			// Задачи для эмуляции пользователя
			tasks := chromedp.Tasks{
				chromedp.Navigate(url), // Переход на сайт
				chromedp.ActionFunc(func(ctx context.Context) error { // Установка User-Agent
					return chromedp.Evaluate(fmt.Sprintf(`navigator.__defineGetter__('userAgent', function(){ return '%s'; });`, userAgent), nil).Do(ctx)
				}),
				chromedp.Sleep(time.Duration(rand.Intn(5)+2) * time.Second),  // Задержка
				chromedp.ScrollIntoView(`body`),                              // Прокрутка страницы вниз
				chromedp.Sleep(time.Duration(rand.Intn(10)+5) * time.Second), // Ещё задержка
			}

			// Запуск задач
			err := chromedp.Run(browserCtx, tasks)
			if err != nil {
				log.Printf("Ошибка при загрузке страницы: %v", err)
				return
			}

			fmt.Printf("[%d/%d] Просмотр добавлен! (User-Agent: %s, Proxy: %s)\n", i+1, views, userAgent, proxy)
		}(i)
	}

	wg.Wait() // Ждём завершения всех горутин
}
func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Привет! Перейдите на /boost для накрутки просмотров.")
	})

	app.Post("/boost", func(c *fiber.Ctx) error {
		type Request struct {
			URL   string `json:"url"`
			Views int    `json:"views"`
		}

		var req Request
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).SendString("Ошибка парсинга JSON")
		}

		go increaseViews(req.URL, req.Views)
		return c.JSON(fiber.Map{"message": "Накрутка запущена!"})
	})

	log.Println("Сервер запущен на :8080")
	log.Fatal(app.Listen(":8080"))
}
