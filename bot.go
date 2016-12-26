package main
import (
	"fmt"
  	"github.com/Syfaro/telegram-bot-api"
	"github.com/carbocation/go-instagram/instagram"
  	"log"
)
type InstagramResponse struct {
	url string
	m   tgbotapi.Message
}
var (
	// Канал, куда будем помещать запросы к инстаграму
	instagram_req chan tgbotapi.Message
	// Канал, куда будем помещать ответы на запросы к инстаграму
	instagram_res chan InstagramResponse
)

// Читает канал с запросами к инстаграму
// Получает популярное фото и отправляет в канал ответов ссылку
func get_ig_popular() {
	// создаем и настраиваем клиента
	instagram_client := instagram.NewClient(nil)
	instagram_client.AccessToken = "983256223.99c2eff.db2234abcd7048c5b5d4960ef7e932c6"
	for {
		for msg := range instagram_req {
			// запрашиваем популярные фото
			media, _, err := instagram_client.Media.Popular()
			if err != nil {
				log.Println(err)
				continue
			}
			// возвращаем ссылку на первое фото из списка
			if len(media) > 0 {
				instagram_res <- InstagramResponse{media[0].Link, msg}
			}
		}
	}
}

func main() {
	// создаем канал для запросов к инстаграму
	instagram_req = make(chan tgbotapi.Message, 5)
	go get_ig_popular()
	instagram_res = make(chan InstagramResponse, 5)

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
			msg := tgbotapi.NewMessage(ChatID, reply)
			if update.Message.NewChatMember  != nil {
				// В чат вошел новый пользователь
				// Поприветствуем его
				reply = fmt.Sprintf(`Hi @%s! Bee good.`,
					update.Message.NewChatMember.UserName)
			}
			if Text == "/instarandom" {
				instagram_req <- *update.Message
				//reply = get_ig_popular()
				//msg.ReplyToMessageID = update.Message.MessageID
			} else {
				reply = fmt.Sprintf(`YoHoHo @%s! %s.`,
					UserName, Text)
			}


			if reply != "" {
				msg.Text = reply
				// и отправляем его
				bot.Send(msg)
			}
		case resp := <-instagram_res:
		// отправляем инстаграм-ссылку пользователю
			msg := tgbotapi.NewMessage(resp.m.Chat.ID, resp.url)
			msg.ReplyToMessageID = resp.m.MessageID
			bot.Send(msg)
		}


	}
}
