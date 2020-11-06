package helpers

import (
	"os"

	"github.com/op/go-logging"
)

var Log = logging.MustGetLogger("velobike")

func init() {
	loggingBackend := logging.NewLogBackend(os.Stdout, "", 0)
	format := logging.MustStringFormatter(
		`%{color}%{time:2006-01-02 15:04:05.000} %{level:.4s}%{color:reset} %{message}`,
	)

	logging.SetBackend(logging.NewBackendFormatter(loggingBackend, format))
}
