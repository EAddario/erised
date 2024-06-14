package main

import (
	"flag"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func setupFlags(f *flag.FlagSet) {
	log.Debug().Msg("entering setupFlags")

	f.Usage = func() {
		fmt.Println("Simple http server to test arbitrary responses (" + version + ")")
		fmt.Println("Usage examples at https://github.com/EAddario/erised")
		fmt.Println("\nerised [options]")
		fmt.Println("\nParameters:")
		flag.PrintDefaults()
		fmt.Println("\nHTTP Headers:")
		fmt.Println("X-Erised-Content-Type:\t\tSets the response Content-Type")
		fmt.Println("X-Erised-Data:\t\t\tReturns the same value in the response body")
		fmt.Println("X-Erised-Headers:\t\tReturns the value(s) in the response header(s). Values must be in a JSON array")
		fmt.Println("X-Erised-Location:\t\tSets the response Location when 300 ≤ X-Erised-Status-Code < 310")
		fmt.Println("X-Erised-Response-Delay:\tNumber of milliseconds to wait before sending response back to client")
		fmt.Println("X-Erised-Response-File:\t\tReturns the contents of file in the response body. If present, X-Erised-Data is ignored")
		fmt.Println("X-Erised-Status-Code:\t\tSets the HTTP Status Code")
		fmt.Println()
	}

	log.Debug().Msg("leaving setupFlags")
}
