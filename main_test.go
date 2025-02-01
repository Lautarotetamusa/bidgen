package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestGetUser(t *testing.T) {
    t.Parallel()
    userId := "81907341"
    userOk, _ := os.ReadFile("api_mocks/user_ok.json")
    userFailed, _ := os.ReadFile("api_mocks/user_failed.json")

    cases := []struct {
        name             string
        status           int
        body             []byte
        wantErr          error
	}{
        {
            name: "ok request",
            status: http.StatusOK,
            body: userOk,
        },
        {
            name: "api error",
            status: http.StatusInternalServerError,
            wantErr: fmt.Errorf("freelancer request failed with status code: %d", http.StatusInternalServerError),
            body: userOk,
        },
        {
            name: "status is not success",
            status: http.StatusOK,
            wantErr: fmt.Errorf("freelancer request status is not success. status: failed"),
            body: userFailed,
        },
    }

    expectedUser := &User{
        Username:"sebabatsom1", 
        DisplayName:"sebabatsom1", 
        Location: Location{
            City:"Johannesburg", 
            Country: Country{
                Name: "",
            },
        },
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            t.Parallel()

            testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                w.WriteHeader(tc.status)
                w.Write(tc.body)
            }))
            freelancer := NewTestFreelancer(testServer.URL)

            user, err := freelancer.GetUser(userId)
            if tc.wantErr != nil { // Expect an error
                if err == nil {
                    t.Errorf("expected error but found success")
                } else {
                    if tc.wantErr.Error() != err.Error() {
                        t.Errorf("\nexpected error: %#v\nreceived error: %#v\n", tc.wantErr, err)
                    }

                    if user != nil {
                        t.Errorf("expected nil user in error")
                    }
                }
            }else {
                if !reflect.DeepEqual(user, expectedUser) {
                    t.Errorf("\nexpected user: %#v\nreceived user: %#v\n", expectedUser, user)
                }
            }
        })
    }
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
    seoUrl := "iphone-app-development/commerce-Mobile-App-Development-39023417"
    expected := "https://www.freelancer.com/api/projects/0.1/projects?full_description=true&limit=1&seo_urls%5B%5D=iphone-app-development%2Fcommerce-Mobile-App-Development-39023417"
    u := makeUrl(defaultBaseUrl + getProjectsPath, seoUrl)

    if expected != u.String() {
        t.Errorf("\nexpected %s\nrecieved %s", expected, u.String())
    }
}
