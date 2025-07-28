package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

// ParseAllFiles parses all files in a directory into a list of Export structs
func ParseAllFiles(input string) ([]*Export, error) {
	files, err := os.ReadDir(input)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no files found in %s", input)
	}
	log.Printf("Found %d files...\n", len(files))
	exports := make([]*Export, 0)
	var wg sync.WaitGroup
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}
		bytes, err := os.ReadFile(fmt.Sprintf("%s/%s", input, file.Name()))
		if err != nil {
			return nil, err
		}
		log.Printf("Parsing file %s...\n", file.Name())
		wg.Add(1)
		go func() {
			defer wg.Done()
			export, err := ParseFile(bytes)
			if err != nil {
				log.Printf("Failed to parse file %s: %v\n", file.Name(), err)
				return
			}
			log.Println("Parsed file", file.Name())
			exports = append(exports, export)
		}()
	}
	wg.Wait()
	log.Printf("Parsed %d files...\n", len(exports))
	return exports, nil
}

// ParseFile parses a file into an Export struct
func ParseFile(bytes []byte) (*Export, error) {
	var export Export
	if err := json.Unmarshal(bytes, &export); err != nil {
		return nil, err
	}
	return &export, nil
}

type userLock struct {
	mu sync.Mutex
	u  map[string]*User
}

// CollectUsers collects all users from a list of exports
func CollectUsers(exports []*Export) (map[string]*User, error) {
	users := &userLock{
		u: make(map[string]*User),
	}
	var wg sync.WaitGroup
	for _, export := range exports {
		log.Printf("Collecting users from %s...\n", export.Channel.Name)
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, message := range export.Messages {
				users.mu.Lock()
				if _, ok := users.u[message.Author.ID]; !ok {
					users.u[message.Author.ID] = &User{
						ID:       message.Author.ID,
						Name:     message.Author.Name,
						Messages: make([]string, 0),
						Roles:    message.Author.Roles,
					}
				}
				users.u[message.Author.ID].Messages = append(users.u[message.Author.ID].Messages, message.Content)
				users.mu.Unlock()
			}
		}()
	}
	wg.Wait()
	log.Printf("Users collected from %d exports...\n", len(exports))
	return users.u, nil
}

// WriteUsers writes a map of users to a file
func WriteUsers(users map[string]*User, path string) error {
	// check if file exists and delete it
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		log.Printf("File %s exists, deleting...\n", path)
		os.Remove(path)
	}
	log.Printf("Writing users to %s...\n", path)
	bytes, err := json.Marshal(users)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, bytes, 0o644)
	if err != nil {
		return err
	}
	return nil
}
