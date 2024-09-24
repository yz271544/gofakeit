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

import "testing"

// 	contents := []string{"i", "love", "you"}
//
//	contents := []string{"我", "爱", "你"}

func Test_joinWords(t *testing.T) {
	type args struct {
		words []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test joinWords with English words",
			args: args{words: []string{"i", "love", "you"}},
			want: "i love you",
		},
		{
			name: "Test joinWords with Chinese words",
			args: args{words: []string{"我", "爱", "你"}},
			want: "我爱你",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := joinWords(tt.args.words); got != tt.want {
				t.Errorf("joinWords() = %v, want %v", got, tt.want)
			}
		})
	}
}
