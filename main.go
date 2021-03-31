package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/andygrunwald/go-jira"
	"gopkg.in/yaml.v2"
)

type Config struct {
	URL              string `yaml:"url"`
	Username         string `yaml:"username"`
	Password         string `yaml:"password"`
	BranchNameRegExp string `yaml:"replace_regexp"`
	Format           string `yaml:"format"`
}

func getConfig() Config {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Open(home + "/.config/jiraf/config.yml")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	var config Config
	decoder := yaml.NewDecoder(f)

	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}

	return config

}

func main() {
	config := getConfig()

	tp := jira.BasicAuthTransport{
		Username: config.Username,
		Password: config.Password,
	}

	jiraClient, _ := jira.NewClient(tp.Client(), config.URL)

	issue, _, err := jiraClient.Issue.Get(os.Args[1], &jira.GetQueryOptions{})
	if err != nil {
		log.Fatal(err)
	}

	branch := fmt.Sprintf(config.Format, issue.Key, beatifySummary(issue.Fields.Summary, config.BranchNameRegExp))
	fmt.Println(branch)
}

func beatifySummary(s string, regex string) string {
	reg, err := regexp.Compile(regex)
	if err != nil {
		log.Fatal(err)
	}

	processedString := reg.ReplaceAllString(strings.ReplaceAll(s, " ", "-"), "")

	return processedString
}
