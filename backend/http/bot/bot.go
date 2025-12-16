package bot

import (
	"context"
	"fmt"
	"strings"

	"summarizer/backend/http/config"
	"summarizer/backend/http/models"
	"summarizer/backend/http/services"

	tbot "github.com/go-telegram/bot"
	tbotmodels "github.com/go-telegram/bot/models"
)

func CreateBotHandler(authService services.AuthService, chatService services.ChatService) tbot.HandlerFunc {
	fNewText := func(ctx context.Context, b *tbot.Bot, msg *tbotmodels.Message, user *models.User) {
		b.SendMessage(ctx, &tbot.SendMessageParams{
			ChatID: msg.Chat.ID,
			Text:   "Отправьте текст для анализа",
		})
	}

	fWeb := func(ctx context.Context, b *tbot.Bot, msg *tbotmodels.Message, user *models.User) {
		b.SendMessage(ctx, &tbot.SendMessageParams{
			ChatID:    msg.Chat.ID,
			Text:      fmt.Sprintf("**Ваша ссылка для входа:**\n\n`http://localhost:8000/login?key=%s`\n\nСрок действия: 15 минут", authService.IssueLoginToken(user)),
			ParseMode: tbotmodels.ParseModeMarkdown,
		})
	}

	return func(ctx context.Context, b *tbot.Bot, update *tbotmodels.Update) {
		msg := update.Message
		if msg == nil {
			return
		}

		tid := uint(msg.From.ID)
		user := authService.GetUserByTelegramID(tid)

		if "/start" == msg.Text || user == nil {
			user := authService.CreateUser(tid, strings.TrimSpace(fmt.Sprintf("%s %s", msg.From.FirstName, msg.From.LastName)))
			fNewText(ctx, b, msg, user)
			return
		}

		if "/web" == msg.Text {
			fWeb(ctx, b, msg, user)
			return
		}

		if "/newchat" == msg.Text {
			chatService.SetUserTelegramActiveChat(user, nil)
			fNewText(ctx, b, msg, user)
			return
		}

		if user.CurrentTelegramChatID == nil {
			chat := chatService.NewChat(user.ID, msg.Text)
			chatService.SetUserTelegramActiveChat(user, chat)

			b.SetMessageReaction(ctx, &tbot.SetMessageReactionParams{
				ChatID:    msg.Chat.ID,
				MessageID: msg.ID,
				Reaction: []tbotmodels.ReactionType{{
					Type: tbotmodels.ReactionTypeTypeEmoji,
					ReactionTypeEmoji: &tbotmodels.ReactionTypeEmoji{
						Type:  "emoji",
						Emoji: "✍️",
					},
				}},
			})
			b.SendMessage(ctx, &tbot.SendMessageParams{
				ChatID: msg.Chat.ID,
				Text:   "Присылайте вопросы к тексту",
			})
		} else {
			b.SetMessageReaction(ctx, &tbot.SetMessageReactionParams{
				ChatID:    msg.Chat.ID,
				MessageID: msg.ID,
				Reaction: []tbotmodels.ReactionType{{
					Type: tbotmodels.ReactionTypeTypeEmoji,
					ReactionTypeEmoji: &tbotmodels.ReactionTypeEmoji{
						Type:  "emoji",
						Emoji: "🤔",
					},
				}},
				IsBig: func() *bool {
					b := true
					return &b
				}(),
			})
			chatService.Summarize(user.CurrentTelegramChat, msg.Text, func(answer string) {
				b.SendMessage(ctx, &tbot.SendMessageParams{
					ChatID: msg.Chat.ID,
					Text:   answer,
					ReplyParameters: &tbotmodels.ReplyParameters{
						MessageID: msg.ID,
						ChatID:    msg.Chat.ID,
					},
				})
			})
		}
	}
}

func CreateBot(cfg *config.BotConfig, handler tbot.HandlerFunc) *tbot.Bot {
	b, err := tbot.New(cfg.Token, tbot.WithDefaultHandler(handler), tbot.WithDebug())
	if err != nil {
		panic(err)
	}
	return b
}
