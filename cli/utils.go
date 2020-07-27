package cli

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func showSuccess(msg string, args ...interface{}) {
	fmt.Printf(msg+"\n", args...)
}

func showError(fields logrus.Fields, msg string, args ...interface{}) {
	logrus.WithFields(fields).Errorf(msg, args...)
}
