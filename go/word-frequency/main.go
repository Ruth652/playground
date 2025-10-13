package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func is_alphanumeric(s string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9]$`).MatchString(s)

}

func word_frq(s string) map[string]int {
	freq := map[string]int{}
	words := strings.Fields(s)

	for _, word := range words {
		word = strings.ToLower(word)
		cleanWord := ""

		for _, char := range word {
			if is_alphanumeric(string(char)) {
				cleanWord += string(char)
			}
		}

		if cleanWord != "" {
			freq[cleanWord] += 1
		}
	}

	return freq
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter a text: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	fmt.Println("%v\n", word_frq(input))
}
