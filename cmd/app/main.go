package main

import (
	"fmt"
	"log"
	"github.com/Syfaro/telegram-bot-api"
	"github.com/opesun/goquery"
)

func main() {
	// создаем бота по токену (здесь токен заменен строкой в целях безопасности)
	bot, err := tgbotapi.NewBotAPI("5703665093:AAF_VZELQd7gZ_c_HQqcUexW8y6MC6iYyJk")
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// получаем обновления
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// GetUpdatesChan запускает и возвращает канал для получения обновлений
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	// вытаскиваем по порядку обновления
	for update := range updates {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		getCommand(bot, update)
		// // для отправки ответа на конкретное твое сообщение
		// msg.ReplyToMessageID = update.Message.MessageID
	}
}

func getCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	command := update.Message.Command()
	if command == "" {
		replyUnknowCommand(bot, update)
	} else {
		switch command {
		case "start":
			doStart(bot, update)
		case "get_weather_forecast":
			doWeatherForecast(bot, update)
		}
	}
}

func replyUnknowCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		"К сожалению, я не знаю такой команды :(\nПопробуйте выбрать другую")
	msg.ReplyToMessageID = update.Message.MessageID
	bot.Send(msg)
}

func doStart(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		"Привет, я - твой личный ассистент!\nТы можешь выбрать любую команду в меню, а я ее выполню :)")
	bot.Send(msg)
}

func doWeatherForecast(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	resp, err := goquery.ParseUrl("https://weather.rambler.ru")
	if err != nil {
		log.Panic(err)
	}

	date := resp.Find(".w4bT").Last().Text()
	summary := resp.Find(".TWnE").Last().Text()
	degree := resp.Find(".HhSR").Last().Text()
	feeling := resp.Find(".iO0y").Last().Text()
	weather := fmt.Sprintf(
		`Твой прогноз погоды готов!

		• Дата: %s
		• Погода: %s
		• Градусы: %s
		• %s`,
		date, summary, degree, feeling)
		
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, weather)
	bot.Send(msg)
}
