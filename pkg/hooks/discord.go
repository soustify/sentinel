package hooks

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/sirupsen/logrus"
	"runtime"
)

type (
	discordCriticalHook struct {
		resource    string
		queueUrl    string
		endpointUrl string
	}

	message struct {
		Title       string
		Description string
		Metadata    []metadata
	}

	metadata struct {
		Name  string
		Value string
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
	meta := make([]metadata, 0)

	meta = append(meta, metadata{
		Name:  "severity",
		Value: entry.Level.String(),
	})

	meta = append(meta, metadata{
		Name:  "resource",
		Value: hook.resource,
	})

	if ok {
		meta = append(meta, metadata{
			Name:  "file",
			Value: fmt.Sprintf("%s:%d", file, line),
		})
	}

	_, err := sendMessage(hook.queueUrl, hook.resource, hook.endpointUrl, message{
		Title:       fmt.Sprintf("Houve um erro na execução!"),
		Description: entry.Message,
		Metadata:    meta,
	})

	if err != nil {
		fmt.Printf("falha ao enviar a mensagem para a fila SQS: %v", err)
	}

	return err
}

func sendMessage(queueUrl, messageGroupId, endpointUrl string, messageBody message) (*sqs.SendMessageOutput, error) {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	var result *sqs.SendMessageOutput
	if err != nil {
		return nil, fmt.Errorf("falha ao carregar a configuração da AWS: %w", err)
	}

	sqsClient := sqs.NewFromConfig(cfg, func(options *sqs.Options) {
		options.BaseEndpoint = aws.String(endpointUrl)
		options.EndpointResolverV2 = sqsEndpointResolver{}
	})
	converted, err := json.Marshal(messageBody)

	if err != nil {
		return nil, err
	}

	messageDeduplicationId := fmt.Sprintf("%x", sha256.Sum256(converted))

	result, err = sqsClient.SendMessage(context.TODO(), &sqs.SendMessageInput{
		QueueUrl:               aws.String(queueUrl),
		MessageBody:            aws.String(string(converted)),
		MessageGroupId:         aws.String(messageGroupId),
		MessageDeduplicationId: aws.String(messageDeduplicationId),
	})
	if err != nil {
		return nil, fmt.Errorf("falha ao enviar a mensagem para a fila SQS: %w", err)
	}
	return result, nil
}
