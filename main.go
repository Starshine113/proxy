package main

import (
	"flag"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/Starshine113/proxy/bot"
	"github.com/Starshine113/proxy/commands"
	"github.com/Starshine113/proxy/db"
	"github.com/Starshine113/proxy/proxy"
	"github.com/Starshine113/proxy/router"
	"github.com/Starshine113/proxy/structs"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

var proxybot *bot.Bot
var database *db.Db
var config *structs.BotConfig
var sugar *zap.SugaredLogger

func main() {
	dg := initialize()

	p := &proxy.Proxy{Bot: proxybot, Session: dg}
	r := router.NewRouter(proxybot, p)
	dg.AddHandler(p.ReactionAdd)

	commands.Init(r)

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages | discordgo.IntentsGuilds | discordgo.IntentsDirectMessages | discordgo.IntentsGuildMessageReactions | discordgo.IntentsDirectMessageReactions)

	err := dg.Open()
	if err != nil {
		panic(err)
	}

	// Defer this to make sure that things are always cleanly shutdown even in the event of a crash
	defer func() {
		dg.Close()
		sugar.Infof("Disconnected from Discord.")
		database.Pool.Close()
		sugar.Infof("Closed database connection.")
	}()

	sugar.Infof("Connected to Discord. Press Ctrl-C or send an interrupt signal to stop.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	sugar.Infof("Interrupt signal received. Shutting down...")
}

func initialize() *discordgo.Session {
	token := flag.String("token", "", "Override the token in config.toml")
	databaseURL := flag.String("db", "", "Override the database URL in config.toml")
	flag.Parse()

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.RedirectStdLog(logger)
	sugar = logger.Sugar()

	if _, err := os.Stat("config.toml"); os.IsNotExist(err) {
		sampleConf, err := ioutil.ReadFile("config.sample.toml")
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile("config.toml", sampleConf, 0644)
		if err != nil {
			panic(err)
		}
		sugar.Errorf("config.toml was not found, created sample configuration.")
		os.Exit(1)
		return nil
	}
	configFile, err := ioutil.ReadFile("config.toml")
	if err != nil {
		panic(err)
	}
	err = toml.Unmarshal(configFile, &config)
	sugar.Infof("Loaded configuration file.")

	if *token != "" {
		config.Auth.Token = *token
	}
	if *databaseURL != "" {
		config.Auth.DatabaseURL = *databaseURL
	}
	if os.Getenv("CB_DB_URL") != "" {
		config.Auth.DatabaseURL = os.Getenv("CB_DB_URL")
	}

	database, err = db.Init(config, sugar)
	if err != nil {
		panic(err)
	}
	sugar.Infof("Loaded database")

	dg, err := discordgo.New("Bot " + config.Auth.Token)
	if err != nil {
		panic(err)
	}

	proxybot = bot.New(database, dg, config, sugar)

	return dg
}
