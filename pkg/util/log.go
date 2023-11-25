package util

import (
	"os"

	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
)

var Logger = slog.Make(sloghuman.Sink(os.Stdout))
