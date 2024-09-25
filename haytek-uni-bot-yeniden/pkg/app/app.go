package app

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/haytek-uni-bot-yeniden/app/service"
	"github.com/haytek-uni-bot-yeniden/pkg/bot"
	"github.com/haytek-uni-bot-yeniden/pkg/config"
	"github.com/haytek-uni-bot-yeniden/pkg/database"
	"github.com/haytek-uni-bot-yeniden/pkg/utils"
	"github.com/uptrace/bun"
	"google.golang.org/api/sheets/v4"
	"log"
	"time"
)

type Ctx struct {
	tgbotapi.Update
	//Logger
}

type CommandHandler func(ctx *Ctx, params []string) (string, error)

type CommandsMap map[string]CommandHandler

type ScheduledHandler func() []ScheduledResponse //aynı vakitte farklı mesajlar farklı kişilere gönderilebilir. ör: saat 11'de herkese özel kontrol mesajı gitmesi gibi. ya da hatırlatma mesajı falan.

type ScheduledResponse struct {
	UserID int64
	Result string
	Error  error
}

type ScheduledMap map[time.Time]ScheduledHandler

type App struct {
	Bot       *tgbotapi.BotAPI
	DB        *bun.DB
	Sheets    *sheets.Service
	Commands  CommandsMap
	Scheduled ScheduledMap
}

func New() *App {
	config.Setup()
	cfg := config.Get()
	database.New(cfg.Database)
	service.InitSheetsService(cfg.Sheets)
	bot.New(cfg.Bot)
	app := App{
		Bot:       bot.Get(),
		DB:        database.Get(),
		Sheets:    service.Get(),
		Commands:  make(CommandsMap),
		Scheduled: make(ScheduledMap),
	}
	return &app
}

func (app *App) Start() {
	//sessionService service.ISessionService

	app.StartScheduledJobs()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	commands := app.Commands
	//var m model.Session
	var err error
	var str string
	log.Println("listening for updates")
	updates := app.Bot.GetUpdatesChan(u)
	for update := range updates {
		//---------------------BURASI AZCIK PROTOTİP ---------
		//message session al
		//userid := update.Message.From.ID
		//err = app.DB.NewSelect().Model(&m).Where("user_id = ?", userid).Scan(context.Background())
		//m, err = sessionService.GetByUserID(context.Background(), userid)
		//if err != nil {
		//	yok ise oluştur.
		//m = model.Session{
		//	UserID:      userid,
		//	LastCommand: update.Message.Command(),
		//	State:       model.StateNone,
		//}
		//err = sessionService.Create(context.Background(), &m)
		//if err != nil {
		//	log.Println(err)
		//	continue
		//}
		//}
		//--------------------------------
		if update.Message == nil {
			continue
		}
		if !update.Message.IsCommand() {
			continue
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		if !commands.isRegistered(update.Message.Command()) {
			msg.Text = "command not found"
		} else {
			str, err = commands.ExecuteCommand(update.Message.Command(), &Ctx{Update: update}, update.Message.CommandArguments())
			msg.Text = str
			if err != nil {
				msg.Text = err.Error()
			}
		}

		if _, err = app.Bot.Send(msg); err != nil {
			log.Panic(err)
		}

		//SESSION UPDATE
		//m.LastCommand = update.Message.Command()
		//err = sessionService.Update(context.Background(), m)
		//if err != nil {
		//	log.Panic(err)
		//}
	}
}

func (app *App) StartScheduledJobs() {
	for t, handler := range app.Scheduled {
		go startTask(t, handler, app.Bot)
	}
	log.Println("initial")
}

func (commands CommandsMap) RegisterCommand(command string, handler CommandHandler) {
	if _, ok := commands[command]; ok {
		log.Println("command already registered")
		return
	}
	commands[command] = handler
}

func (commands CommandsMap) ExecuteCommand(command string, ctx *Ctx, params string) (string, error) {
	return commands[command](ctx, utils.ParseCommandArguments(params))
}
func (commands CommandsMap) resolve(update *tgbotapi.Update, command string) (string, error) {
	str, err := commands.ExecuteCommand(command, &Ctx{Update: *update}, update.Message.CommandArguments())
	return str, err
}

func (commands CommandsMap) isRegistered(command string) bool {
	_, ok := commands[command]
	return ok
}

func (s ScheduledMap) RegisterScheduled(timestr string, handler ScheduledHandler) {
	t := utils.ParseStrTime(timestr)
	if _, ok := s[t.Round(time.Minute)]; ok {
		log.Println("scheduled task already registered for this time\tmake sure you dont register a handler for same hh:mm")
		return
	}
	s[t] = handler
}

// her bir task için &bot göndermek mantıklı olmayabilir ama bakalım.
//bide burası yine aduket oldu ya

func startTask(t time.Time, handler ScheduledHandler, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(0, "")
	timer := newTimer(t)
	var startTime time.Time
	for {
		select {
		case _ = <-timer.C:
			startTime = time.Now()
			res := handler()
			for _, result := range res {
				msg.ChatID = result.UserID
				if result.Error != nil {
					msg.Text = result.Error.Error()
				} else {
					msg.Text = result.Result
				}
				if _, err := bot.Send(msg); err != nil {
					log.Println(err)
					continue
				}
			}
			timer.Reset(24*time.Hour - time.Now().Sub(startTime))
		}
	}
}

func newTimer(t time.Time) *time.Timer {
	sub := t.Sub(time.Now())
	if sub < 0 {
		sub += 24 * time.Hour
	}
	return time.NewTimer(sub)
}
