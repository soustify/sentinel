package parameters

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type (
	ParameterStoreRepository struct {
		client *ssm.Client
	}
)

func NewParameterStoreRepository(cfg aws.Config) ParameterStoreRepository {
	return ParameterStoreRepository{
		client: ssm.NewFromConfig(cfg),
	}
}

func (p ParameterStoreRepository) GetParameter(ctx context.Context, name string) (string, error) {
	input := &ssm.GetParameterInput{
		Name: aws.String(name),
	}

	output, err := p.client.GetParameter(ctx, input)
	if err != nil {
		return "", fmt.Errorf("erro ao buscar parâmetro %s: %w", name, err)
	}

	if output.Parameter == nil || output.Parameter.Value == nil {
		return "", fmt.Errorf("parâmetro %s não encontrado ou sem valor", name)
	}

	return aws.ToString(output.Parameter.Value), nil
}

func (p ParameterStoreRepository) GetDataGatewayTokenPermission(name string) (string, error) {
	ctx := context.Background()
	token, err := p.GetParameter(ctx, "/data-gateway/tokens/"+name)
	if err != nil {
		return "", fmt.Errorf("erro ao buscar parâmetro %s: %w", name, err)
	}
	if token == "" {
		return "", fmt.Errorf("token vazio")
	}
	return token, nil
}
