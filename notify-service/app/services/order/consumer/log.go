package consumer

import (
	"fmt"
	"io"
	"os"

	"github.com/cg917658910/fzkj-wallet/notify-service/lib/log"
	"github.com/cg917658910/fzkj-wallet/notify-service/lib/log/base"
)

var (
	_logFile    io.WriteCloser
	_errLogFile io.WriteCloser
	logger      base.MyLogger
	errLogger   base.MyLogger
	_logPath         = "./tmp/log/notify/"
	_logname         = "consumer.log"
	_errlogname      = "consumer_error.log"
	_useStdout  bool = false
)

func init() {
	setupLogger()
}
func setupLogger() {
	_logFile = os.Stdout
	if !_useStdout {
		logout, err := getLogOut(_logname)
		if err != nil {
			fmt.Printf("get logout %s err: %s\n", _logname, err)
		}
		if logout != nil {
			_logFile = logout
		}
	}
	logger = log.DLoggerWithWriter(_logFile)

	_errLogFile = os.Stdout
	errlogout, err := getLogOut(_errlogname)
	if err != nil {
		fmt.Printf("get logout %s err: %s\n", _errlogname, err)
	}
	if errlogout != nil {
		_errLogFile = errlogout
	}
	errLogger = log.DLoggerWithWriter(_errLogFile)
}

func cleanupLogger() {
	if _logFile != nil {
		_logFile.Close()
	}
	if _errLogFile != nil {
		_errLogFile.Close()
	}
}

func getFileAndMakeDir(path, filename string) (logout io.WriteCloser, err error) {
	logFilename := path + filename
	if _, err = os.Stat(path); err != nil {
		err = os.MkdirAll(path, 0711)
		if err != nil {
			return
		}
	}
	logout, err = os.OpenFile(logFilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	return logout, nil

}

func getLogOut(filename string) (log io.WriteCloser, err error) {
	return getFileAndMakeDir(_logPath, filename)
}
