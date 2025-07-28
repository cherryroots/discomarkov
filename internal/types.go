package internal

// Internal use structs

// User struct is a combined struct of the author of a message and all their messages
type User struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Messages []string `json:"messages"`
	Roles    []Role   `json:"roles"`
}

// UserMarkov struct is a struct of a user and their markov chains
type UserMarkov struct {
	Name          string         `json:"name"`
	WordFrequency map[string]int `json:"wordFrequency"`
	Ngrams        []*Ngram       `json:"ngrams"`
}

// Ngram struct is a struct of ngrams
type Ngram struct {
	N     int            `json:"n"`
	Grams map[string]int `json:"grams"` // words and their frequency
}

// Discord Chat Exporter JSON Format

// Export struct is the root of the JSON file
type Export struct {
	Guild        Guild     `json:"guild"`
	Channel      Channel   `json:"channel"`
	DateRange    DateRange `json:"dateRange"`
	ExportedAt   string    `json:"exportedAt"`
	Messages     []Message `json:"messages"`
	MessageCount int       `json:"messageCount"`
}

// Guild struct is the guild of the export
type Guild struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	IconURL string `json:"iconUrl"`
}

// Channel struct is the channel of the export
type Channel struct {
	ID         string  `json:"id"`
	Type       string  `json:"type"`
	CategoryID *string `json:"categoryId,omitempty"`
	Category   *string `json:"category,omitempty"`
	Name       string  `json:"name"`
	Topic      *string `json:"topic,omitempty"`
}

// DateRange struct is the date range of the export
type DateRange struct {
	After  *string `json:"after,omitempty"`
	Before *string `json:"before,omitempty"`
}

// Role struct is the role of the author of a message
type Role struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Color    *string `json:"color,omitempty"`
	Position int     `json:"position"`
}

// Author struct is the author of a message
type Author struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Discriminator string  `json:"discriminator"`
	Nickname      *string `json:"nickname,omitempty"`
	Color         *string `json:"color,omitempty"`
	IsBot         bool    `json:"bot"`
	Roles         []Role  `json:"roles"`
	AvatarURL     *string `json:"avatarUrl,omitempty"`
}

// Message struct is the message of the export
type Message struct {
	ID                 string  `json:"id"`
	Type               string  `json:"type"`
	Timestamp          string  `json:"timestamp"`
	TimestampEdited    *string `json:"timestampEdited,omitempty"`
	CallEndedTimestamp *string `json:"callEndedTimestamp,omitempty"`
	IsPinned           bool    `json:"pinned"`
	Author             Author  `json:"author"`
	Content            string  `json:"content"`
	// Omitting attachments, embeds, stickers, reactions, mentions, inlineEmojis
}
