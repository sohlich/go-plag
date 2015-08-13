package parser

import (
	"fmt"
	"io"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
)

var plugins map[string]string = map[string]string{
	"java": "plugin/JPlag.jar",
}

type NoSuchPluginError struct {
	Lang string
}

func (err *NoSuchPluginError) Error() string {
	return fmt.Sprintf("Plugin for language %s does not exist", err.Lang)
}

func TokenizeContent(content, lang string) (string, error) {
	reader := strings.NewReader(content)
	return execJavaPlugin(reader, lang)
}

func execJavaPlugin(input io.Reader, pluginLanguage string) (string, error) {
	path := plugins[pluginLanguage]
	if path == "" {
		return path, &NoSuchPluginError{pluginLanguage}
	}
	log.Info(path)
	subProcess := exec.Command("java", "-jar", path)
	subProcess.Stdin = input
	bs, err := subProcess.Output()
	if err != nil {
		log.Error(err)
		return "", err
	}
	return string(bs), nil
}
