package main

import (
	"flag"
	"fmt"
	"os"
)

const envFieldName  = "AI_APIKEY"

func main() {
    apiKey, exists := os.LookupEnv(envFieldName)
    if !exists {
        panic(fmt.Sprintf("%s must be present in the enviroment. try with:\nexport %s=api_key", envFieldName, envFieldName))
    }

    projectUrl := flag.String("u", "", "project url")
    // model := flag.String("model", "deepSeek", "model can be deepSeek or gpt")
    temp := flag.Float64("temp", defaultTemp, "temperature of ai request")
    flag.Parse()

    if len(*projectUrl) < 2 {
        panic("project url must not be empty")
    }

    if *temp < 0 {
        panic("temp parameter cannot be less than 0")
    }

    freelancer := NewFreelancer()

    project, err := freelancer.GetProyect(*projectUrl)
    if err != nil {
        panic(err)
    }
    fmt.Printf("%#v\n", project)

    user, err := freelancer.GetUser(project.UserId)
    if err != nil {
        println(err.Error())
        os.Exit(0)
    }
    fmt.Printf("%#v\n", user)

    fmt.Printf("%s\n%s\ncliente: %s", project.Title, project.Description, user.Username)

    // println("\n\n--- PROYECT DESCRIPTION ---")
    // println(project.Description)
    //
    ai := NewAIModel(apiKey, DeepSeekChat)

    println("\n\n--- BID ---")
    res, err := ai.CreateBid(project, *temp)
    if err != nil {
        panic(err)
    }
    println(res)
}
