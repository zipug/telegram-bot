package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type MiniO struct {
	Host              string        `mapstructure:"host"`
	Port              int           `mapstructure:"port"`
	User              string        `mapstructure:"user"`
	Password          string        `mapstructure:"password"`
	BucketArticles    string        `mapstructure:"articles_bucket"`
	BucketAttachments string        `mapstructure:"attachments_bucket"`
	BucketAvatars     string        `mapstructure:"avatars_bucket"`
	UrlLifetime       time.Duration `mapstructure:"url_lifetime"`
	UseSsl            bool          `mapstructure:"use_ssl"`
}

type Search struct {
	SearchUrl string `mapstructure:"search_url"`
}

type Postgres struct {
	Host           string `mapstructure:"host"`
	Port           int    `mapstructure:"port"`
	User           string `mapstructure:"user"`
	Password       string `mapstructure:"password"`
	DBName         string `mapstructure:"db_name"`
	SSLMode        string `mapstructure:"ssl_mode"`
	MigrationsPath string `mapstructure:"migrations_path"`
}

type Telegram struct {
	BotToken       string   `mapstructure:"bot_token"`
	InformationUrl string   `mapstructure:"information_url"`
	HelloMessage   []string `mapstructure:"hello_message"`
	ErrorMessage   []string `mapstructure:"error_message"`
	MainImage      string   `mapstructure:"main_image"`
	MainButtons    []Button `mapstructure:"main_buttons"`
}

type Container struct {
	BotId       int64  `mapstructure:"bot_id"`
	ProjectId   int64  `mapstructure:"project_id"`
	UserId      int64  `mapstructure:"user_id"`
	Name        string `mapstructure:"name"`
	Description string `mapstructure:"description"`
	Icon        string `mapstructure:"icon"`
}

type Button struct {
	Text string `mapstructure:"text"`
	Data string `mapstructure:"data"`
}

type TgBotConfig struct {
	Env        ENV       `mapstructure:"env"`
	Postgres   Postgres  `mapstructure:"postgres"`
	MiniO      MiniO     `mapstructure:"minio"`
	Search     Search    `mapstructure:"search"`
	Container  Container `mapstructure:"container"`
	Telegram   Telegram  `mapstructure:"telegram"`
	configPath string
}

type ENV string

const (
	ENV_DEVELOPMENT ENV = "development"
	ENV_PRODUCTION  ENV = "production"
)

func NewConfigService() *TgBotConfig {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/app/configs/")
	viper.AddConfigPath("configs/")
	viper.BindEnv("env", "ENV")
	viper.BindEnv("postgres.host", "POSTGRES_HOST")
	viper.BindEnv("postgres.port", "POSTGRES_PORT")
	viper.BindEnv("postgres.user", "POSTGRES_USER")
	viper.BindEnv("postgres.password", "POSTGRES_PASSWORD")
	viper.BindEnv("postgres.db_name", "POSTGRES_DB_NAME")
	viper.BindEnv("postgres.ssl_mode", "POSTGRES_SSL_MODE")
	viper.BindEnv("postgres.migrations_path", "POSTGRES_MIGRATIONS_PATH")
	viper.BindEnv("minio.host", "MINIO_HOST")
	viper.BindEnv("minio.port", "MINIO_PORT")
	viper.BindEnv("minio.user", "MINIO_ROOT_USER")
	viper.BindEnv("minio.password", "MINIO_ROOT_PASSWORD")
	viper.BindEnv("minio.articles_bucket", "MINIO_ARTICLES_BUCKET")
	viper.BindEnv("minio.attachments_bucket", "MINIO_ATTACHMENTS_BUCKET")
	viper.BindEnv("minio.avatars_bucket", "MINIO_AVATARS_BUCKET")
	viper.BindEnv("minio.use_ssl", "MINIO_USE_SSL")
	viper.BindEnv("minio.url_lifetime", "MINIO_URL_LIFETIME")
	viper.BindEnv("search.search_url", "SEARCH_URL")
	viper.BindEnv("container.bot_id", "CONTAINER_BOT_ID")
	viper.BindEnv("container.project_id", "CONTAINER_PROJECT_ID")
	viper.BindEnv("container.user_id", "CONTAINER_USER_ID")
	viper.BindEnv("container.name", "CONTAINER_NAME")
	viper.BindEnv("container.description", "CONTAINER_DESCRIPTION")
	viper.BindEnv("container.icon", "CONTAINER_ICON")
	viper.BindEnv("telegram.bot_token", "TELEGRAM_BOT_TOKEN")
	viper.BindEnv("telegram.information_url", "TELEGRAM_INFORMATION_URL")
	viper.BindEnv("telegram.hello_message", "TELEGRAM_HELLO_MESSAGE")
	viper.BindEnv("telegram.error_message", "TELEGRAM_ERROR_MESSAGE")
	viper.BindEnv("telegram.main_image", "TELEGRAM_MAIN_IMAGE")
	viper.BindEnv("telegram.main_buttons", "TELEGRAM_MAIN_BUTTONS")
	viper.AutomaticEnv()
	viper.SetDefault("telegram.information_url", "")
	viper.SetDefault("telegram.hello_message", []string{
		"*👨‍💻 Привет! Я бот поддержки*",
		"Я могу ответить на все ваши вопросы!",
		"*Напишите мне,* и я постараюсь помочь вам!",
	})
	viper.SetDefault("telegram.error_message", []string{"Если Вы не нашли ответа на свой вопрос или нуждаетесь в консультации наших технических специалистов, оставьте запрос для службы поддержки"})
	viper.SetDefault("telegram.main_image", "")
	viper.SetDefault("telegram.main_buttons", []Button{{Text: "✏️  Задать вопрос", Data: "ask_question"}})

	if err := viper.ReadInConfig(); err != nil {
		if strings.Contains(err.Error(), "Not Found in") {
			fmt.Println("Config file not found; ignore error if running in CI/CD")
		} else if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Config file not found; ignore error if running in CI/CD")
		} else {
			panic(err)
		}
	}

	var cfg TgBotConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}

	fmt.Println("Config loaded successfully")

	return &cfg
}
