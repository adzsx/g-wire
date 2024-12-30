package utils

import (
	"fmt"

	"log"
	"os"
	"strings"
)

var (
	Verbose int
	Format  bool = true
)

func Err(err error, critical bool) {
	Ansi("\033[31m")
	if err != nil {
		log.Println("Error:", err)
		if critical {
			os.Exit(0)
		}
	}
	Ansi("\033[0m")
}

func InSlice(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// Verbose print
func Print(v any, level int) {

	// Ansi("\033[33m")
	// Ansi("\x1b[s\033[999B")
	// Ansi("\033[2A\033[999D\033[K\033[L")
	if Verbose >= level {
		log.SetFlags(0)

		Ansi("\033[36m")
		fmt.Print("System: ")

		log.Print(v)
		Ansi("\033[0m")
	}

	// Ansi("\033[0m\x1b[u")
}

func Ansi(inp string) {
	if Format {
		fmt.Print(inp)
	}
}

func Remove(list []string, str string) []string {
	var final []string
	for _, element := range list {
		if element != str {
			final = append(final, element)
		}
	}

	return final
}

func FilterChar(str string, char string, before bool) string {
	var final string

	for index, element := range str {
		if before {
			if string(element) != char {
				final += string(element)
			} else {
				return final
			}
		} else {

			if string(element) == char {
				final += str[index+1:]
			}
		}

	}

	return final
}

func aton(letter rune) int {
	uppercaseLetter := strings.ToUpper(string(letter))
	if len(uppercaseLetter) != 1 || uppercaseLetter < "A" || uppercaseLetter > "Z" {
		return 0 // Return 0 for non-letter characters or invalid input
	}

	num := int(uppercaseLetter[0] - 'A' + 1)
	return num
}
