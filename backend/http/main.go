package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"summarizer/backend/http/bot"
	"summarizer/backend/http/config"
	"summarizer/backend/http/middleware"
	"summarizer/backend/http/resources"
	"summarizer/backend/http/routes"
	"summarizer/backend/http/services"

	tbot "github.com/go-telegram/bot"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

func NewGinEngine(authHandler *routes.AuthHandler, chatHandler *routes.ChatHandler) *gin.Engine {
	r := gin.Default()

	store := cookie.NewStore([]byte(os.Getenv("APP_KEY")))
	r.Use(sessions.Sessions("summarizer", store))

	r.LoadHTMLGlob("templates/*")

	r.Static("/assets", "./dist/assets")
	r.GET("/", func(c *gin.Context) {
		c.File("./dist/index.html")
	})
	r.NoRoute(func(c *gin.Context) {
		c.File("./dist/index.html")
	})

	{
		publicApi := r.Group("/api")

		authHandler.Bind(publicApi)
	}
	{
		protectedApi := r.Group("/api")
		protectedApi.Use(middleware.AuthRequired())

		chatHandler.Bind(protectedApi)
	}

	return r
}

func serve(lc fx.Lifecycle, dbClient *gorm.DB, engine *gin.Engine, bot *tbot.Bot) {
	host := os.Getenv("APP_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8000"
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: engine.Handler(),
	}

	go bot.Start(context.TODO())

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			err := srv.ListenAndServe()
			return err
		},
		OnStop: func(ctx context.Context) error {
			if err := dbClient.Close(); err != nil {
				log.Printf("Failed to close DB connection: %v", err)
			}
			if err := srv.Close(); err != nil {
				log.Printf("Failed to shutdown HTTP server: %v", err)
			}
			return nil
		},
	})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	app := fx.New(
		fx.Provide(
			config.LoadDBConfig,
			resources.NewDBConnection,
			//
			config.LoadWorkerNodesConfig,
			resources.NewWorkerNodeOrchestrator,
			//
			services.NewAuthService,
			routes.NewAuthHandler,
			//
			services.NewChatService,
			routes.NewChatHandler,
			//
			config.LoadBotConfig,
			bot.CreateBotHandler,
			bot.CreateBot,
			//
			NewGinEngine,
		),
		fx.Invoke(serve),
	)

	if err := app.Start(context.Background()); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}

	<-app.Done()
}
