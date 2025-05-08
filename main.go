package main

import "github.com/MRegterschot/trackmania-server-fm/app"

func main() {
	err := app.SetupAndRunApp()
	if err != nil {
		panic(err)
	}
}
