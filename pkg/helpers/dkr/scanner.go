package dkr

import (
	"bufio"
	"io"

	"github.com/sirupsen/logrus"
	"github.com/vitorfhc/bob/pkg/docker/outputs"
)

// ScanBody reads the body of the response of a Docker command. The output will be
// the last line at the end. If logger is not nil, the output will be logged.
func ScanBody(body io.ReadCloser, output outputs.Output, logger *logrus.Entry) {
	var lastLine string
	scanner := bufio.NewScanner(body)
	for scanner.Scan() {
		lastLine = scanner.Text()
		err := output.LoadFromJSON(lastLine)
		if err != nil {
			logger.Error(err)
		}
		if logger != nil {
			logger.Log(logrus.DebugLevel, output.String())
		}
	}
}
