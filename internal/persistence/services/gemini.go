/*
 * Copyright (C) 2026 Simone Pezzano
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package services

import (
	"context"
	"os"

	"cloud.google.com/go/auth/credentials"
	"github.com/samber/lo"
	"github.com/theirish81/meta/internal/config"
	"google.golang.org/genai"
)

type GeminiService struct {
	client *genai.Client
}

func NewGeminiService() (*GeminiService, error) {
	credsBytes, err := os.ReadFile("etc/keys/gemini.json")
	if err != nil {
		return nil, err
	}
	creds, err := credentials.DetectDefault(&credentials.DetectOptions{
		Scopes:          []string{"https://www.googleapis.com/auth/cloud-platform"},
		CredentialsJSON: credsBytes,
	})
	if err != nil {
		return nil, err
	}
	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		Project:     config.Instance.GeminiProjectID,
		Location:    "global",
		Credentials: creds,
		Backend:     genai.BackendVertexAI,
	})
	if err != nil {
		return nil, err
	}
	return &GeminiService{client: client}, nil
}

func (s *GeminiService) ExtractEmbeddings(input []string) ([]Embedding, error) {
	contents := lo.Map[string, *genai.Content](input, func(item string, index int) *genai.Content {
		return genai.NewContentFromText(item, "")
	})
	result, err := s.client.Models.EmbedContent(context.Background(),
		config.Instance.EmbeddingModel,
		contents,
		nil,
	)
	if err != nil {
		return nil, err
	}
	embeddings := make([]Embedding, len(input))
	for i := 0; i < len(input); i++ {
		embeddings[i] = Embedding{
			Text:   input[i],
			Vector: result.Embeddings[i].Values,
		}
	}
	return embeddings, nil
}
