package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	bot     *tgbotapi.BotAPI
	cache   *cache.Cache
	storage *Storage
}

func NewHandler(bot *tgbotapi.BotAPI, cache *cache.Cache) *Handler {
	return &Handler{cache: cache, bot: bot}
}

func (h *Handler) Handle(u *tgbotapi.Update) {
	if !u.Message.IsCommand() {
		return
	}
	chatId := u.Message.Chat.ID
	var reply tgbotapi.MessageConfig

	switch u.Message.Command() {
	case "list":
		reply = tgbotapi.NewMessage(chatId, h.list(chatId))
	case "add":
		reply = tgbotapi.NewMessage(chatId, "ADD MESSAGE")
		//добавить обработчик создания новыой записи
		logrus.Info(u.Message.Text)
	case "remove":
		reply = tgbotapi.NewMessage(chatId, "AREMOVE MESSAGE")
	default:
		reply = helpMessage(chatId)
	}

	h.bot.Send(reply)
}

func (h *Handler) list(chatId int64) string {
	games := make([]Game, 0)

	items, ok := h.cache.Get(strconv.FormatInt(chatId, 10))
	if ok {
		json.Unmarshal([]byte(items.(string)), games)
	} else {
		//get from file
		var err error
		games, err = h.storage.Get(strconv.FormatInt(chatId, 10))
		if err != nil {
			logrus.Error(err)
			//дропнуть ошибку в чат?
		}
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
