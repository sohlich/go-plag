package parser

import (
	"encoding/json"
	"strings"
	"testing"

	log "github.com/Sirupsen/logrus"
)

func TestParseDocument(t *testing.T) {
	input := strings.NewReader(fakeFile)
	pluginLang := "java"
	output, err := execJavaPlugin(input, pluginLang)
	if err != nil {
		t.Error(err)
	}
	decoder := json.NewDecoder(strings.NewReader(output))

	decodedDoc := NGramDoc{}
	decoder.Decode(&decodedDoc)

	log.Info(decodedDoc)
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
