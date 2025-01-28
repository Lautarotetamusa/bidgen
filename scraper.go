package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

const proyectsUrl = "https://www.freelancer.com/api/projects/0.1/projects"

type FreelanceResponse struct {
    Status  string  `json:"status"`
    Result  Result  `json:"result"`
}

type Result struct {
    Projects []Project `json:"projects"`
}

type Project struct {
    Title       string `json:"title"`
    Description string `json:"description"`
    Budget      Budget `json:"budget"`
    Status      string `json:"status"`
}

type Budget struct {
    Min float32 `json:"minimum"`
    Max float32 `json:"maximum"`
}

var client *http.Client = &http.Client{
    Timeout: 15 * time.Second,
}

func GetProyect(projectUrl string) (*Project, error) {
    seoUrl := extractSeo(projectUrl)
    u := makeUrl(seoUrl)
    fmt.Println(u)

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var data FreelanceResponse
	if err = json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}

    if res.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("freelancer project request is not ok. status code: %d", res.StatusCode)
    }

    if data.Status != "success" {
        return nil, fmt.Errorf("freelancer project request status is not success. status: %s", data.Status)
    }

    return &data.Result.Projects[0], nil
}

func extractSeo(u string) string {
// https://www.freelancer.com.ar/projects/iphone-app-development/commerce-Mobile-App-Development-39023417/proposals
    regex := regexp.MustCompile(`projects/([^/]+)/([^/]+)/`)
	match := regex.FindStringSubmatch(u)

	if len(match) > 0 {
		extracted := match[1] + "/" + match[2]
		fmt.Println("Extracted part:", extracted)
        return extracted
	}	

    fmt.Println("No match found")
    return ""
}

func makeUrl(seoUrl string) url.URL {
    u, err := url.Parse(proyectsUrl)
	if err != nil {
        panic(err)
	}

    q := u.Query()
    q.Add("seo_urls[]", seoUrl)
    q.Add("limit", "1")
    q.Add("full_description", "true")

	u.RawQuery = q.Encode()    
    return *u
}
