package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
)

var UnsupportedLangError = errors.New("Language not supported")

type Tokenizer interface {
	ProcessFile(input io.Reader) (string, error)
	ExtensionsFilter() map[string]bool
}

type Plugin struct {
	Path       string
	Language   string
	Extensions []string
	FileFilter map[string]bool
}

func (p *Plugin) ProcessFile(input io.Reader) (string, error) {
	cmd := exec.Command("java", "-jar", p.Path, "parse")
	cmd.Stdin = input
	bs, err := cmd.CombinedOutput()
	if err != nil {
		cmd.Process.Signal(os.Kill)
		Log.Errorf(string(bs))
		return "", err
	}
	return string(bs), nil
}

//Return the map of supported file extensions
func (p *Plugin) ExtensionsFilter() map[string]bool {
	return p.FileFilter
}

var Log = log.StandardLogger()

func SetLogger(logger *log.Logger) {
	Log = logger
}

//Returnd initalized Plugin struct pointer
//with created map
func NewPlugin() *Plugin {
	return &Plugin{FileFilter: make(map[string]bool)}
}

type PluginMap struct {
	mutex     *sync.Mutex
	pluginMap map[string]Tokenizer
}

func NewPluginMap() *PluginMap {
	pluginMap := &PluginMap{
		new(sync.Mutex),
		map[string]Tokenizer{},
	}

	return pluginMap
}

func (p *PluginMap) PutPlugin(lang string, tok Tokenizer) {
	p.mutex.Lock()
	p.pluginMap[lang] = tok
	p.mutex.Unlock()
}

func (p *PluginMap) GetTokenizer(lang string) (Tokenizer, bool) {
	p.mutex.Lock()
	tok, ok := p.pluginMap[lang]
	p.mutex.Unlock()
	return tok, ok
}

func (p *PluginMap) GetLangs() []string {
	p.mutex.Lock()
	langs := make([]string, 0)
	for lang, _ := range p.pluginMap {
		langs = append(langs, lang)
	}

	p.mutex.Unlock()
	return langs
}

var pluginMap = NewPluginMap()

//Tokenize document via java plugins
func TokenizeContent(content, lang string) ([]uint32, map[string]int, error) {
	Log.Debugf("Starting to tokenize")

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
	plugin, ok := pluginMap.GetTokenizer(pluginKey)
	if !ok {
		return "", &NoSuchPluginError{pluginKey}
	}

	Log.Debugf("Executing plugin %v", plugin)
	return plugin.ProcessFile(input)
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
		Log.Errorln("Cannot load any plugin")
	}

	for _, plugPath := range filePaths {
		plugin, loadErr := loadPlugin(plugPath)
		if loadErr == nil {
			Log.Infof("Loaded plugin for %c", plugin.Language)
			pluginKey := strings.ToLower(plugin.Language)
			pluginMap.PutPlugin(pluginKey, plugin)
		} else {
			Log.Errorf("Cannot load plugin %s because of %v", plugPath, loadErr)
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
	defer dir.Close()

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

func GetLangFileFilter(lang string) (map[string]bool, error) {
	plugin, ok := pluginMap.GetTokenizer(lang)
	if ok {
		return plugin.ExtensionsFilter(), nil
	} else {
		return nil, UnsupportedLangError
	}

}

func GetSupportedLangs() []string {
	return pluginMap.GetLangs()
}

func IsSupportedLang(lang string) bool {
	_, ok := pluginMap.GetTokenizer(lang)
	return ok
}
