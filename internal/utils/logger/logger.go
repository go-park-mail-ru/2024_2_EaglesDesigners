package logger

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

type contextKey string

const RequestIDKey contextKey = "requestID"

var Log *logrus.Logger

func init() {
	Log = logrus.New()
	Log.SetFormatter(&CustomFormatter{})
}

func LoggerWithCtx(ctx context.Context, log *logrus.Logger) *logrus.Entry {
	reqID, ok := ctx.Value(RequestIDKey).(string)
	if !ok {
		reqID = "errors"
	}

	pc, _, _, _ := runtime.Caller(1)
	funcPath := runtime.FuncForPC(pc).Name()

	parts := strings.Split(funcPath, "/")

	var funcName string
	var feature string

	if len(parts) > 2 {
		funcName = parts[len(parts)-1]
		feature = parts[len(parts)-2]
	}

	fields := logrus.Fields{
		"request_id": reqID[:6],
		"function":   funcName,
		"feature":    feature,
	}

	entry := log.WithFields(fields)

	return entry
}

type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var levelColor string
	switch entry.Level {
	case logrus.InfoLevel:
		levelColor = "\033[34m" // Синий
	case logrus.WarnLevel:
		levelColor = "\033[33m" // Желтый
	case logrus.ErrorLevel:
		levelColor = "\033[31m" // Красный
	case logrus.FatalLevel, logrus.PanicLevel:
		levelColor = "\033[35m" // Фиолетовый
	default:
		levelColor = "\033[0m" // Без цвета
	}

	timestamp := entry.Time.Format("2006-01-02T15:04:05")

	requestID := entry.Data["request_id"]
	feature := entry.Data["feature"]
	function := entry.Data["function"]

	level := fmt.Sprintf("%s[%s]%s", levelColor, strings.ToUpper(entry.Level.String()), "\033[0m")
	feature = fmt.Sprintf("%s[%s]%s", "\033[33m", feature, "\033[0m")
	function = fmt.Sprintf("%s[%s]%s", "\033[36m", function, "\033[0m")

	log := fmt.Sprintf("%s[%s][%s]%s%s %s\n",
		level,
		timestamp,
		requestID,
		feature,
		function,
		entry.Message,
	)

	return []byte(log), nil
}
