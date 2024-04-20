package fvtime

import (
	"log/slog"
	"os"
	"time"
)

const FV_NOW_ENV_NAME = "FV_NOW"

func New(now time.Time) time.Time {
	if os.Getenv(FV_NOW_ENV_NAME) == "" {
		return now
	}

	time, err := time.Parse(time.DateOnly, os.Getenv(FV_NOW_ENV_NAME))
	if err != nil {
		slog.Error(err.Error())
		return now
	}

	return time
}
