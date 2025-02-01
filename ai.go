package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const defaultTemp float64 = 1.3

const (
    openAiapiUrl   string = "https://api.openai.com/v1/chat/completions"
    deepSeekApiUrl string = "https://api.deepseek.com/chat/completions"
    llamaApiUrl    string = "http://localhost:11434/api/chat" 
)

type Model string
const (
    GPT4o           Model = "gpt-4o"
    GPT4oMini       Model = "gpt-4o-mini"
    DeepSeekChat    Model = "deepseek-chat"
    DeepSeekReasone Model = "deepseek-reasone"
    LlamaDeepSeekR1 Model = "deepseek-r1"
)

type Role string
const (
    UserRole        Role = "user"
    SystemRole      Role = "system"
    AssistantRole   Role = "assistant"
)

type AIModel struct {
    apiKey  string
    model   Model
    client  *http.Client
    chatUrl string
}

type Message struct {
    Role    Role   `json:"role"`
    Content string `json:"content"`
}

type Payload struct {
    Model       Model      `json:"model"`
    Messages    []Message   `json:"messages"`
    Temp        float64     `json:"temperature"`
}

type Choice struct {
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
	Index        int     `json:"index"`
}

type Response struct {
	Choices []Choice `json:"choices"`
}

func NewAIModel(apiKey string, model Model) *AIModel {
    chatUrl := openAiapiUrl
    if model == DeepSeekChat || model == DeepSeekReasone {
        chatUrl = deepSeekApiUrl
    }

    if model == LlamaDeepSeekR1 {
        chatUrl = llamaApiUrl
    }

    return &AIModel {
        apiKey: apiKey,
        model: model,
        chatUrl: chatUrl,
        client: &http.Client{
            Timeout: 60 * time.Second,
        },
    }
}

func (m *AIModel) CreateBid(project *Project, temp float64) (string, error) {
    if temp == 0.0 {
        temp = defaultTemp
    }
    fmt.Println("temp", temp)

    messages := []Message{
        {
            Role: SystemRole,
            Content: prompt,
        },
        {
            Role: UserRole,
            Content: project.Title + "\n" + project.Description,
        },
    }

    res, err := m.makeRequest(&Payload{
        Model: m.model,
        Messages: messages,
        Temp: temp,
    })
    if err != nil {
        return "", err
    }

    return res.Choices[0].Message.Content, nil
}

func (m *AIModel) makeRequest(p *Payload) (*Response, error)  {
    var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(p)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, m.chatUrl, &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + m.apiKey)

	res, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var data Response
	if err = json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API request failed with status code: %d", res.StatusCode)
	}

	return &data, nil
}
