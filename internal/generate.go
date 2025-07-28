package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

/*

This is the structure of the clips deftemplate for ngrams
(deftemplate word
   (slot name)
   (slot frequency (type INTEGER)))

(deftemplate bigram
   (slot first-word)
   (slot second-word)
   (slot frequency (type INTEGER)))

(deftemplate trigram
   (slot first-word)
   (slot second-word)
   (slot third-word)
   (slot frequency (type INTEGER)))

For the deffacts
(deffacts message-data
   ; Word frequencies
   (word (name hey) (frequency 2))
   ;  and so on for all words

   ; Bigram frequencies
   (bigram (first-word hey) (second-word what) (frequency 1))
   (bigram (first-word what) (second-word s) (frequency 1))
   (bigram (first-word s) (second-word up) (frequency 1))
   (bigram (first-word lol) (second-word no) (frequency 1))
   ;  and so on for all bigrams
)

*/

// UserMarkovsToClipsFiles converts a UserMarkovs to a clips deffacts file
func UserMarkovsToClipsFiles(userMarkovs map[string]*UserMarkov) (string, error) {
	_ = userMarkovs
	return "", nil
}

// WriteUserMarkovs writes a map of user markovs to a file
func WriteUserMarkovs(userMarkovs map[string]*UserMarkov, path string) error {
	// check if file exists and delete it
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		log.Printf("File %s exists, deleting...\n", path)
		os.Remove(path)
	}
	log.Printf("Writing users to %s...\n", path)
	bytes, err := json.Marshal(userMarkovs)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, bytes, 0o644)
	if err != nil {
		return err
	}
	return nil
}

// ReadUserMarkovs reads a map of user markovs from a file
func ReadUserMarkovs(path string) (map[string]*UserMarkov, error) {
	log.Printf("Reading users from %s...\n", path)
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var userMarkovs map[string]*UserMarkov
	if err := json.Unmarshal(bytes, &userMarkovs); err != nil {
		return nil, err
	}
	return userMarkovs, nil
}

// GenerateMarkovChains generates markov chains from a list of users
func GenerateMarkovChains(users map[string]*User, filters []string) (map[string]*UserMarkov, error) {
	userMarkovs := make(map[string]*UserMarkov)
	log.Printf("Generating markov chains from %d users...\n", len(users))
UserLoop:
	for _, user := range users {
		if len(filters) > 0 {
			for _, filter := range filters {
				if strings.HasPrefix(filter, "uid:") {
					if !strings.Contains(user.ID, filter[4:]) {
						continue UserLoop
					}
				}
				if strings.HasPrefix(filter, "rid:") {
					for _, role := range user.Roles {
						if !strings.Contains(role.ID, filter[4:]) {
							continue UserLoop
						}
					}
				}
				if strings.HasPrefix(filter, "u:") {
					if !strings.Contains(user.Name, filter[2:]) {
						continue UserLoop
					}
				}
				if strings.HasPrefix(filter, "r:") {
					for _, role := range user.Roles {
						if !strings.Contains(role.Name, filter[2:]) {
							continue UserLoop
						}
					}
				}
			}
		}
		wordFrequency, err := GetWordFrequency(user.Messages)
		if err != nil {
			return nil, err
		}
		userMarkov := &UserMarkov{
			Name:          user.Name,
			WordFrequency: wordFrequency,
			Ngrams:        make([]*Ngram, 0),
		}

		bigram, err := GetNgram(user.Messages, 2)
		if err != nil {
			return nil, err
		}
		trigram, err := GetNgram(user.Messages, 3)
		if err != nil {
			return nil, err
		}
		userMarkov.Ngrams = append(userMarkov.Ngrams, bigram)
		userMarkov.Ngrams = append(userMarkov.Ngrams, trigram)
		userMarkovs[user.ID] = userMarkov
	}
	log.Printf("Generated markov chains from %d users...\n", len(userMarkovs))
	log.Printf("Filtered out %d users...\n", len(users)-len(userMarkovs))
	return userMarkovs, nil
}

// GetNgram creates an ngram from a list of messages
func GetNgram(messages []string, n int) (*Ngram, error) {
	if n <= 0 {
		return nil, fmt.Errorf("n must be greater than 0")
	}
	ngram := &Ngram{
		N:     n,
		Grams: make(map[string]int),
	}
	for _, message := range messages {
		words := strings.FieldsFunc(message, func(r rune) bool {
			return r == ' ' || r == '\n' || r == '\t' || r == '\r'
		})
		if len(words) < n {
			continue
		}
		for i := 0; i <= len(words)-n; i++ {
			gram := strings.Join(words[i:i+n], " ")
			ngram.Grams[gram]++
		}
	}
	return ngram, nil
}

// GetWordFrequency creates a word frequency from a list of messages
func GetWordFrequency(messages []string) (map[string]int, error) {
	wordFrequency := make(map[string]int)
	for _, message := range messages {
		words := strings.FieldsSeq(message)
		for word := range words {
			wordFrequency[word]++
		}
	}
	return wordFrequency, nil
}
