package main

import "server/internal/app"

func main() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}
