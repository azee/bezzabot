package main
import (
  "github.com/Syfaro/telegram-bot-api"
  "log"
	"fmt"
)

func main() {
  // подключаемся к боту с помощью токена
  bot, err := tgbotapi.NewBotAPI("243313633:AAEJ9gRzlkht8yIjRfCKItk-WFvkhztOazc")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// инициализируем канал, куда будут прилетать обновления от API
	var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	// err = bot.UpdatesChan(ucfg)
	// читаем обновления из канала
	updates, err := bot.GetUpdatesChan(ucfg)

	for {
		select {
		case update := <-updates:
			if update.Message == nil {
				continue
			}

			// Пользователь, который написал боту
			UserName := update.Message.From.UserName

			// ID чата/диалога.
			// Может быть идентификатором как чата с пользователем
			// (тогда он равен UserID) так и публичного чата/канала
			ChatID := update.Message.Chat.ID

			// Текст сообщения
			Text := update.Message.Text

			log.Printf("[%s] %d %s", UserName, ChatID, Text)

			// Ответим пользователю его же сообщением
			var reply string

			if update.Message.NewChatMember  != nil {
				// В чат вошел новый пользователь
				// Поприветствуем его
				reply = fmt.Sprintf(`Hi @%s! Bee good.`,
					update.Message.NewChatMember.UserName)
			} else {
				reply = fmt.Sprintf(`YoHoHo @%s! %s.`,
					UserName, Text)
			}


			if reply != "" {
				// Созадаем сообщение
				msg := tgbotapi.NewMessage(ChatID, reply)
				// и отправляем его
				bot.Send(msg)
			}
		}

	}
}
