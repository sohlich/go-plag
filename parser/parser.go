package parser

import (
	"encoding/json"
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

func TokenizeContent(content, lang string) (*NGramDoc, error) {
	reader := strings.NewReader(content)
	out, err := execJavaPlugin(reader, lang)
	if err != nil {
		return nil, err
	}
	decodedDoc := &NGramDoc{}
	decoder := json.NewDecoder(strings.NewReader(out))
	decoder.Decode(decodedDoc)

	return decodedDoc, nil
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
