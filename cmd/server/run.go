package main

import (
	"library-api/cmd/server/bootstrap"
	"library-api/kit/log"
)

func main() {
	if err := bootstrap.Run(); err != nil {
		log.FatalErr("error running server", err)
	}
}
