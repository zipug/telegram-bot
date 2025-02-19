package application

import (
	"bot/internal/common/service/config"
	logger "bot/internal/common/service/logger/zerolog"
	"bot/internal/core/models"
	"bot/internal/core/ports"
	"bot/internal/infrastucture/repository/postgres"
	"bot/pkg/slices"
	telegramhtml "bot/pkg/telegram_html"
	"context"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/minio/minio-go/v7"
)

const (
	NoRelevantAnswer string = "no_relevant_answer"
	AskAI            string = "ask_ai"
	TmpDir           string = "tmp"
)

type Application struct {
	ctx         context.Context
	config      *config.TgBotConfig
	api         ports.ApiService
	attachments ports.AttachmentsService
	statistics  ports.StatisticsService
	tg_users    ports.TelegramUsersService
	articles    ports.ArticlesService
	bot         *tgbotapi.BotAPI
	log         *logger.Logger
	minio       ports.MinioService
	db          *postgres.PostgresRepository
	mainKB      *tgbotapi.ReplyKeyboardMarkup
	actions     map[int64]Action
}

type Action struct {
	IsQuestion   bool
	IsAIQuestion bool
	IsAnswer     bool
	Question     string
	Answers      map[string]AnswerWithId
}

type AnswerWithId struct {
	Id      int64
	Content string
}

func New(
	ctx context.Context,
	config *config.TgBotConfig,
	logger *logger.Logger,
	minio ports.MinioService,
	db *postgres.PostgresRepository,
	api ports.ApiService,
	attachments ports.AttachmentsService,
	tg_users ports.TelegramUsersService,
	articles ports.ArticlesService,
	stats ports.StatisticsService,
) *Application {
	bot, err := tgbotapi.NewBotAPI(config.Telegram.BotToken)
	if err != nil {
		panic(err)
	}
	menuBtns := [][]tgbotapi.KeyboardButton{}
	for _, btn := range config.Telegram.MainButtons {
		menuBtns = append(menuBtns, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(btn.Text)))
	}
	actions := make(map[int64]Action)
	keyboard := tgbotapi.NewReplyKeyboard(menuBtns...)
	return &Application{
		ctx:         ctx,
		config:      config,
		api:         api,
		attachments: attachments,
		statistics:  stats,
		tg_users:    tg_users,
		articles:    articles,
		bot:         bot,
		log:         logger,
		minio:       minio,
		db:          db,
		mainKB:      &keyboard,
		actions:     actions,
	}
}

