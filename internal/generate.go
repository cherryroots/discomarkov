package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

var ngramTemplates = map[int]string{
	1:  "word",
	2:  "bigram",
	3:  "trigram",
	4:  "quadgram",
	5:  "pentagram",
	6:  "hexagram",
	7:  "heptagram",
	8:  "octagram",
	9:  "nonogram",
	10: "dectogram",
}

var ngramSlots = map[int]string{
	1:  "first-word",
	2:  "second-word",
	3:  "third-word",
	4:  "fourth-word",
	5:  "fifth-word",
	6:  "sixth-word",
	7:  "seventh-word",
	8:  "eighth-word",
	9:  "ninth-word",
	10: "tenth-word",
}

// UserMarkovsToClips converts a UserMarkovs to a clips deffacts file
func UserMarkovsToClips(userMarkovs map[string]*UserMarkov) (map[string]string, error) {
	// function input is a map of users with their ngram and a word frequencies
	// function output is a map of deffacts for each user
	deffacts := make(map[string]string)
	for _, userMarkov := range userMarkovs {
		log.Printf("Generating deffacts for %s...\n", userMarkov.Name)
		
		// Pre-calculate total size for better memory allocation
		estimatedSize := 1000 + len(userMarkov.WordFrequency)*50
		for _, ngram := range userMarkov.Ngrams {
			estimatedSize += len(ngram.Grams) * 100
		}
		
		var builder strings.Builder
		builder.Grow(estimatedSize)
		builder.WriteString("(deffacts message-data\n")
		
		log.Printf("Generating word frequencies for %s...\n", userMarkov.Name)
		// Process word frequencies
		for word, frequency := range userMarkov.WordFrequency {
			builder.WriteString("   (word (name \"")
			builder.WriteString(FormatClipsString(word))
			builder.WriteString("\") (frequency ")
			builder.WriteString(fmt.Sprintf("%d", frequency))
			builder.WriteString("))\n")
		}
		
		log.Printf("Generating ngram frequencies for %s...\n", userMarkov.Name)
		// Process ngrams with optimized string operations
		for _, ngram := range userMarkov.Ngrams {
			log.Printf("Generating ngram frequencies for %s...\n", ngramTemplates[ngram.N])
			ngramTemplate := ngramTemplates[ngram.N]
			
			// Pre-split all grams to avoid repeated string.Split calls
			type gramData struct {
				words     []string
				frequency int
			}
			gramList := make([]gramData, 0, len(ngram.Grams))
			for gram, frequency := range ngram.Grams {
				gramList = append(gramList, gramData{
					words:     strings.Split(gram, " "),
					frequency: frequency,
				})
			}
			
			// Build ngram entries efficiently
			for _, gd := range gramList {
				builder.WriteByte('(')
				builder.WriteString(ngramTemplate)
				builder.WriteByte(' ')
				
				for i, word := range gd.words {
					builder.WriteByte('(')
					builder.WriteString(ngramSlots[i+1])
					builder.WriteString(" \"")
					builder.WriteString(FormatClipsString(word))
					builder.WriteString("\") ")
				}
				
				builder.WriteString("(frequency ")
				builder.WriteString(fmt.Sprintf("%d", gd.frequency))
				builder.WriteString("))\n")
			}
		}
		
		builder.WriteString(")\n")
		deffacts[userMarkov.Name] = builder.String()
	}
	return deffacts, nil
}

// FormatClipsString formats a string for clips
func FormatClipsString(str string) string {
	// escape \ with \\"
	str = strings.ReplaceAll(str, `\`, `\\`)
	// escape " with \"
	str = strings.ReplaceAll(str, `"`, `\"`)
	return str
}

// WriteClipsDeffacts writes a map of deffacts to a file
func WriteClipsDeffacts(deffacts map[string]string, path string) error {
	for name, deffact := range deffacts {
		path := fmt.Sprintf("%s/%s.clp", path, name)
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			log.Printf("File %s exists, deleting...\n", path)
			os.Remove(path)
		}
		log.Printf("Writing deffacts to %s...\n", path)
		os.WriteFile(path, []byte(deffact), 0o644)
	}
	return nil
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
				f := strings.Split(filter, ":")
				if len(f) != 2 {
					log.Printf("Invalid filter: %s\n", filter)
					continue
				}
				fType := f[0]
				fVal := f[1]
				switch fType {
				case "uid":
					if !strings.Contains(user.ID, fVal) {
						continue UserLoop
					}
				case "rid":
					for _, role := range user.Roles {
						if !strings.Contains(role.ID, fVal) {
							continue UserLoop
						}
					}
				case "u":
					if !strings.Contains(user.Name, fVal) {
						continue UserLoop
					}
				case "r":
					for _, role := range user.Roles {
						if !strings.Contains(role.Name, fVal) {
							continue UserLoop
						}
					}
				default:
					return nil, fmt.Errorf("unknown filter type: %s", fType)
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
	if n > 10 {
		return nil, fmt.Errorf("n must be less than 10")
	}
	ngram := &Ngram{
		N:     n,
		Grams: make(map[string]int),
	}
	for _, message := range messages {
		words := strings.Fields(message)
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
