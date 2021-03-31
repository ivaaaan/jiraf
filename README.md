# Jiraf

Jiraf finds an issue in Jira by key and generates a git branch name from its summary.

## Configuration

Configuration file stored in `~/.config/jiraf/config.yml` in the YAML format.

```yaml
url: <url to your Jira>
username: <username>
password: <Jira API token>
regexp: <regexp to replace characters from summary>
format: <format for fmt.Sprintf, the first argument is an issue key, the second one is generated summary>
```

### Example

```yaml
url:  <url to your Jira>
username: <username>
password: <Jira api token>
regexp: "[^a-zA-Z0-9-]+"
format: "%s_%s"
```

This configuration will generate a branch name like this: `ISSUE-1_Summary-of-your-issue`

## Usage

One way to integrate jiraf with git is to create a simple script like this:

```bash
#!/bin/bash

git checkout -b $(jiraf $2)
```

And put this script into `~/bin/git-cb` file (replace "cb" with preferred subcommand name). Now you can just run `git cb ISSUE-1`, which will create a new branch with a generated name by jiraf.

## TODO

[] Allow customizing generation of summary

[] Add a limit for the max amount of words in the branch name
