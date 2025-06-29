package utils

import (
	"encoding/json"
	"fmt"
	"unicode"
)

func PrintJSON(v any) {
	jsonBytes, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(jsonBytes))
}

func CapitalizeFirstLetter(s string) string {
	if s == "" {
		return ""
	}

	// Convert string to a slice of runes to handle Unicode characters correctly
	runes := []rune(s)

	// If the first character is a letter, convert it to uppercase
	if unicode.IsLetter(runes[0]) {
		runes[0] = unicode.ToUpper(runes[0])
	}

	return string(runes)
}
