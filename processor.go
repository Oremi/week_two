package main

import (
	"fmt"
	"strconv"
	"strings"
)

//func ProcessText(text string) string {
//	return strings.Join(text, " ")
//}

func DecimalConversion(words []string) []string {

	for i := 0; i < len(words); i++ {
		// for hexadecimal conversion
		if words[i] == "(hex)" && i > 0 {
			if val, err := strconv.ParseInt(words[i-1], 16, 64); err == nil {
				words[i-1] = strconv.FormatInt(val, 10)
				// Remove the "(hex)" from the slice and adjust the index (before the index of "(hex)" and after the index of "(hex)")
			}
			words = RemoveWords(words, i)
			i--

		} else if words[i] == "(bin)" && i > 0 {
			if val, err := strconv.ParseInt(words[i-1], 2, 64); err == nil {
				words[i-1] = strconv.FormatInt(val, 10)
				// Remove the "(bin)" from the slice and adjust the index (before the index of "(bin)" and after the index of "(bin)")
			}
			words = RemoveWords(words, i)
			i--

		}
	}
	return words

}

func CaseConversion(words []string) []string {

	for i := 0; i < len(words); i++ {
		match := words[i]
		// strings.Tolower, Toupper
		if match == "(up)" && i > 0 {
			words[i-1] = strings.ToUpper(words[i-1])
			words = RemoveWords(words, i)
			i--
		}

		// change to lowercase
		if match == "(low)" && i > 0 {
			words[i-1] = strings.ToLower(words[i-1])
			words = RemoveWords(words, i)
			i--
		}

		// // for capitalize
		if match == "(cap)" && i > 0 {
			words[i-1] = strings.ToUpper(words[i-1][0:1]) + strings.ToLower(words[i-1][1:])
			words = RemoveWords(words, i)
			i--
		}

		// for (up, N), (low, N) and (cap, N)
		if match == "(up," || match == "(low," || match == "(cap," {

			//ensure that the words exist
			if i+1 >= len(words) {
				continue
			}
			ndx := strings.TrimSuffix(words[i+1], ")")
			ndxInt, err := strconv.Atoi(ndx)
			if err != nil {
				fmt.Printf("Error converting %s to integer: %v\n", ndx, err)
				continue
			}
			if i-ndxInt < 0 {
				fmt.Printf("Error: Not enough words to convert for %s\n", words[i])
				continue
			}

			for j := i - ndxInt; j < i; j++ {
				if match == "(up," {
					words[j] = strings.ToUpper(words[j])
				} else if match == "(low," {
					words[j] = strings.ToLower(words[j])
				} else if match == "(cap," {
					words[j] = strings.ToUpper(words[j][0:1]) + strings.ToLower(words[j][1:])
				}
			}
			words = RemoveWords(words, i)
			words = RemoveWords(words, i)
			i--
		}

	}
	return words
}

func RemoveWords(words []string, index int) []string {
	result := append(words[:index], words[index+1:]...)
	return result

}

func main() {
	text := "this is a, sample , string for (up, 3) and 1a (hex) and 1010 (bin)"
	words := strings.Fields(text)
	words = DecimalConversion(words)
	words = CaseConversion(words)
	fmt.Println(strings.Join(words, " "))
}
