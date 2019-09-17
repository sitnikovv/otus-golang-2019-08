package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Word struct {
	value string
	number int
}

func main() {

	var (
		err error
		text string
		minLetter = -1
	)

	//	Что считать словом? Сколько минимально букв?
	for ; minLetter == -1; {
		fmt.Printf("Enter the minimum letter number for word and press Enter:\n")
		if text, err := readText(true); err == nil {
			if num, err := strconv.Atoi(strings.Trim(text, "\n\r")); err == nil {
				minLetter = num
				fmt.Println("Word minimum letter: " + strconv.Itoa(minLetter))
			} else {
				fmt.Println("Error: " + err.Error() + "\n\n")
			}
		} else {
			panic("Failed input text, reason: " + err.Error())
		}
	}
	fmt.Printf("\n\n")

	//	Вводим текст
	for ; text == ""; {
		fmt.Println("Enter the text, for end press Enter twice:")
		text, err = readText(false)
		text = strings.Trim(text, "\r\n")
		if err != nil {
			panic("Failed input text, reason: " + err.Error())
		} else if text == "" {
			fmt.Printf("You must enter text\n\n")
		}
	}
	fmt.Println("\n\nProcessing..." + text)
	fmt.Printf("\n\n%v\n", getWords(text, minLetter))
}

func readText(onlyOneLine bool) (string, error) {
	buffer := make([]byte, 0, 1024)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		buffer = append(buffer, scanner.Bytes()...)
		buffer = append(buffer, '\n')
		if (onlyOneLine && len(buffer) > 0 && buffer[len(buffer) - 1 ] == '\n') || (len(buffer) > 1 && string(buffer[len(buffer) - 2 : ]) == "\n\n") {
			break
		}
	}
	return strings.Trim(string(buffer), "\r\n"), nil
}

//	Parse text to words, count it and return slice of top 10 words
//	text - text to parse
//	minLetter - number of minimum letter in word
func getWords(text string, minLetter int) []Word {

	//	Split text to words and count each
	reg := regexp.MustCompile("([a-zA-Zа-яА-Я\\-]{" + strconv.Itoa(minLetter) + ",})").FindAllString(strings.ToLower(text), -1)
	cnt := make(map[string]int)
	for _, word := range reg {
		if word != "" {
			cnt[word]++
		}
	}

	//	Sort words
	result := make([]Word, 0)
	for v, n := range cnt {
		result = append(result, Word{v, n})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].number >= result[j].number
	})

	if len(result) > 10 {
		return result[:10]
	}
	return result
}
