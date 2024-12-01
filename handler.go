package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	bot     *tgbotapi.BotAPI
	cache   *cache.Cache
	storage *Storage
}

func NewHandler(bot *tgbotapi.BotAPI, cache *cache.Cache, storage *Storage) *Handler {
	return &Handler{cache: cache, bot: bot, storage: storage}
}

func (h *Handler) Handle(u *tgbotapi.Update) {
	if u.CallbackQuery != nil {
		if u.CallbackQuery.Data == "close" {
			deleteMsg := tgbotapi.NewDeleteMessage(u.CallbackQuery.Message.Chat.ID, u.CallbackQuery.Message.MessageID)
			h.bot.Send(deleteMsg)

			return
		}
	}

	if !u.Message.IsCommand() {
		return
	}

	chatId := u.Message.Chat.ID

	var reply tgbotapi.MessageConfig

	switch u.Message.Command() {
	case "list":
		reply = h.handleList(chatId)
	case "add":
		if err := h.handleAdd(u.Message); err != nil {
			reply = tgbotapi.NewMessage(chatId, "Ошибка добавления")
		} else {
			reply = tgbotapi.NewMessage(chatId, "Успешно добавлено")
		}
	case "remove":
		reply = tgbotapi.NewMessage(chatId, "AREMOVE MESSAGE")
	default:
		reply = helpMessage(chatId)
	}

	h.bot.Send(reply)
}

func (h *Handler) handleList(chatId int64) tgbotapi.MessageConfig {
	reply := tgbotapi.NewMessage(chatId, h.list(chatId))

	rowButtons := make([]tgbotapi.InlineKeyboardButton, 0)
	closeButton := tgbotapi.NewInlineKeyboardButtonData("Close", "close")
	rowButtons = append(rowButtons, closeButton)

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup()
	inlineKeyboard.InlineKeyboard = append(inlineKeyboard.InlineKeyboard, rowButtons)
	reply.ReplyMarkup = inlineKeyboard

	return reply
}

func (h *Handler) handleAdd(message *tgbotapi.Message) error {
	sChatId := strconv.FormatInt(message.Chat.ID, 10)
	h.cache.Delete(sChatId)

	shards := strings.Split(message.Text, " ")
	game := Game{Owner: message.From.UserName, Name: strings.Join(shards[1:], " ")}
	if err := h.storage.Add(sChatId, game); err != nil {
		//send error message
		logrus.Error(err)
		return err
	}

	return nil
}

func (h *Handler) list(chatId int64) string {
	games := make([]Game, 0)

	items, ok := h.cache.Get(strconv.FormatInt(chatId, 10))
	if ok {
		games = items.([]Game)
		/*if err := json.Unmarshal([]byte(items.(string)), games); err != nil {
			logrus.Error(err)
		}*/

	} else {
		//get from file
		games = h.storage.Get(strconv.FormatInt(chatId, 10))

		h.cache.Set(strconv.FormatInt(chatId, 10), games, 5*time.Minute)
	}

	sb := strings.Builder{}
	for _, g := range games {
		sb.WriteString("#" + fmt.Sprint(g.Id) + " " + g.Name + " " + g.Owner)
	}

	return sb.String()
}

func helpMessage(chatId int64) tgbotapi.MessageConfig {
	b := strings.Builder{}
	b.WriteString("Комманды бота: \n")
	b.WriteString("/list - выводит список игр и их владельцев \n")
	b.WriteString("/add - команда для добавления игры в список \n")
	b.WriteString("/remove - команда для удаления игры по id \n")

	return tgbotapi.NewMessage(chatId, b.String())
}
