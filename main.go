package main

import (
	"flag"
	"fmt"
	"image/png"
	"log"
	"os"
	"strings"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("1"),
		tgbotapi.NewKeyboardButton("2"),
		tgbotapi.NewKeyboardButton("3"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("4"),
		tgbotapi.NewKeyboardButton("5"),
		tgbotapi.NewKeyboardButton("6"),
	),
)

var yesnoKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Yes"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("No"),
	),
)

func main() {
	myToken := flag.String("token", "none", "Telegram API token")
	flag.Parse()
	if *myToken == "none" {
		fmt.Println("Add flag with Telegram Bot token")
		os.Exit(1)
	}
	bot, err := tgbotapi.NewBotAPI(*myToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "help":
				msg.Text = "TODO: mostrar help"
				bot.Send(msg)
			case "addTask":
				splitText := strings.Split(update.Message.Text, " ")
				if len(splitText) > 1 {
					task := splitText[1:]
					fmt.Printf("task: %s\n", strings.Join(task, " "))
					msg.Text = "TODO: agregar tarea a la lista: " + strings.Join(task, " ")
				}
				bot.Send(msg)
			case "getTasks":
				msg.Text = "TODO: mostrar lista de tareas"
				bot.Send(msg)
			case "status":
				msg.Text = "I'm ok."
				bot.Send(msg)
			case "open":
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
				msg.ReplyMarkup = numericKeyboard
				bot.Send(msg)
			case "close":
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				bot.Send(msg)
			case "yesOpen":
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
				msg.ReplyMarkup = yesnoKeyboard
				bot.Send(msg)
			case "getImage":
				msg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, "images/bot_image_1.png")
				msg.Caption = "Test image"
				bot.Send(msg)
			case "genQR":
				splitText := strings.Split(update.Message.Text, " ")
				text := ""
				if len(splitText) > 1 {
					text = strings.Join(splitText[1:], " ")
					fmt.Printf("text: %s\n", text)
				} else {

				}
				// Create the barcode
				qrCode, _ := qr.Encode(text, qr.M, qr.Auto)
				// Scale the barcode to 200x200 pixels
				qrCode, _ = barcode.Scale(qrCode, 200, 200)
				// create the output file
				file, _ := os.Create("images/qrcode.png")
				defer file.Close()
				// encode the barcode as png
				png.Encode(file, qrCode)

				msg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, "images/qrcode.png")
				msg.Caption = "QR image"
				bot.Send(msg)
			default:
				msg.Text = "I don't know that command"
				bot.Send(msg)
			}

		}

		/*
			splitText := strings.Split(update.Message.Text, " ")
			command := splitText[0]
			if command == "/addTask" {
				if len(splitText) > 1 {
					task := splitText[1:]
					fmt.Printf("task: %s\n", strings.Join(task, " "))
				}

			} else if command == "/getTasks" {
				txt_out := "Devolver lista de tareas"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, txt_out)
				bot.Send(msg)

			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
				msg.ReplyToMessageID = update.Message.MessageID

				bot.Send(msg)
			}
		*/

	}
}
