package message

type (
	DiscordMessage struct {
		Title       string
		Description string
		Metadata    []DiscordMetadata
	}

	DiscordMetadata struct {
		Name  string
		Value string
	}
)
