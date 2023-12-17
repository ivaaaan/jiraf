package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/andygrunwald/go-jira"
	"gopkg.in/yaml.v2"
)

var help = `Jiraf finds an issue in Jira by key and generates a git branch name from its summary.

Usage:
	jiraf <jira ticket id>

Options:
	--help, Print this message
`

type Config struct {
	URL              string              `yaml:"url"`
	Username         string              `yaml:"username"`
	Password         string              `yaml:"password"`
	BranchNameRegExp string              `yaml:"replace_regexp"`
	Format           string              `yaml:"format"`
	Pipeline         map[string][]string `yaml:"pipeline"`
}

func defaultConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", nil
	}

	return home + "/.config/jiraf/config.yml", nil
}

func getConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	decoder := yaml.NewDecoder(f)

	var config Config
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func main() {
	configPath, err := defaultConfigPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting a default config path: %v", err)
		os.Exit(1)
	}

	config, err := getConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing config: %v", err)
		os.Exit(1)
	}

	tp := jira.BasicAuthTransport{
		Username: config.Username,
		Password: config.Password,
	}

	jiraClient, err := jira.NewClient(tp.Client(), config.URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error initializing a Jira client: %v", err)
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, "expected single argument, a ticket ID (e.g., PROJECT-123)")
		os.Exit(1)
	}

	if os.Args[1] == "--help" {
		fmt.Fprint(os.Stdout, help)
		os.Exit(0)
	}

	issue, _, err := jiraClient.Issue.Get(os.Args[1], &jira.GetQueryOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting a Jira issue: %v", err)
		os.Exit(1)
	}

	branchName := issue.Fields.Summary

	var errors []string
	for pipeName, args := range config.Pipeline {
		if pipeFunc, ok := PipelineMap[pipeName]; ok {
			branchName, err = pipeFunc(branchName, args...)
			if err != nil {
				errors = append(errors, fmt.Sprintf("%s function returned an error: %v", pipeName, err))
			}
		}
	}

	if len(errors) > 0 {
		fmt.Fprintf(os.Stderr, "one or more functions returned errors: %s", strings.Join(errors, ","))
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, config.Format, issue.Key, branchName)
}
