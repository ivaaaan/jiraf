package main

import (
	"fmt"
	"log"
	"os"

	"github.com/andygrunwald/go-jira"
	"gopkg.in/yaml.v2"
)

type Config struct {
	URL              string              `yaml:"url"`
	Username         string              `yaml:"username"`
	Password         string              `yaml:"password"`
	BranchNameRegExp string              `yaml:"replace_regexp"`
	Format           string              `yaml:"format"`
	Pipeline         map[string][]string `yaml:"pipeline"`
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

	branchName := issue.Fields.Summary

	for pipeName, args := range config.Pipeline {
		if pipeFunc, ok := PipelineMap[pipeName]; ok {
			branchName, err = pipeFunc(branchName, args...)
			if err != nil {
				log.Fatal(err)
			}
		}

	}

	branch := fmt.Sprintf(config.Format, issue.Key, branchName)
	fmt.Println(branch)
}
