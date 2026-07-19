package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

const commitTypes = "build|chore|ci|docs|feat|fix|perf|refactor|revert|test"

var (
	typePattern = regexp.MustCompile(`^(` + commitTypes + `)`)
	fullPattern = regexp.MustCompile(
		`^(` + commitTypes + `)(\([a-z0-9._/-]+\))?(!)?: [^\s].+$`,
	)
)

func main() {
	if len(os.Args) < 2 {
		fail("commit message file not provided")
	}

	content, err := os.ReadFile(os.Args[1])
	if err != nil {
		fail("could not read commit message file")
	}

	message := normalize(string(content))

	if shouldSkip(message) {
		return
	}

	validate(message)
}

func validate(message string) {
	if fullPattern.MatchString(message) {
		return
	}

	if !typePattern.MatchString(message) {
		fail("missing or invalid commit type")
	}

	validateScope(message)
	validateSeparator(message)

	fail("invalid commit message format")
}

func validateScope(message string) {
	open := strings.Count(message, "(")
	close := strings.Count(message, ")")

	if open > 0 && close == 0 {
		fail("scope is malformed (missing closing parenthesis)")
	}

	if close > 0 && open == 0 {
		fail("scope is malformed (missing opening parenthesis)")
	}

	if open != close {
		fail("scope is malformed")
	}

	if strings.Contains(message, "()") {
		fail("scope cannot be empty")
	}
}

func validateSeparator(message string) {
	index := strings.Index(message, ":")

	if index == -1 {
		fail("missing ':' after type or scope")
	}

	after := message[index+1:]

	if after == "" || after[0] != ' ' {
		fail("commit description is required")
	}

	if strings.TrimSpace(after[1:]) == "" {
		fail("use ': ' (colon + space) before description")
	}
}

func normalize(value string) string {
	lines := strings.Split(strings.ReplaceAll(value, "\r\n", "\n"), "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			return trimmed
		}
	}

	return ""
}

func shouldSkip(message string) bool {
	return strings.HasPrefix(message, "Merge ") ||
		strings.HasPrefix(message, "Revert \"")
}

func fail(msg string) {
	fmt.Fprintf(os.Stderr, "error: %s\n", msg)
	os.Exit(1)
}
