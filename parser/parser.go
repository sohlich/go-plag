package parser

import (
	"encoding/json"
	"io"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
)

//Map with plugin paths
var plugins map[string]string = map[string]string{
	"java": "plugin/JPlag.jar",
}

//Tokenize document via java plugins
func TokenizeContent(content, lang string) ([]uint32, error) {
	reader := strings.NewReader(content)
	out, err := execJavaPlugin(reader, lang)
	if err != nil {
		return nil, err
	}
	decodedDoc := &NGramDoc{}
	decoder := json.NewDecoder(strings.NewReader(out))
	decoder.Decode(decodedDoc)

	//Create hashes from token strings
	hashes := make([]uint32, 0)
	for _, ngram := range decodedDoc.NGrams {
		hashes = append(hashes, hash(ngram))
	}

	log.Infof("Array of hashes\n%s", hashes)

	//Run through winnowing
	winnowing := Winnowing{4}
	fp, err := winnowing.processTokensToFingerPrint(hashes)
	log.Infof("Array of hashes\n%s", fp.FingerPrint)

	return fp.FingerPrint, err
}

//Executes java plugin based
//on assignment language
func execJavaPlugin(input io.Reader, pluginLanguage string) (string, error) {
	path := plugins[pluginLanguage]
	if path == "" {
		return path, &NoSuchPluginError{pluginLanguage}
	}
	log.Debugf("Path for plugin %s", path)
	subProcess := exec.Command("java", "-jar", path)
	subProcess.Stdin = input
	bs, err := subProcess.Output()
	if err != nil {
		log.Error(err)
		return "", err
	}
	return string(bs), nil
}
