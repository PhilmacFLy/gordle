package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	var file string

	// Get the file name from the command line
	flag.StringVar(&file, "f", "config/openthesaurus.txt", "Specify the path of the openthesaurus file")

	//Read Open Thesaurus File and scan it
	readFile, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	var words []string
	existingsWords := make(map[string]bool)
	//For every line in the file
	for fileScanner.Scan() {
		//If the line contains # then we skip it as it is a comment
		if strings.Contains(fileScanner.Text(), "#") {
			continue
		}
		//Split the line by semicolon
		split := strings.Split(fileScanner.Text(), ";")
		for _, s := range split {
			//Trim the spaces
			s = strings.TrimSpace(s)
			//If the line contains brackets remove the content inside the brackets
			//As those are only explanation of the word
			if strings.Contains(s, "(") {
				i1 := strings.Index(s, "(")
				i2 := strings.Index(s, ")")
				s = s[:i1] + s[i2+1:]
				s = strings.TrimSpace(s)
			}
			//If we still have spaces then its a saying so we skip it
			if strings.Contains(s, " ") {
				continue
			}
			//Remove all ...
			s = strings.TrimPrefix(s, "...")
			s = strings.TrimSuffix(s, "...")
			//Replace german characters with wordle characters
			s = strings.ReplaceAll(s, "ä", "ae")
			s = strings.ReplaceAll(s, "ö", "oe")
			s = strings.ReplaceAll(s, "ü", "ue")
			s = strings.ReplaceAll(s, "ß", "ss")
			//Remove ! and ?
			s = strings.ReplaceAll(s, "!", "")
			s = strings.ReplaceAll(s, "?", "")
			//Because wordle only has 5 characters we dont nee the rest
			if len(s) != 5 {
				continue
			}
			//Append the word to the slice if it is not already in the slice
			if !existingsWords[s] {
				words = append(words, strings.ToLower(s))
				existingsWords[s] = true
			}
		}
	}
	//Close the file
	readFile.Close()

	//Ask user what we know about the word
	fmt.Println("Enter all needed Characters: ")
	var needed string
	fmt.Scanln(&needed)

	fmt.Println("Enter all forbidden Characters: ")
	var forbidden string
	fmt.Scanln(&forbidden)

	fmt.Println("Enter word format: ")
	var format string
	fmt.Scanln(&format)

	needed = strings.TrimSpace(strings.ToLower(needed))
	forbidden = strings.TrimSpace(strings.ToLower(forbidden))
	format = strings.TrimSpace(strings.ToLower(format))

	//Split everything into single characters
	neededslice := strings.Split(needed, "")
	forbiddenslice := strings.Split(forbidden, "")
	formatslice := strings.Split(format, "")

	var possibleWords []string
	for _, word := range words {
		add := true
		//Check if the word contains all needed characters
		for _, needed := range neededslice {
			if !strings.Contains(word, needed) {
				add = false
				break
			}
		}
		//Check if the word contains any forbidden characters
		for _, forbidden := range forbiddenslice {
			if strings.Contains(word, forbidden) {
				add = false
				break
			}
		}
		//Check if the word matches the format
		for i, c := range formatslice {
			if c != "_" {
				if c != string(word[i]) {
					add = false
					break
				}
			}
		}
		//If all checks are passed we add the word to the slice
		if add {
			possibleWords = append(possibleWords, word)
		}
	}
	//Print all possible words
	for _, word := range possibleWords {
		fmt.Println(word)
	}
}
