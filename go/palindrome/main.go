package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func is_alphanumeric(s string) bool {

	reg := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return reg.MatchString(s)
}

func is_palindrome(s string) bool {
	l, r := 0, len(s)-1

	for l < r {
		left, right := string(s[l]), string(s[r])
		if !is_alphanumeric(left) {
			l += 1
			continue
		}
		if !is_alphanumeric(right) {
			r -= 1
			continue
		}

		if !strings.EqualFold(left, right) {
			return false
		}

		l++
		r--
	}

	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter a Word: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	fmt.Println()
	fmt.Printf("Is %s a palindrome? %v\n", input, is_palindrome(input))

}
