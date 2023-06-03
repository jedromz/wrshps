package main

import (
	"context"
	"warships/pkg/game"
)

func main() {
	c := make(chan game.GameEvent)
	s := make(chan string)
	t := make(chan int)
	ctx := context.Background()
	app := game.Setup(c, s, t)
	app.Run(ctx)
}
