package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/liwei1dao/lego/sys/log"
)

func newSys(options *Options) (sys *Bot, err error) {
	sys = &Bot{options: options}
	if sys.api, err = tgbotapi.NewBotAPI(options.ApiToken); err != nil {
		log.Errorln(err)
	}
	return
}

type Bot struct {
	options *Options
	api     *tgbotapi.BotAPI
}

func (this *Bot) run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := this.api.GetUpdatesChan(u)
	if err != nil {
		log.Errorln(err)
		return
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello! I am your bot.")
				this.api.Send(msg)
			}
		}
	}
}
