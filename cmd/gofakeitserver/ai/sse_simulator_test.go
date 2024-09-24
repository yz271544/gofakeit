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
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/sashabaranov/go-openai"
)

func Test_Rand(t *testing.T) {

	// 设置随机数种子
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// 生成10以内的随机整数
	for i := 0; i < 10; i++ {
		randomInt := rand.Intn(10) + 1 // 生成0到9的随机整数
		fmt.Println(randomInt)
	}
}

func Test_openai_req(t *testing.T) {
	req := openai.ChatCompletionRequest{
		Model: "2024-0618-sensenova",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "system",
				Content: "test111",
			},
			{
				Role:    "user",
				Content: "hello world",
			},
		},
	}

	marshal, _ := json.Marshal(req)
	fmt.Printf("curl -N -X POST -d '%s' http://localhost:8085/openai-simulator\n", string(marshal))
}
