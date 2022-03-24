package telegabot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"gobot/config"
	"strconv"
)

type TelegaApi struct {
	Bot   *tgbotapi.BotAPI
	Log   *logrus.Logger
	Chats []int64
}

func NewTelegaBot(log *logrus.Logger, config *config.Config) (*TelegaApi, error) {
	bot, err := tgbotapi.NewBotAPI(config.Telega.ApiToken)
	if err != nil {
		return nil, err
	}

	bot.Debug = config.Telega.Debug
	// tgbotapi.SetMyCommandsConfig{}
	t := &TelegaApi{
		Bot: bot,
		Log: log,
	}

	t.Chats = append(t.Chats, -1001771210076)

	go t.processMessages()
	t.Log.Infof("Telegram bot %s connected!", bot.Self.UserName)

	return t, err
}

func (t *TelegaApi) processMessages() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := t.Bot.GetUpdatesChan(u)
	for update := range updates {
		t.Log.Info("update")
		if update.Message != nil {
			t.Log.Infof("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			switch update.Message.Command() {
			case "start":
				t.Chats = append(t.Chats, update.Message.Chat.ID)
				msg.Text = "🤖 " + strconv.FormatInt(update.Message.Chat.ID, 10)
			default:
				msg.Text = "I don't know that command"
			}

			_, _ = t.Bot.Send(msg)
		}
	}
}

func (t *TelegaApi) SendMessages(text string) {
	for _, id := range t.Chats {
		msg := tgbotapi.NewMessage(id, text)
		_, _ = t.Bot.Send(msg)
	}
}
