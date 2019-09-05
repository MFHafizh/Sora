package jira

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/andygrunwald/go-jira"
)

var jiraClient = initJiraClient()

func initJiraClient() *jira.Client {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	tp := jira.BasicAuthTransport{
		Username: os.Getenv("JIRA_USER"),
		Password: os.Getenv("JIRA_TOKEN"),
	}

	jiraClient, err := jira.NewClient(tp.Client(), "https://mfhafizh.atlassian.net/")
	if err != nil {
		panic(err)
	}
	return jiraClient
}

func CreateIssue(projectKey string, title string, description string) {
	field := jira.IssueFields{
		Summary:     title,
		Description: description,
		Type: jira.IssueType{
			Name: "Bug",
		},
		Project: jira.Project{
			Key: projectKey,
		},
	}
	issue := jira.Issue{
		Fields: &field,
	}
	_, _, err := jiraClient.Issue.Create(&issue)
	if err != nil {
		panic(err)
	}
	fmt.Println("Create " + title + " issue")
}

func GetIssues(projectKey string) []jira.Issue {
	issues, _, err := jiraClient.Issue.Search("project="+projectKey, &jira.SearchOptions{})
	if err != nil {
		panic(err)
	}
	return issues
}
