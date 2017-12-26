package main

import (
	"flag"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/run"

	"github.com/coinlet/coinlet/internal/ticker"
)

func runTicker(args []string) error {
	flagset := flag.NewFlagSet("ticker", flag.ExitOnError)
	var (
		debug = flagset.Bool("debug", false, "debug logging")
	)
	flagset.Usage = usageFor(flagset, "coinlet ticker [flags]")
	if err := flagset.Parse(args); err != nil {
		return err
	}

	// Logging.
	var logger log.Logger
	{
		logLevel := level.AllowInfo()
		if *debug {
			logLevel = level.AllowAll()
		}
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = level.NewFilter(logger, logLevel)
	}

	var g run.Group
	{
		t := ticker.NewTicker(log.With(logger, "component", "ticker"))
		g.Add(func() error {
			t.Run()
			return nil
		}, func(error) {
			t.Stop()
		})
	}
	{
		cancel := make(chan struct{})
		g.Add(func() error {
			return interrupt(cancel)
		}, func(error) {
			close(cancel)
		})
	}
	return g.Run()

	return nil
}