func (a *Application) Run() error {
	const op = "application.Run"
	a.bot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := a.bot.GetUpdatesChan(u)

	for update := range updates {
		if data := update.CallbackData(); data != "" {
			if data == NoRelevantAnswer {
				a.onNoRelevantAnswer(
					update.CallbackQuery.From.ID,
					update.CallbackQuery.Message.Chat.ID,
					update.CallbackQuery.Message.MessageID,
				)
				continue
			}
			if data == AskAI {
				del := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
				a.sendMessage(op, del)
				a.onAskAIQuestion(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
				continue
			}
			if action, ok := a.actions[update.CallbackQuery.Message.Chat.ID]; ok {
				a.log.Log("info", "callback data", logger.WithStrAttr("data", data), logger.WithStrAttr("content", action.Question))
				a.onAnswer(
					data,
					action,
					update.CallbackQuery.From.ID,
					update.CallbackQuery.Message.Chat.ID,
					update.CallbackQuery.Message.MessageID,
				)
				continue
			}
			continue
		}

		if update.Message == nil {
			continue
		}

		cmd := update.Message.Command()
		if cmd != "" {
			switch cmd {
			case "start":
				a.messageOnStart(
					update.Message.Chat.ID,
					update.SentFrom(),
				)
				continue
			default:
				continue
			}
		}

		if action, ok := a.actions[update.Message.Chat.ID]; ok {
			a.log.Log(
				"info",
				"action",
				logger.WithStrAttr("content", action.Question),
				logger.WithBoolAttr("is_question", action.IsQuestion),
				logger.WithBoolAttr("is_ai_question", action.IsAIQuestion),
				logger.WithBoolAttr("is_answer", action.IsAnswer),
			)
			if action.IsAIQuestion {
				a.messageOnAIQuestion(
					update.Message.Chat.ID,
					update.Message.Text,
				)
				continue
			}
			if action.IsQuestion {
				a.messageOnQuestion(
					update.Message.Chat.ID,
					update.Message.Text,
				)
				continue
			}
		}

		if update.Message.Text == "✏️  Задать вопрос" {
			a.onAskQuestion(update.Message.Chat.ID, update.Message.MessageID)
			continue
		}

		if update.Message.Text == "🤖 Задать вопрос ИИ" {
			a.onAskAIQuestion(update.Message.Chat.ID, update.Message.MessageID)
			continue
		}

	}

	return nil
}

func (a *Application) sendMessage(op string, msg tgbotapi.Chattable) {
	if rmsg, err := a.bot.Send(msg); err != nil {
		if err.Error() != "json: cannot unmarshal bool into Go value of type tgbotapi.Message" {
			a.log.Log(
				"error",
				"failed to send message to user",
				logger.WithErrAttr(err),
				logger.WithStrAttr("message", rmsg.Text),
				logger.WithStrAttr("op", op),
			)
		}
	}
}

func (a *Application) messageOnStart(chat_id int64, user *tgbotapi.User) {
	const op = "application.sendMessage"
	msg := tgbotapi.NewMessage(chat_id, strings.Join(a.config.Telegram.HelloMessage, "\n"))
	msg.ReplyMarkup = a.mainKB
	msg.ParseMode = tgbotapi.ModeMarkdown
	a.bot.Send(msg)
	usr := models.Telegram{
		TelegramId: user.ID,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Username:   user.UserName,
		ChatId:     chat_id,
	}
	if err := a.tg_users.AddTelegramUser(usr); err != nil {
		a.log.Log(
			"error",
			"cound not add telegram user",
			logger.WithErrAttr(err),
		)
	}
}

func (a *Application) messageOnQuestion(chat_id int64, msg string) {
	const op = "application.messageOnQuestion"
	chat_action := tgbotapi.NewChatAction(chat_id, "typing")
	a.sendMessage(op, chat_action)
	answers, err := a.getQuestions(msg)
	if err != nil {
		a.log.Log(
			"error",
			"failed to get questions",
			logger.WithErrAttr(err),
			logger.WithStrAttr("op", op),
		)
		return
	}
	a.actions[chat_id] = Action{IsQuestion: false, IsAIQuestion: false, IsAnswer: true, Question: msg, Answers: answers}
	new_msg := tgbotapi.NewMessage(chat_id, "*По вашему вопросу найдены следующие статьи*:\nВыберите подходящий или оставьте запрос поддержке")
	new_msg.ReplyMarkup = a.createQuestinosKeyboard(answers)
	new_msg.ParseMode = tgbotapi.ModeMarkdown
	a.sendMessage(op, new_msg)
}

func (a *Application) messageOnAIQuestion(chat_id int64, msg string) {
	const op = "application.messageOnAIQuestion"
	chat_action := tgbotapi.NewChatAction(chat_id, "typing")
	a.sendMessage(op, chat_action)
	answer, err := a.api.AISearch(msg)
	if err != nil {
		a.log.Log(
			"error",
			"failed to get questions",
			logger.WithErrAttr(err),
			logger.WithStrAttr("op", op),
		)
		return
	}
	a.actions[chat_id] = Action{IsQuestion: false, IsAIQuestion: false, IsAnswer: false, Question: msg, Answers: nil}
	new_msg := tgbotapi.NewMessage(chat_id, answer.Choices[0].Message.Content)
	new_msg.ParseMode = tgbotapi.ModeMarkdown
	a.sendMessage(op, new_msg)
	record := models.Statistic{
		BotId:       a.config.Container.BotId,
		TelegramId:  chat_id,
		Question:    msg,
		ArticleName: "generated_by_ai:" + answer.Model,
		IsResolved:  true,
	}
	if err := a.statistics.AddStatisticRecord(record); err != nil {
		a.log.Log(
			"error",
			"cound not add statistic record",
			logger.WithErrAttr(err),
			logger.WithStrAttr("op", op),
		)
	}
	if err := a.articles.CreateArticle(models.Article{
		Name:        msg,
		Description: "generated_by_ai:" + answer.Model,
		Content:     answer.Choices[0].Message.Content,
		ProjectId:   a.config.Container.ProjectId,
	}); err != nil {
		a.log.Log(
			"error",
			"cound not add article",
			logger.WithErrAttr(err),
			logger.WithStrAttr("op", op),
		)
	}
}

func (a *Application) createQuestinosKeyboard(answers map[string]AnswerWithId) *tgbotapi.InlineKeyboardMarkup {
	btns := []tgbotapi.InlineKeyboardButton{}
	for key := range answers {
		btns = append(btns, tgbotapi.NewInlineKeyboardButtonData(key, key))
	}
	rows := slices.ChunkBy(btns, 2)
	rows = append(rows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("❌ Нет подходящего ответа", NoRelevantAnswer),
		tgbotapi.NewInlineKeyboardButtonData("🤖 Задать вопрос ИИ", AskAI),
	})
	res := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &res
}

func (a *Application) getQuestions(msg string) (map[string]AnswerWithId, error) {
	const op = "application.getQuestions"
	project_id := strconv.FormatInt(a.config.Container.ProjectId, 10)
	answers, err := a.api.Search(project_id, a.config.Search.SearchUrl, msg)
	if err != nil {
		a.log.Log(
			"error",
			"failed to get questions",
			logger.WithErrAttr(err),
			logger.WithStrAttr("op", op),
		)
		return nil, err
	}
	a.log.Log("info", "answers", logger.WithInt64Attr("len", int64(len(answers))))
	answerMap := make(map[string]AnswerWithId)
	for _, answer := range answers {
		answerMap[answer.Name] = AnswerWithId{Id: answer.Id, Content: answer.Content}
	}
	return answerMap, nil
}

