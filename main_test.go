package main

import (
	"fmt"
	"testing"
)

func TestGetProyect(t *testing.T) {
    seo := "iphone-app-development/commerce-Mobile-App-Development-39023417"

    p, err := GetProyect(seo)
    if err != nil {
        t.Error(err)
    }

    fmt.Printf("%#v\n", p)
}

func TestExtractSeo(t *testing.T) {
    url := "https://www.freelancer.com.ar/projects/iphone-app-development/commerce-Mobile-App-Development-39023417/details"
    expected := "iphone-app-development/commerce-Mobile-App-Development-39023417" 

    seo := extractSeo(url)
    if seo != expected {
        t.Errorf("\nexpected %s\nrecieved %s", expected, seo)
    }
}

func TestMakeUrl(t *testing.T) {
    projectUrl := "https://www.freelancer.com.ar/projects/iphone-app-development/commerce-Mobile-App-Development-39023417/details"
    expected := "https://www.freelancer.com/api/projects/0.1/projects?full_description=true&limit=1&seo_urls%5B%5D=iphone-app-development%2Fcommerce-Mobile-App-Development-39023417"
    u := makeUrl(projectUrl)

    if expected != u.String() {
        t.Errorf("\nexpected %s\nrecieved %s", expected, u.String())
    }
}

func TestGetProject(t *testing.T) {
    projectUrl := "https://www.freelancer.com.ar/projects/iphone-app-development/commerce-Mobile-App-Development-39023417/details"

    GetProyect(projectUrl)
}
