package keybaser

import "github.com/sirupsen/logrus"

func newDefaultLogger() logrus.FieldLogger {
	return logrus.New()
}
