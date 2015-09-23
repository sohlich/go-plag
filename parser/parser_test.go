package parser

import (
	"encoding/json"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadPlugins(t *testing.T) {
	files, _ := readFilesFromDir("../plugin/")

	for _, file := range files {
		log.Println(loadPlugin(file))
	}
}

//Fake testing file for JAVA
var fakeFile = `
import com.google.gson.Gson;
import cz.fai.utb.lang.api.ParseResultWrapper;
import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;

/**
 *
 * @author radek
 */
public class CmdMain {

    private static final JavaProcessor processor = new JavaProcessor();

    public static void main(String[] args) {
        StringBuilder inputAppender = new StringBuilder();
        try {
            BufferedReader br
                    = new BufferedReader(new InputStreamReader(System.in));
            String input;
            while ((input = br.readLine()) != null) {
                inputAppender.append(input);
            }

            ParseResultWrapper wrapper = processor.parseSource(inputAppender.toString());
            System.out.println(new Gson().toJson(wrapper));

        } catch (IOException io) {

        }
    }
}
`
