package main

import (
	"context"
	"warships/pkg/api"
	"warships/pkg/game"
)

func main() {
	c := make(chan api.GameStatus)
	s := make(chan string)
	state := make(chan api.GameState)

	ctx := context.Background()
	app := game.NewApp(c, s, state)
	app.Menu(ctx)
}
