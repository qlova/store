package db

import (
	"math/big"
	"strings"
	"sync"
	"text/scanner"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

//Dictionary is required for efficient natural-language search.
//How you setup/store the dictionary is up to you.
//But DictionaryMap & DictionaryViewer implementations are provided by this package.
type Dictionary interface {

	//IndexWords returns related concept-indicies for the given words.
	//The indicies should be constant (for a given word within a given dictionary)
	//and unique to the concept that the word represents.
	//prime multiples can represent sub-concepts, ie animal = 3, dog = 6.
	//prime combinations can represent intersecting concepts.
	//meaningless words should be set to 1.
	//words that are used at a higher frequency should return lower numbers.
	//Words without dictionary entries should return 0.
	LookupWords(results []int64, words ...string) error
}

func Tokenise(text string) (words []string) {
	var s scanner.Scanner
	s.Init(strings.NewReader(text))

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		words = append(words, normalise(s.TokenText()))
	}

	return words
}

//DictionaryMap implements a map-backed Dictionary.
type DictionaryMap struct {
	sync.RWMutex

	Map map[string]int64

	NextPrime int64
}

//NewDictionary returns a new DictionaryMap.
func NewDictionary() *DictionaryMap {
	return &DictionaryMap{
		Map:       make(map[string]int64),
		NextPrime: 2,
	}
}

func nextPrime(input int64) int64 {
	for i := big.NewInt(input + 1); true; i.Add(i, big.NewInt(1)) {
		if i.ProbablyPrime(0) {
			return i.Int64()
		}
	}

	panic(nil)
}

func normalise(word string) string {
	word = strings.ToLower(word)

	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	word, _, _ = transform.String(t, word)

	return word
}

//Add adds a word to the dictionary.
func (d *DictionaryMap) Add(words ...string) {
	d.Lock()
	defer d.Unlock()

	for _, word := range words {
		if _, ok := d.Map[word]; !ok {

			if strings.Contains(word, "/") {
				splits := strings.Split(word, "/")

				word = normalise(splits[0])

				d.Map[word] = d.NextPrime
				d.NextPrime = nextPrime(d.NextPrime)

				for _, synonym := range splits[1:] {
					d.Map[normalise(synonym)] = d.Map[word]
				}

				continue
			}

			d.Map[normalise(word)] = d.NextPrime
			d.NextPrime = nextPrime(d.NextPrime)
		}
	}
}

//AddSynonym adds a synonym to the dictionary.
func (d *DictionaryMap) AddSynonym(synonym string, of string) {
	d.Lock()
	defer d.Unlock()

	d.Map[normalise(synonym)] = d.Map[of]
}

//LookupWords implements Dictionary.LookupWords
func (d *DictionaryMap) LookupWords(results []int64, words ...string) error {
	d.RLock()
	defer d.RUnlock()

	for i, word := range words {
		results[i] = d.Map[word]
	}
	return nil
}
