package main

import (
	"github.com/lugobots/lugo4go/v2"
	"github.com/lugobots/lugo4go/v2/coach"
	"github.com/lugobots/the-dummies-go/v2/bot"
	"log"
	"os"
	"os/signal"
)

func main() {
	var err error
	// DefaultBundle is a shot cut for stuff that usually we define in init functions
	config, logger, err := lugo4go.DefaultBundle()
	if err != nil {
		log.Fatalf("could not init default config or logger: %s", err)
	}

	// Creating a bot to play
	myBot, err := bot.NewBot(config, logger)
	if err != nil {
		logger.Fatalf("did not connected to the gRPC server at '%s': %s", config.GRPCAddress, err)
	}

	config.InitialPosition = myBot.InitialPosition

	// open the connection to the server
	ctx, client, err := lugo4go.NewClient(config)
	if err != nil {
		logger.Fatalf("did not connected to the gRPC server at '%s': %s", config.GRPCAddress, err)
	}
	// defining the bot as the "decider" interface to be used by the Turn Handler
	client.OnNewTurn(coach.DefaultTurnHandler(myBot, config, logger), logger)

	// keep the process alive
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	select {
	case <-signalChan:
		logger.Warnf("got interruption signal")
		if err := client.Stop(); err != nil {
			logger.Errorf("error stopping the player client: %s", err)
		}
	case <-ctx.Done():
		logger.Infof("player client stopped")
	}
	logger.Infof("process finished")
}
