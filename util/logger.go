package util

import (
	loggerUtil "github.com/koko990/logcollector/util/logger"
)
var Logger loggerUtil.NewLog
func init() {
	Logger.LogRegister(loggerUtil.LevelDebug)
}
