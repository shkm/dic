package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"io"
	"net/http"
	"os"
	"strings"
)

const BaseUrl = "https://api.dictionaryapi.dev/api/v2/entries/en/"

type Phonetic struct {
	Text  string `json:"text"`
	Audio string `json:"audio"`
}

type Definition struct {
	Definition string   `json:"definition"`
	Example    string   `json:"example"`
	Synonyms   []string `json:"synonyms"`
	Antonyms   []string `json:"antonyms"`
}

type Meaning struct {
	PartOfSpeech string       `json:"partOfSpeech"`
	Definitions  []Definition `json:"definitions"`
}

type Word struct {
	Word      string     `json:"word"`
	Phonetic  string     `json:"phonetic"`
	Phonetics []Phonetic `json:"phonetics"`
	Origin    string     `json:"origin"`
	Meanings  []Meaning  `json:"meanings"`
}

func Lookup(phrase string) ([]Word, error) {
	url := BaseUrl + phrase

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		if resp.StatusCode == 404 {
			return nil, errors.New("Couldn't find word.")
		}
		return nil, errors.New("Something went wrong.")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var words []Word
	err = json.Unmarshal(body, &words)
	if err != nil {
		return nil, err
	}
	return words, nil
}

func usage() {
	fmt.Print(`dic

Usage: dic <query>`)
}

func printWord(word Word, num int) {
	c := color.New(color.Bold).Add(color.FgGreen)
	c.Printf("%d. %s", num+1, word.Word)

	printPhonetics(word.Phonetics)

	if len(word.Meanings) > 0 {
		fmt.Println("")
		for _, meaning := range word.Meanings {
			printMeaning(meaning)
		}
	}
}

func printMeaning(meaning Meaning) {
	c := color.New(color.FgYellow)
	if len(meaning.PartOfSpeech) > 0 {
		c.Printf("\n%s\n", meaning.PartOfSpeech)
	}

	contentAfterDefinition := false
	for _, definition := range meaning.Definitions {
		if contentAfterDefinition {
			fmt.Println("")
		}

		c := color.New(color.Italic)
		c.Print("def. ")
		c = color.New(color.Bold).Add(color.FgBlue)
		c.Printf("%s\n", definition.Definition)

		if len(definition.Example) > 0 {
			contentAfterDefinition = true
			c := color.New(color.Italic)
			c.Print("ex.  ")
			c.DisableColor()

			c.Printf("%s\n", definition.Example)
		}

		c = color.New(color.Italic)

		var presentSynonyms []string
		for _, synonym := range definition.Synonyms {
			if len(synonym) > 0 {
				presentSynonyms = append(presentSynonyms, synonym)
			}
		}

		if len(presentSynonyms) > 0 {
			contentAfterDefinition = true
			synonyms := strings.Join(presentSynonyms, ", ")

			c.EnableColor()
			c.Print("syn.")
			c.DisableColor()
			c.Printf(" %s\n", synonyms)
		}

		var presentAntonyms []string
		for _, antonym := range definition.Antonyms {
			if len(antonym) > 0 {
				presentAntonyms = append(presentAntonyms, antonym)
			}
		}

		if len(presentAntonyms) > 0 {
			contentAfterDefinition = true
			antonyms := strings.Join(presentAntonyms, ", ")

			c.EnableColor()
			c.Print("ant.")
			c.DisableColor()
			c.Printf(" %s\n", antonyms)
		}

	}
}

func printPhonetics(phonetics []Phonetic) {
	var phoneticTexts []string
	for _, phonetic := range phonetics {
		if len(phonetic.Text) > 0 {
			phoneticTexts = append(phoneticTexts, phonetic.Text)
		}
	}

	if len(phoneticTexts) > 0 {
		phoneticsString := strings.Join(phoneticTexts, ", ")
		fmt.Printf(" [%s]", phoneticsString)
	}
}

func main() {
	if len(os.Args) == 1 {
		usage()
		os.Exit(0)
	}

	words, err := Lookup(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for i, word := range words {
		if i > 0 {
			fmt.Println("")
		}
		printWord(word, i)
	}
}
