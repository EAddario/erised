package main

import (
	"time"

	"github.com/rs/zerolog/log"
)

func elapsedTime(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Debug().Msg(name + " ran for " + elapsed.Round(time.Second).String())
}
