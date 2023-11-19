package core

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

type CodeGenerator interface {
}

type InitializeParameters struct {
	Model  string
	Client *openai.Client
}

type codeGeneratorModel struct {
	InitializeParameters
}

func NewCodeGenerator(p InitializeParameters) (*codeGeneratorModel, error) {
	m := &codeGeneratorModel{p}
	return m, nil
}

func (m *codeGeneratorModel) LoadContext(ctx context.Context) error {
	return nil
}

func (m *codeGeneratorModel) Predict(ctx context.Context) (string, error) {
	return "success", nil
}
