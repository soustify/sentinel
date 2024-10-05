package publisher

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/soustify/sentinel/pkg/message"
)

func Publish(queueUrl, messageGroupId, endpointUrl string, messageBody message.DiscordMessage) (*sqs.SendMessageOutput, error) {
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
