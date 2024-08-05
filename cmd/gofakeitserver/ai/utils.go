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
	"strings"
	"unicode"
)

// isChinese checks if a rune is a Chinese character
func isChinese(r rune) bool {
	// Check if the rune is within the range of Chinese characters
	return unicode.Is(unicode.Han, r)
}

// joinWords joins a slice of strings into a sentence based on the language
func joinWords(words []string) string {
	// Convert the first word to rune slice to check the language
	if len(words) == 0 {
		return ""
	}

	// Check the first rune to determine language type
	firstRune := []rune(words[0])[0]
	if isChinese(firstRune) {
		// Join without spaces for Chinese
		return strings.Join(words, "")
	}
	// Default to joining with spaces for other languages
	return strings.Join(words, " ")
}
