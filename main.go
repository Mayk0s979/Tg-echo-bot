package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("не удалось загрузить .env")
	}

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("BOT_TOKEN pust")
	}

	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println("Ошибка создания бота:", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go sendWelcomeMessage(bot, ctx)

	updates, err := bot.UpdatesViaLongPolling(ctx, &telego.GetUpdatesParams{})
	if err != nil {
		fmt.Println("Ошибка при запуске long polling:", err)
		os.Exit(1)
	}

	for update := range updates {
		if update.Message != nil {
			chatId := tu.ID(update.Message.Chat.ID)

			_, err := bot.CopyMessage(ctx, tu.CopyMessage(
				chatId,
				chatId,
				update.Message.MessageID,
			))
			if err != nil {
				fmt.Printf("Ошибка копирования сообщения: %v\n", err)
			}
		}
	}
}

func sendWelcomeMessage(bot *telego.Bot, ctx context.Context) {
	adminChatIDstr := os.Getenv("ADMIN_CHAT_ID")
	if adminChatIDstr == "" {
		fmt.Println("ADMIN_CHAT_ID pust")
	}

	adminChatID, err := strconv.ParseInt(adminChatIDstr, 10, 64)
	if err != nil {
		fmt.Printf("Ошибка преобразования DMIN_CHAT_ID: %v\n", err)
	}

	time.Sleep(1 * time.Second)

	text := "🤖 Бот успешно запущен!"

	_, err = bot.SendMessage(ctx, tu.Message(
		tu.ID(adminChatID),
		text,
	))

	if err != nil {
		fmt.Printf("Ошибка отправки приветствия: %v\n", err)
	} else {
		fmt.Println("Приветствие отправлено.")
	}

}
