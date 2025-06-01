package timezone

import (
	"log/slog"
	"os"
	"time"
)

var (
	defaultLocation = time.UTC
)

func Init() (*time.Location, error) {
	tz := os.Getenv("TMZ")
	slog.Info(tz)
	if tz == "" {
		return defaultLocation, nil
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return nil, err
	}
	defaultLocation = loc
	return loc, nil
}

func Get() *time.Location {
	return defaultLocation
}

func Now() time.Time {
	now := time.Now().In(defaultLocation)
	slog.Info("Time with location",
		slog.String("time", now.String()),
		slog.String("location", defaultLocation.String()),
	)
	return now
}
