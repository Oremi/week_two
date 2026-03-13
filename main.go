package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run . input.txt output.txt")
		return
	}

	input := os.Args[1]
	output := os.Args[2]

	data, err := os.ReadFile(input)
	if err != nil {
		fmt.Println("Error reading file")
		return
	}

	text := string(data)

	text = processText(text)
	text = fixPunctuation(text)
	text = fixArticles(text)
	text = fixQuotes(text)

	os.WriteFile(output, []byte(text), 0644)
}

func processLine(line string) string {
	words := strings.Fields(line)
	result := []string{}

	for i := 0; i < len(words); i++ {
		w := words[i]

		if w == "(hex)" && len(result) > 0 {
			n, err := strconv.ParseInt(result[len(result)-1], 16, 64)
			if err == nil {
				result[len(result)-1] = strconv.FormatInt(n, 10)
			}
			continue
		}

		if w == "(bin)" && len(result) > 0 {
			n, err := strconv.ParseInt(result[len(result)-1], 2, 64)
			if err == nil {
				result[len(result)-1] = strconv.FormatInt(n, 10)
			}
			continue
		}

		if w == "(up)" && len(result) > 0 {
			result[len(result)-1] = strings.ToUpper(result[len(result)-1])
			continue
		}

		if w == "(low)" && len(result) > 0 {
			result[len(result)-1] = strings.ToLower(result[len(result)-1])
			continue
		}

		if w == "(cap)" && len(result) > 0 {
			word := result[len(result)-1]
			result[len(result)-1] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
			continue
		}

		if w == "(up," && i+1 < len(words) {
			n := getNumber(words[i+1])
			for j := len(result) - 1; j >= 0 && n > 0; j-- {
				result[j] = strings.ToUpper(result[j])
				n--
			}

			i++
			continue
		}

		if strings.HasPrefix(w, "(low,") {
			n := getNumber(words[i+1])
			for j := len(result) - 1; j >= 0 && n > 0; j-- {
				result[j] = strings.ToLower(result[j])
				n--
			}
			i++
			continue
		}

		if strings.HasPrefix(w, "(cap,") {
			n := getNumber(words[i+1])
			for j := len(result) - 1; j >= 0 && n > 0; j-- {
				word := result[j]
				result[j] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
				n--
			}
			i++
			continue
		}

		result = append(result, w)
	}

	return strings.Join(result, " ")
}

func processText(text string) string {
	linesArray := strings.Split(text, "\n")
	var result []string

	//  init. ;  condition         ; post
	for i := 0; i < len(linesArray); i++ {
		newLine := processLine(linesArray[i])
		result = append(result, newLine)
	}

	// fmt.Println(strings.Join(result, "\n"))
	return strings.Join(result, "\n")
}

func getNumber(s string) int {
	n := 0
	if strings.HasSuffix(s, ")") {
		s = strings.TrimSuffix(s, ")")
		number, _ := strconv.Atoi(s)
		n = number
	} else if strings.HasSuffix(s, "),") {
		s = strings.TrimSuffix(s, "),")
		number, _ := strconv.Atoi(s)
		n = number
	}
	return n
}

func fixPunctuation(text string) string {
	p := []string{".", ",", "!", "?", ":", ";"}

	for _, mark := range p {
		// Remove space before punctuation
		for strings.Contains(text, " "+mark) {
			text = strings.ReplaceAll(text, " "+mark, mark)
		}
			
		// Ensure space after the punctuation if followed by a letter
			for i := 0; i < len(text)-1; i++ {
			if text[i] == mark[0] && text[i-1] != ' ' && ((text[i+1] >= 'A' && text[i+1] <= 'Z') || (text[i+1] >= 'a' && text[i+1] <= 'z')) {
				text = text[:i+1] + " " + text[i+1:]
			}
		}
	}

	// fix multiple punctuation groups
	text = strings.ReplaceAll(text, " ... ", "... ")
	text = strings.ReplaceAll(text, " ...", "... ")
	text = strings.ReplaceAll(text, " !?" , "!?")

	return text
}

func fixQuotes(text string) string {
        text = strings.ReplaceAll(text, "' ", " '")
        text = strings.ReplaceAll(text, " '", "'")
        return text
}

func fixArticles(text string) string {
	textArray := strings.Split(text, "\n")

	for eachIndex, eachLine := range textArray {
		arrayOfWords := strings.Fields(eachLine)
		for i := 0; i < len(arrayOfWords)-1; i++ {
			if strings.ToLower(arrayOfWords[i]) == "a" {
				first := strings.ToLower(string(arrayOfWords[i+1][0]))
				if strings.Contains("aeiouh", first) {
					arrayOfWords[i] = "an"
				}
			}
		}
		textArray[eachIndex] = strings.Join(arrayOfWords, " ")
	}

	return strings.Join(textArray, "\n")
}
