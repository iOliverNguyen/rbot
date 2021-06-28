package l

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// production | development
const EnvKey = "LOGGING"

// debug | info | warn | error | dpanic | panic | fatal
const EnvLevel = "LOGGING_LEVEL"

const EnvProd = "production"

var zapCfg zap.Config

func init() {
	if isProd() {
		zapCfg = zap.NewProductionConfig()

	} else {
		zapCfg = zap.NewDevelopmentConfig()
	}

	lvl, ok := parseLevel(os.Getenv(EnvLevel))
	if ok {
		zapCfg.Level.SetLevel(lvl)
	}
}

func isProd() bool {
	envLog := os.Getenv(EnvKey)
	return strings.ToLower(envLog) == EnvProd
}

func parseLevel(s string) (lvl zapcore.Level, ok bool) {
	err := lvl.Set(s)
	return lvl, err == nil
}
