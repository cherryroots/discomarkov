package internal

// Internal use structs

type User struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Messages []string `json:"messages"`
	Roles    []string `json:"roles"`
}

// Discord Chat Exporter JSON Format

type Export struct {
	Guild        Guild     `json:"guild"`
	Channels     []Channel `json:"channels"`
	DateRange    DateRange `json:"dateRange"`
	ExportedAt   string    `json:"exportedAt"`
	Messages     []Message `json:"messages"`
	MessageCount int       `json:"messageCount"`
}

type Guild struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	IconUrl string `json:"iconUrl"`
}

type Channel struct {
	ID         string  `json:"id"`
	Type       string  `json:"type"`
	CategoryID *string `json:"categoryId,omitempty"`
	Category   *string `json:"category,omitempty"`
	Name       string  `json:"name"`
	Topic      *string `json:"topic,omitempty"`
}

type DateRange struct {
	After  *string `json:"after,omitempty"`
	Before *string `json:"before,omitempty"`
}

type Role struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Color    *string `json:"color,omitempty"`
	Position int     `json:"position"`
}

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
