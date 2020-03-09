package log

import (
	"io"
	"os"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
)

func out(cfg Config) io.Writer {
	if len(cfg) == 0 {
		return os.Stdout
	}

	if strings.Contains(strings.ToLower(strings.Join(os.Args, "-")), "go-build") {
		return os.Stdout
	}

	rs, err := rotatelogs.New(cfg+".%Y_%m_%d",
		rotatelogs.WithLinkName(cfg),
		rotatelogs.WithMaxAge(time.Duration(31*24)*time.Hour))
	if err != nil {
		panic(errors.Wrap(err, "Set rotate log error"))
	}
	return rs
}
