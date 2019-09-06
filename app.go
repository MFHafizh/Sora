package main

import (
	"Sora/jira"
	"Sora/sonarqube"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"

	"github.com/julienschmidt/httprouter"
)

var template = "%s on %s. %s"

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Fail to load env")
	}
	port := ":8080"
	router := httprouter.New()
	router.POST("/payload", handlers)
	log.Print("running in port: " + port)
	log.Fatal(http.ListenAndServe(port, router))

}

func handlers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	projectKey := r.Header.Get("X-SonarQube-Project")
	body, _ := ioutil.ReadAll(r.Body)
	var payload sonarqube.Payload
	err := json.Unmarshal(body, &payload)
	if err != nil {
		panic(err)
	}
	serverURL := payload.ServerURL
	fmt.Println(os.Getenv("JIRA_PROJECT_KEY"))
	jiraProject := os.Getenv("JIRA_PROJECT_KEY")
	issues := sonarqube.GetIssues(serverURL, projectKey)
	jiraIssues := jira.GetIssues(jiraProject)
	for _, sonarIssue := range issues {
		isExist := false
		component := strings.Split(sonarIssue.Component, ":")
		file := component[len(component)-1]
		message := fmt.Sprintf("[SQ] %s on %s line:%d", sonarIssue.Message, file, sonarIssue.Line)
		for _, jiraIssue := range jiraIssues {
			if jiraIssue.Fields.Summary == message {
				isExist = true
			}
		}
		if !isExist {
			author := sonarqube.GetAuthor(serverURL, sonarIssue.Component, sonarIssue.Line, sonarIssue.Line)
			description := fmt.Sprintf("%s %s on %s line:%d created by: %s", sonarIssue.Severity, sonarIssue.Type, file, sonarIssue.Line, author)
			jira.CreateIssue(jiraProject, message, description)
		}
	}
}
