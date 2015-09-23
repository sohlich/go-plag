package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"

	log "github.com/Sirupsen/logrus"
)

type Plugin struct {
	Path       string
	Language   string
	Extensions []string
	FileFilter map[string]bool
}

//Returnd initalized Plugin struct pointer
//with created map
func NewPlugin() *Plugin {
	return &Plugin{FileFilter: make(map[string]bool)}
}

var pluginMap map[string]*Plugin = map[string]*Plugin{}

//Tokenize document via java plugins
func TokenizeContent(content, lang string) ([]uint32, map[string]int, error) {
	log.Debugf("Starting to tokenize")

	reader := strings.NewReader(content)
	out, err := execJavaPlugin(reader, lang)
	if err != nil {
		return nil, nil, err
	}
	decodedDoc := &NGramDoc{}
	decoder := json.NewDecoder(strings.NewReader(out))
	decoder.Decode(decodedDoc)

	//Create hashes from token strings
	hashes := make([]uint32, 0)
	for _, ngram := range decodedDoc.NGrams {
		hashes = append(hashes, hash(ngram))
	}

	//Run through winnowing
	winnowing := Winnowing{4} //initialize winnowing for window n=4
	fp, err := winnowing.processTokensToFingerPrint(hashes)
	outM := make(map[string]int)
	for _, hash := range fp.FingerPrint {
		outM[fmt.Sprint(hash)] += 1
	}

	return fp.FingerPrint, outM, err
}

//Executes java plugin based
//on assignment language
func execJavaPlugin(input io.Reader, pluginLanguage string) (string, error) {
	pluginKey := strings.ToLower(pluginLanguage)
	plugin, ok := pluginMap[pluginKey]
	if !ok {
		return "", &NoSuchPluginError{pluginKey}
	}

	log.Debugf("Executing plugin %s", plugin.Path)

	cmd := execJava(plugin.Path, "parse")
	cmd.Stdin = input
	bs, err := cmd.CombinedOutput()
	if err != nil {
		cmd.Process.Signal(os.Kill)
		log.Errorf(string(bs))
		return "", err
	}
	return string(bs), nil
}

func execJava(path string, command string) *exec.Cmd {
	subProcess := exec.Command("java", "-jar", path, command)
	return subProcess
}

//Load plugins from given directory
func LoadPlugins(dirPath string) {
	//TODO better handle errors
	filePaths, err := readFilesFromDir(dirPath)
	if err != nil {
		log.Errorln("Cannot load any plugin")
	}

	for _, plugPath := range filePaths {
		plugin, loadErr := loadPlugin(plugPath)
		if loadErr == nil {
			log.Infof("Loaded plugin for %c", plugin.Language)
			pluginKey := strings.ToLower(plugin.Language)
			pluginMap[pluginKey] = plugin
		} else {
			log.Errorf("Cannot load plugin %s because of %v", plugPath, loadErr)
		}
	}

}

//Reads the content of given folder and returns the
//string slice with file paths
func readFilesFromDir(dirPath string) ([]string, error) {
	dir, err := os.Open(dirPath)

	if err != nil {
		return nil, err
	}

	files, err := dir.Readdir(-1)

	if err != nil {
		return nil, err
	}
	filePaths := make([]string, 0)
	for _, file := range files {
		filePath := path.Join(dirPath, file.Name())
		filePaths = append(filePaths, filePath)
	}
	return filePaths, nil
}

//Loads the java Plag plugin, the plugin must return the
//info JSON with identifier and info from plugin jar
func loadPlugin(pluginPath string) (*Plugin, error) {
	subProcess := execJava(pluginPath, "info")
	bs, err := subProcess.Output()

	if err != nil {
		return nil, err
	}

	plugin := NewPlugin()
	decoder := json.NewDecoder(strings.NewReader(string(bs)))
	err = decoder.Decode(plugin)

	if err != nil {
		return nil, err
	}

	plugin.Path = pluginPath

	for _, ext := range plugin.Extensions {
		plugin.FileFilter[ext] = true
	}

	return plugin, nil
}

func GetLangFileFilter(lang string) map[string]bool {
	return pluginMap[lang].FileFilter
}

func GetSupportedLangs() []string {
	langs := make([]string, 0)
	for lang, _ := range pluginMap {
		langs = append(langs, lang)
	}

	return langs
}
