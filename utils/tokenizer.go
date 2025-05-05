package utils

import (
	"strings"
)

func TokenizeStep(step string) []string {
	var tokens []string
	var buff strings.Builder
	var inParenthesis bool = false

	for _, ch := range step {
		if ch == ' ' && !inParenthesis {
			newToken := buff.String()
			tokens = append(tokens, newToken)
			buff.Reset()
		} else if ch == '"' {
			inParenthesis = !inParenthesis
		} else {
			buff.WriteString(string(ch))
		}
	}

	if buff.Len() > 0 {
		tokens = append(tokens, buff.String())
	}

	return tokens
}