func (a *Application) onAnswer(data string, action Action, telegram_id, chat_id int64, message_id int) error {
	const op = "application.onAnswer"
	del := tgbotapi.NewDeleteMessage(chat_id, message_id)
	a.sendMessage(op, del)
	if content, ok := action.Answers[data]; ok {
		text_content := telegramhtml.ToTelegramHTML(content.Content)
		msg := tgbotapi.NewMessage(chat_id, text_content)
		msg.ParseMode = tgbotapi.ModeHTML
		a.sendMessage(op, msg)
		attachs, err := a.attachments.GetAllAttachmentsByArticleId(content.Id)
		if err != nil {
			a.log.Log(
				"error",
				"failed to get attachments",
				logger.WithErrAttr(err),
				logger.WithStrAttr("op", op),
			)
			return err
		}
		for _, attach := range attachs {
			a.log.Log("info", "attachment", logger.WithStrAttr("name", attach.Name), logger.WithStrAttr("object_id", attach.ObjectId), logger.WithStrAttr("mimetype", attach.Mimetype))
			client := a.minio.GetClient()
			if err := client.FGetObject(a.ctx, a.config.MiniO.BucketAttachments, attach.ObjectId, TmpDir+"/"+attach.Name, minio.GetObjectOptions{}); err != nil {
				a.log.Log(
					"error",
					"failed to get object",
					logger.WithErrAttr(err),
					logger.WithStrAttr("op", op),
				)
				return err
			}
			file, err := os.Open(TmpDir + "/" + attach.Name)
			if err != nil {
				a.log.Log(
					"error",
					"failed to open file",
					logger.WithErrAttr(err),
					logger.WithStrAttr("op", op),
				)
				return err
			}
			reader := tgbotapi.FileReader{
				Name:   attach.Name,
				Reader: file,
			}
			if attach.Mimetype != "image/png" &&
				attach.Mimetype != "image/jpg" &&
				attach.Mimetype != "image/jpeg" &&
				attach.Mimetype != "image/gif" {
				file := tgbotapi.NewDocument(chat_id, reader)
				a.sendMessage(op, file)
			} else {
				image := tgbotapi.NewPhoto(chat_id, reader)
				a.sendMessage(op, image)
			}
		}
		record := models.Statistic{
			BotId:       a.config.Container.BotId,
			TelegramId:  telegram_id,
			ArticleId:   content.Id,
			Question:    action.Question,
			ArticleName: data,
			IsResolved:  true,
		}
		if err := a.statistics.AddStatisticRecord(record); err != nil {
			a.log.Log(
				"error",
				"cound not add statistic record",
				logger.WithErrAttr(err),
				logger.WithStrAttr("op", op),
			)
		}
	}
	return nil
}

func (a *Application) onNoRelevantAnswer(telegram_id, chat_id int64, message_id int) {
	const op = "application.onNoRelevantAnswer"
	var question string
	if action, ok := a.actions[chat_id]; ok {
		question = action.Question
	}
	a.log.Log("info", "no relevant answer", logger.WithStrAttr("question", question))
	msg := tgbotapi.NewMessage(chat_id, "Ваш запрос отправлен в службу поддержки")
	msg.ParseMode = tgbotapi.ModeHTML
	a.sendMessage(op, msg)
	del := tgbotapi.NewDeleteMessage(chat_id, message_id)
	a.sendMessage(op, del)
	record := models.Statistic{
		BotId:       a.config.Container.BotId,
		TelegramId:  telegram_id,
		Question:    question,
		ArticleName: NoRelevantAnswer,
		IsResolved:  false,
	}
	if err := a.statistics.AddStatisticRecord(record); err != nil {
		a.log.Log(
			"error",
			"cound not add statistic record",
			logger.WithErrAttr(err),
			logger.WithStrAttr("op", op),
		)
	}
}

func (a *Application) onAskQuestion(chat_id int64, message_id int) {
	const op = "application.onAskQuestion"
	/* del := tgbotapi.NewDeleteMessage(chat_id, message_id)
	a.sendMessage(op, del) */
	a.actions[chat_id] = Action{IsQuestion: true}
	msg := tgbotapi.NewMessage(chat_id, "Введите ваш вопрос")
	msg.ParseMode = tgbotapi.ModeMarkdown
	a.sendMessage(op, msg)
}

func (a *Application) onAskAIQuestion(chat_id int64, message_id int) {
	const op = "application.onAskAIQuestion"
	/* del := tgbotapi.NewDeleteMessage(chat_id, message_id)
	a.sendMessage(op, del) */
	a.actions[chat_id] = Action{IsAIQuestion: true}
	msg := tgbotapi.NewMessage(chat_id, "Введите ваш вопрос")
	msg.ParseMode = tgbotapi.ModeMarkdown
	a.sendMessage(op, msg)
}
