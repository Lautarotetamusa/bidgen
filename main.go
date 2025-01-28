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
    // temp := flag.Float64("temp", defaultTemp, "temperature of ai request")

    flag.Parse()
    
    fmt.Println(*projectUrl)
    if len(*projectUrl) < 2 {
        panic("project url must not be empty")
    }

    project, err := GetProyect(*projectUrl)
    if err != nil {
        panic(err)
    }

    println("\n\n--- PROYECT DESCRIPTION ---")
    println(project.Description)

    ai := NewAIModel(apiKey, DeepSeekChat)

    println("\n\n--- BID ---")
    res, err := ai.CreateBid(project.Description)
    if err != nil {
        panic(err)
    }
    println(res)
}
