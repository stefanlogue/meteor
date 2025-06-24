package main

func wordWrap(input string, width int) string {
	words := splitIntoWords(input)
	if len(words) == 0 {
		return ""
	}
	var result string
	currentLine := ""
	for _, word := range words {
		if len(currentLine)+len(word)+1 > width {
			if currentLine != "" {
				result += currentLine + "\n"
			}
			currentLine = word
		} else {
			if currentLine != "" {
				currentLine += " "
			}
			currentLine += word
		}
	}
	if currentLine != "" {
		result += currentLine
	}
	return result
}

func splitIntoWords(input string) []string {
	words := []string{}
	currentWord := ""
	for _, char := range input {
		if char == ' ' || char == '\n' || char == '\t' {
			if currentWord != "" {
				words = append(words, currentWord)
				currentWord = ""
			}
		} else {
			currentWord += string(char)
		}
	}
	if currentWord != "" {
		words = append(words, currentWord)
	}
	return words
}
