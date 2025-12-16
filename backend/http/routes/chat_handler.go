package routes

import (
	"net/http"
	"strconv"
	"summarizer/backend/http/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	svc services.ChatService
}

func NewChatHandler(ChatService services.ChatService) *ChatHandler {
	return &ChatHandler{svc: ChatService}
}

func (h *ChatHandler) Bind(r *gin.RouterGroup) {
	r.GET("/chats", h.Show)
	r.POST("/chats", h.Create)
	r.POST("/chats/:id", h.CreateMsg)
	r.DELETE("/chats/:id", h.DeleteChat)
}

func (h *ChatHandler) Show(c *gin.Context) {
	session := sessions.Default(c)
	userID, _ := session.Get("uid").(uint)

	c.IndentedJSON(http.StatusOK, h.svc.GetChats(userID))
}

func (h *ChatHandler) Create(c *gin.Context) {
	var requestBody struct {
		Prompt string `json:"prompt"`
	}
	if err := c.BindJSON(&requestBody); err != nil {
		panic(err)
	}

	session := sessions.Default(c)
	userID, _ := session.Get("uid").(uint)
	chat := h.svc.NewChat(userID, requestBody.Prompt)
	c.IndentedJSON(http.StatusOK, chat)
}

func (h *ChatHandler) CreateMsg(c *gin.Context) {
	session := sessions.Default(c)
	userID, _ := session.Get("uid").(uint)
	chatID, _ := strconv.Atoi(c.Param("id"))

	var body struct {
		Question string `json:"msg"`
	}
	if err := c.BindJSON(&body); err != nil {
		panic(err)
	}

	chat := h.svc.FindChat(userID, uint(chatID))
	entry := h.svc.Summarize(chat, body.Question, func(result string) {})
	c.JSON(http.StatusOK, entry)
}

func (h *ChatHandler) DeleteChat(c *gin.Context) {
	session := sessions.Default(c)
	userID, _ := session.Get("uid").(uint)
	chatID, _ := strconv.Atoi(c.Param("id"))

	chat := h.svc.FindChat(userID, uint(chatID))
	h.svc.DeleteChat(chat)
}
