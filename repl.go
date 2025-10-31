package main

import ( 
	"strings"
)

func cleanInput(text string) []string {
	// Convert to lowercase
	lower := strings.ToLower(text)

	// Trim leading/trailing whitespace
	trimmed := strings.TrimSpace(lower)

	// Split by whitespace
	words := strings.Fields(trimmed)

	return words
}