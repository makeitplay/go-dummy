# MakeItPlay - The Dummies Go

[![GoDoc](https://godoc.org/github.com/makeitplay/the-dummies-go?status.svg)](https://godoc.org/github.com/makeitplay/the-dummies-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/makeitplay/the-dummies-go)](https://goreportcard.com/report/github.com/makeitplay/the-dummies-go)

The Dummies Go is a [Go](http://golang.org/) implementation of a player (bot) for [MakeItPlay football](http://www.makeitplay.ai/football) game.
This bot was made using the [Client Player](https://github.com/makeitplay/client-player-go) for Go.

As this name suggest, **The Dummies** are not that smart, but they may play well enough to help you to test your bot.

### Requirements

0. Docker >= 18.03 (https://docs.docker.com/install/)
0. Docker Compose >= 1.21 (https://docs.docker.com/compose/install/)
0. Go Lang >= 1.10 (https://golang.org/doc/install)

### Usage 

You have two ways to make The Dummies play, they are described below.
However, in both ways we will need to start the game server first:

```bash
    docker run -p 8080:8080  makeitplay/football:1.0.0-alpha
```

Then you may start The Dummies **and** the other team.
 
#### Option A (no Git Clone needed)

You do not have to download The Dummies to make it play against your team.

If you only wish to play with them, you may them as a Docker container:

```bash
docker run makeitplay/the-dummies-go -team=[home|away] -number=[1-11] 
```

However, you will need to execute the command above 11 times (one for each player position). 
So, you may downloading [this script](./start-team-container.sh) and run the command below instead:

```bash
./start-team-container.sh makeitplay/the-dummies-go [home|away]
```

#### Option B  (recommended for developing environment because the startup is a little faster)

If you are working in your bot, and you would like to play against The Dummies several times to test your bot, I recommend
you having a copy of The Dummies in you machine, so it will startup faster than running them as container. 

First, clone the project using the command:
```bash
go get github.com/makeitplay/the-dummies-go
```

Then, you may execute the script `play.sh [home|away]` in that directory always you want to start the team.

### The Dummies vs The Dummies

If you have no other team to play against **The Dummies** or if you are just curious to watch a Make It Play match,
you may start a game using The Dummies as the Home and Away teams.

0. Download the [Demo Docker compose file](./docker-compose-demo.yml)
0. Execute the command bellow:
    ```bash
    HOME_TEAM=makeitplay/the-dummies-go \
    AWAY_TEAM=makeitplay/the-dummies-go \
    docker-compose -f docker-compose-demo.yml up
    ```
0. Watch the game in the browser at the address `http://localhost:8080`

**Important**: You probably want to remove that bunch of containers from your environment later. So, execute the command below:
```bash
HOME_TEAM=makeitplay/the-dummies-go \
AWAY_TEAM=makeitplay/the-dummies-go \
docker-compose -f docker-compose-demo.yml down
```
