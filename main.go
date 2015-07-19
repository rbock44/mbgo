package main

import (
	"fmt"
	//"flag"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"regexp"
	"bufio"
)

func main() {
	//flag.Parse()

	includes := parseIncludes("mbconfig.json")

	mbContent := loadFile("mbconfig.json")
	for _, include := range includes {
		mbContent = strings.Replace(mbContent, "\"INCLUDE="+ include.Key + "\"", include.Value, 1)
	}

	var result map[string]interface{}
	err := json.Unmarshal([]byte(mbContent), &result)
	if err != nil {
		fmt.Printf("invalid json err=<%v>", err)
		os.Exit(1)
	}

	ioutil.WriteFile("mb.json", []byte(mbContent), 0644)
}

type Include struct {
	Key	string
	Value	string
}

func parseIncludes(filename string) []Include {
	includes := make([]Include, 10)

	file, err :=  os.Open(filename) 
	if err != nil {
		fmt.Printf("cannot open config file=<%s>", filename)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "\"INCLUDE=") {
			exp, _ := regexp.Compile(".*INCLUDE=(.*)\"")
			include := exp.FindStringSubmatch(line)[1]
			content := loadFile(include)
			escapedContent := escapeContent(content)
			entry := Include { include, escapedContent }
			includes = append(includes, entry)
		}
	}
	return includes
}

func loadFile(path string) string {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("error cannot open file=<%s> error=<%v>", path, err)
		os.Exit(1)
	}
	return string(content[:])
}

func escapeContent(content string) string {
	escapedContent, err := json.Marshal(content)
	if err != nil {
		fmt.Printf("error cannot encode content=<%s> to json err=<%v>\n", content, err)
		os.Exit(1)
	}
	return string(escapedContent[:])
}
