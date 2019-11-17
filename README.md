# Lugo - The Dummies Go

[![GoDoc](https://godoc.org/github.com/lugobots/the-dummies-go?status.svg)](https://godoc.org/github.com/lugobots/the-dummies-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/lugobots/the-dummies-go)](https://goreportcard.com/report/github.com/lugobots/the-dummies-go)

The Dummies Go is a [Go](http://golang.org/) implementation of a player (bot) for [Lugo](https://lugobots.dev) game.
This bot was made using the [Go Client Player](https://github.com/lugobots/client-player-go) for Go.

As this name suggest, _The Dummies_ are not that smart, but they may play well enough to help you to test your bot.

### Requirements

0. Docker >= 18.03 (https://docs.docker.com/install/)
0. Docker Compose >= 1.21 (https://docs.docker.com/compose/install/)
0. Go Lang >= 1.12 (https://golang.org/doc/install)

### Usage 

You have two ways to make The Dummies play, they are described below.
 
#### Option A - Running them in containers (no Git Clone needed)

Download the [Docker compose file](https://raw.githubusercontent.com/lugobots/the-dummies-go/master/docker-compose.yml) that starts
the server along with 11 instances of _The Dummies_ bot.

Start the set of containers:
```
TEAM_IMAGE=lugobots/the-dummies-go TEAM_PLACE=away docker-compose up
```

It will start them as **away** team (defined by the env variable `TEAM_PLACE`).
Now you may start your bot to play against _The Dummies_ 

#### Option B  (recommended for developing environment because the startup is a faster)

If you are working in your bot, and you would like to play against The Dummies several times to test your bot, I recommend
you having a copy of The Dummies in you machine, so it will startup faster than running them as container. 

1. Clone the repository to your machine
2. Start the game server
   ```
   docker run -p 8080:8080  lugobots/server:v1.1 play --dev-mode
   ```

, and then, you may execute the script `./play.sh [home|away]` in that directory when you want to start the team.

### The Dummies vs The Dummies

If you have no other team to play against _The Dummies_ or if you are just curious to watch a Lugo match,
you may start a game using _The Dummies_ as the Home and Away teams.

0. Download the [Demo Docker compose file](https://raw.githubusercontent.com/lugobots/the-dummies-go/master/docker-compose-demo.yml)
0. Execute the command bellow:
    ```bash
    HOME_TEAM=lugobots/the-dummies-go \
    AWAY_TEAM=lugobots/the-dummies-go \
    docker-compose -f docker-compose-demo.yml up
    ```
0. Watch the game in the browser at the address `http://localhost:8080`

**Important**: You probably want to remove that bunch of containers from your environment later. So, execute the command below:
```bash
HOME_TEAM=lugobots/the-dummies-go \
AWAY_TEAM=lugobots/the-dummies-go \
docker-compose -f docker-compose-demo.yml down
```
