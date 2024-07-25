package logging

import (
	"github.com/sirupsen/logrus"
)

func GetLogger() *logrus.Logger {
	return logrus.StandardLogger()
}
