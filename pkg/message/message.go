package message

import "github.com/soustify/sentinel/pkg/constant"

type (
	DiscordMessage struct {
		Title       string
		Description string
		Metadata    DiscordMetadaList
	}

	DiscordMetadaList []DiscordMetadata

	DiscordMetadata struct {
		Name  string
		Value string
	}
)

func (list DiscordMetadaList) GetResource() string {
	return list.extractByName(constant.Resource)
}

func (list DiscordMetadaList) GetSeverity() string {
	return list.extractByName(constant.Severity)
}

func (list DiscordMetadaList) extractByName(key string) string {
	if list != nil {
		for _, metadata := range list {
			if metadata.Name == key {
				return metadata.Value
			}
		}
	}
	return ""
}
