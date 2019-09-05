package sonarqube

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Payload struct {
	ServerURL  string `json:"serverUrl"`
	TaskID     string `json:"taskId"`
	Status     string `json:"status"`
	AnalysedAt string `json:"analysedAt"`
	Revision   string `json:"revision"`
	ChangedAt  string `json:"changedAt"`
	Project    struct {
		Key  string `json:"key"`
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"project"`
	Branch struct {
		Name   string `json:"name"`
		Type   string `json:"type"`
		IsMain bool   `json:"isMain"`
		URL    string `json:"url"`
	} `json:"branch"`
	QualityGate struct {
		Name       string `json:"name"`
		Status     string `json:"status"`
		Conditions []struct {
			Metric         string `json:"metric"`
			Operator       string `json:"operator"`
			Value          string `json:"value,omitempty"`
			Status         string `json:"status"`
			ErrorThreshold string `json:"errorThreshold"`
		} `json:"conditions"`
	} `json:"qualityGate"`
	Properties struct {
	} `json:"properties"`
}

type IssuesSearch struct {
	Total  int `json:"total"`
	P      int `json:"p"`
	Ps     int `json:"ps"`
	Paging struct {
		PageIndex int `json:"pageIndex"`
		PageSize  int `json:"pageSize"`
		Total     int `json:"total"`
	} `json:"paging"`
	EffortTotal int     `json:"effortTotal"`
	DebtTotal   int     `json:"debtTotal"`
	Issues      []Issue `json:"issues"`
	Components  []struct {
		Organization string `json:"organization"`
		Key          string `json:"key"`
		UUID         string `json:"uuid"`
		Enabled      bool   `json:"enabled"`
		Qualifier    string `json:"qualifier"`
		Name         string `json:"name"`
		LongName     string `json:"longName"`
		Path         string `json:"path,omitempty"`
	} `json:"components"`
	Facets []interface{} `json:"facets"`
}

type Issue struct {
	Key       string `json:"key"`
	Rule      string `json:"rule"`
	Severity  string `json:"severity"`
	Component string `json:"component"`
	Project   string `json:"project"`
	Line      int    `json:"line"`
	Hash      string `json:"hash"`
	TextRange struct {
		StartLine   int `json:"startLine"`
		EndLine     int `json:"endLine"`
		StartOffset int `json:"startOffset"`
		EndOffset   int `json:"endOffset"`
	} `json:"textRange"`
	Flows              []interface{} `json:"flows"`
	Status             string        `json:"status"`
	Message            string        `json:"message"`
	Effort             string        `json:"effort"`
	Debt               string        `json:"debt"`
	Tags               []string      `json:"tags"`
	CreationDate       string        `json:"creationDate"`
	UpdateDate         string        `json:"updateDate"`
	Type               string        `json:"type"`
	Organization       string        `json:"organization"`
	FromHotspot        bool          `json:"fromHotspot"`
	ExternalRuleEngine string        `json:"externalRuleEngine,omitempty"`
}

type Sources struct {
	Scm [][]interface{} `json:"scm"`
}

func GetIssues(serverURL string, projectKey string) []Issue {
	url := fmt.Sprintf("%s/api/issues/search?componentKeys=%s&resolved=false&severities=MAJOR,CRITICAL,BLOCKER", serverURL, projectKey)
	req, _ := http.NewRequest("GET", url, nil)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var issues IssuesSearch
	err = json.Unmarshal(body, &issues)
	if err != nil {
		panic(err)
	}
	return issues.Issues
}

func GetAuthor(serverURL string, key string, from int, to int) string {
	url := fmt.Sprintf("%s/api/sources/scm?key=%s&from=%d&to=%d", serverURL, key, from, to)
	fmt.Println(url)
	req, _ := http.NewRequest("GET", url, nil)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var sources Sources
	err = json.Unmarshal(body, &sources)
	if err != nil {
		panic(err)
	}
	if len(sources.Scm) == 0 {
		return "Unknown"
	}
	return sources.Scm[0][1].(string)
}
