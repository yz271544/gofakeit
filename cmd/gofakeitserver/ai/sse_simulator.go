/*
Copyright 2024 The gofakeit Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/sashabaranov/go-openai"
)

var clients = make(map[*SSEClient]bool)
var addClient = make(chan *SSEClient)
var removeClient = make(chan *SSEClient)
var broadcast = make(chan string)

// 定义一个结构体来表示SSE客户端连接
type SSEClient struct {
	id      string
	message chan string
}

func SseHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	client := &SSEClient{
		id:      r.RemoteAddr,
		message: make(chan string),
	}

	addClient <- client

	defer func() {
		removeClient <- client
	}()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {
		select {
		case msg := <-client.message:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		}
	}
}

func manageClients() {
	for {
		select {
		case client := <-addClient:
			clients[client] = true
		case client := <-removeClient:
			delete(clients, client)
			close(client.message)
		case message := <-broadcast:
			for client := range clients {
				client.message <- message
			}
		}
	}
}

func OpenAIHandler(w http.ResponseWriter, r *http.Request) {
	prompt := r.URL.Query().Get("prompt")
	if prompt == "" {
		http.Error(w, "Prompt is required", http.StatusBadRequest)
		return
	}

	response, err := callOpenAI(prompt)
	if err != nil {
		http.Error(w, "Error calling OpenAI API", http.StatusInternalServerError)
		return
	}

	broadcast <- response
}

func callOpenAI(prompt string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OpenAI API key is not set")
	}

	requestBody, err := json.Marshal(map[string]string{
		"model":  "text-davinci-003",
		"prompt": prompt,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("Unexpected response format from OpenAI API")
	}

	text, ok := choices[0].(map[string]interface{})["text"].(string)
	if !ok {
		return "", fmt.Errorf("Unexpected response format from OpenAI API")
	}

	return text, nil
}

func OpenAIChatCompletionsSimulatorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Error openai api only support POST method", http.StatusMethodNotAllowed)
		return
	}

	rand.New(rand.NewSource(time.Now().UnixNano()))

	var req = &openai.ChatCompletionRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	randomInt := rand.Intn(50) + 1
	for i := range randomInt {
		chatResp, err := callOpenAISimulator(*req, i)
		if err != nil {
			http.Error(w, "Error calling OpenAI API", http.StatusInternalServerError)
			return
		}

		data, err := json.Marshal(chatResp)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		_, err = w.Write([]byte(fmt.Sprintf("data: %s\n\n", data)))
		if err != nil {
			http.Error(w, fmt.Sprintf("Error writing response error: %v", err), http.StatusInternalServerError)
			return
		}
	}

	_, err = w.Write([]byte(fmt.Sprintf("data: %s\n\n", "[DONE]")))
	if err != nil {
		http.Error(w, fmt.Sprintf("output [DONE] Error writing response error: %v", err), http.StatusInternalServerError)
	}

	if fluster, ok := w.(http.Flusher); ok {
		fluster.Flush()
	} else {
		http.Error(w, "Flushing not supported", http.StatusInternalServerError)
		return
	}

}

func callOpenAISimulator(req openai.ChatCompletionRequest, i int) (openai.ChatCompletionStreamResponse, error) {

	car := gofakeit.Car()
	content := fmt.Sprintf("[%d] %s-%s-%s-%s-%s", car.Year, car.Model, car.Type, car.Brand, car.Fuel, car.Transmission)

	resp := openai.ChatCompletionStreamResponse{
		ID:      "example-id",
		Object:  "example-object",
		Created: time.Now().Unix(),
		Model:   req.Model,
		Choices: []openai.ChatCompletionStreamChoice{
			{
				Index: i,
				Delta: openai.ChatCompletionStreamChoiceDelta{
					Content: fmt.Sprintf("This is message %d: content: %s", i, content),
				},
				FinishReason: "length",
			},
		},
		SystemFingerprint: "example-fingerprint",
		PromptAnnotations: []openai.PromptAnnotation{
			{
				PromptIndex: 0,
				ContentFilterResults: openai.ContentFilterResults{
					Hate: openai.Hate{
						Filtered: gofakeit.Bool(),
						Severity: "",
					},
					SelfHarm: openai.SelfHarm{
						Filtered: gofakeit.Bool(),
						Severity: "",
					},
					Sexual: openai.Sexual{
						Filtered: gofakeit.Bool(),
						Severity: "",
					},
					Violence: openai.Violence{
						Filtered: gofakeit.Bool(),
						Severity: "",
					},
				},
			},
		},
		PromptFilterResults: []openai.PromptFilterResult{
			{
				Index: 0,
				ContentFilterResults: openai.ContentFilterResults{
					Hate: openai.Hate{
						Filtered: gofakeit.Bool(),
						Severity: "",
					},
					SelfHarm: openai.SelfHarm{
						Filtered: gofakeit.Bool(),
						Severity: "",
					},
					Sexual: openai.Sexual{
						Filtered: gofakeit.Bool(),
						Severity: "",
					},
					Violence: openai.Violence{
						Filtered: gofakeit.Bool(),
						Severity: "",
					},
				},
			},
		},
		Usage: &openai.Usage{
			PromptTokens:     10,
			CompletionTokens: 20,
			TotalTokens:      30,
		},
	}

	return resp, nil
}
