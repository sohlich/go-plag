package parser

import (
	"io"
	"os/exec"

	log "github.com/Sirupsen/logrus"
)

var plugins map[string]string = map[string]string{
	"java": "../plugin/JPlag.jar",
}

func execJavaPlugin(input io.Reader, pluginLanguage string) (string, error) {
	path := plugins[pluginLanguage]
	subProcess := exec.Command("java", "-jar", path)
	subProcess.Stdin = input
	bs, err := subProcess.Output()
	if err != nil {
		log.Error(err)
		return "", err
	}
	return string(bs), nil
}
