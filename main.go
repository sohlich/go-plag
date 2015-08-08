package main

import (
	"os"
	"os/exec"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

func main() {
	file, fileErr := os.Open("/home/radek/Projekty/Java/Plag/PlagiarismDetection/JPlag/src/main/java/cz/fai/utb/java/api/CmdMain.java")
	if fileErr != nil {
		log.Error(fileErr)
	}
	subProcess := exec.Command("java", "-jar", "plugin/JPlag.jar")
	subProcess.Stdin = file
	bs, err := subProcess.Output()
	if err != nil {
		log.Error(err)
	}
	log.Println(string(bs))

}

func execJavaPlugin(input io.Reader, pluginPath string) (string, err) {
	subProcess := exec.Command("java", "-jar", "plugin/JPlag.jar")
	subProcess.Stdin = file
	bs, err := subProcess.Output()
	if err != nil {
		log.Error(err)
	}
	log.Println(string(bs))
}
