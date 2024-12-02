package main

import (
	"flag"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

const (
	ENV_DEBUG = "ENV_DEBUG" //Для локального запуска в дебаг режиме(все сообщения чата будут ссыпаться в лог)
	ENV_LOCAL = "ENV_LOCAL" //Для локального запуска (токен передается через флаг)
	ENV_POD   = "ENV_POD"   //Для запуска в контейнере (токен прокидывается через переменную окружения)
)

var (
	Token string
	Env   string
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{})

	Env = os.Getenv("ENV")
	if Env == "" {
		if err := godotenv.Load(); err != nil {
			logrus.Warning("error while reading environment %s", err.Error())
		}
	}

	Env = os.Getenv("ENV")
	if Env == "" {
		logrus.Warn("cant set environment, setting to local by default")
		Env = ENV_LOCAL
	}

	logrus.Info("ENVIRONMENT IS " + Env)

	setToken()
}

func main() {
	bot, err := tgbotapi.NewBotAPI(Token)
	if err != nil {
		logrus.Panic(err)
	}

	bot.Debug = Env == ENV_DEBUG

	cache := cache.New(5*time.Minute, 10*time.Minute)
	storage := NewStorage("db")
	handler := NewHandler(bot, cache, storage)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		handler.Handle(&update)
	}
}

func setToken() {
	if Env == ENV_LOCAL || Env == ENV_DEBUG {
		flag.StringVar(&Token, "t", "", "Bot Token")
		flag.Parse()
	}

	if Env == ENV_POD {
		Token = os.Getenv("TOKEN")
	}
}
