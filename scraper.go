package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

// API endpoints urls
const (
    baseUrl = "https://www.freelancer.com/api/"

    getProjects = baseUrl + "projects/0.1/projects"
    getUser     = baseUrl + "users/0.1/users/%s" // %s: user_id
)

type FreelanceResponse struct {
    Status  string  `json:"status"`
    Result  json.RawMessage  `json:"result"`
}

type ProjectResult struct {
    Projects []Project `json:"projects"`
}

type User struct {
    Username    string      `json:"username"`
    DisplayName string      `json:"display_name"`
    Location    Location    `json:"location"`
}


type Project struct {
    Title       string `json:"title"`
    Description string `json:"description"`
    Budget      Budget `json:"budget"`
    Status      string `json:"status"`
}

type Location struct {
    City    string   `json:"city"`
    Country Country  `json:"country"`
}

type Country struct {
    Name    string  `json:"country"`
}

type Budget struct {
    Min float32 `json:"minimum"`
    Max float32 `json:"maximum"`
}

var client *http.Client = &http.Client{
    Timeout: 15 * time.Second,
}

func GetUser(userId string) (*User, error) {
    data, err := makeRequest(fmt.Sprintf(getUser, userId))
    if err != nil {
        return nil, err
    }

    var u User
    if err := json.Unmarshal(data.Result, &u); err != nil {
        return nil, err
    }

    return &u, nil
}

func GetProyect(projectUrl string) (*Project, error) {
    seoUrl := extractSeo(projectUrl)
    u := makeUrl(seoUrl)

    data, err := makeRequest(u.String())
    if err != nil {
        return nil, err
    }

    var p ProjectResult
    if err := json.Unmarshal(data.Result, &p); err != nil {
        return nil, err
    }

    return &p.Projects[0], nil
}

func makeRequest(url string) (*FreelanceResponse, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
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
        return nil, fmt.Errorf("freelancer request failed with status code: %d", res.StatusCode)
    }

    if data.Status != "success" {
        return nil, fmt.Errorf("freelancer request status is not success. status: %s", data.Status)
    }

    return &data, nil
}

func extractSeo(u string) string {
// https://www.freelancer.com.ar/projects/iphone-app-development/commerce-Mobile-App-Development-39023417/proposals
    regex := regexp.MustCompile(`projects/([^/]+)/([^/]+)/`)
	match := regex.FindStringSubmatch(u)

	if len(match) > 0 {
		extracted := match[1] + "/" + match[2]
        return extracted
	}	

    fmt.Println("No match found")
    return ""
}

func makeUrl(seoUrl string) url.URL {
    u, err := url.Parse(getProjects)
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
