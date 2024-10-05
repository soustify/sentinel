package hooks

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/soustify/sentinel/pkg/message"
	"github.com/soustify/sentinel/pkg/publisher"
	"runtime"
)

type (
	discordCriticalHook struct {
		resource    string
		queueUrl    string
		endpointUrl string
	}
)

func NewDiscordCriticalHook(resource, queueUrl, endpointUrl string) logrus.Hook {
	return discordCriticalHook{
		resource:    resource,
		queueUrl:    queueUrl,
		endpointUrl: endpointUrl,
	}
}

func (p discordCriticalHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
	}
}

func (hook discordCriticalHook) Fire(entry *logrus.Entry) error {
	_, file, line, ok := runtime.Caller(9)
	meta := make([]message.DiscordMetadata, 0)

	meta = append(meta, message.DiscordMetadata{
		Name:  "severity",
		Value: entry.Level.String(),
	})

	meta = append(meta, message.DiscordMetadata{
		Name:  "resource",
		Value: hook.resource,
	})

	if ok {
		meta = append(meta, message.DiscordMetadata{
			Name:  "file",
			Value: fmt.Sprintf("%s:%d", file, line),
		})
	}

	_, err := publisher.Publish(hook.queueUrl, hook.resource, hook.endpointUrl, message.DiscordMessage{
		Title:       fmt.Sprintf("Houve um erro na execução!"),
		Description: entry.Message,
		Metadata:    meta,
	})

	if err != nil {
		fmt.Printf("falha ao enviar a mensagem para a fila SQS: %v", err)
	}

	return err
}
