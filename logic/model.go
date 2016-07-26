package logic

// BotResult represents response from HERE API
type BotResult struct {
	Location string
	Places   []*BotPlace
}

// BotPlace represent a single item from HERE API response
type BotPlace struct {
	Title        string
	Distance     int
	OpeningHours string
	URL          string
	HereURL      string
	IconURL      string
}
