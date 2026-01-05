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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/theirish81/meta/internal/config"
)

type EmbeddingService struct {
	baseURL string
	client  *http.Client
}

type OllamaEmbeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Params Params `json:"params"`
}

type Params struct {
	OutputDim int `json:"output_dim"`
}

func (r *OllamaEmbeddingRequest) Reader() *bytes.Reader {
	data, _ := json.Marshal(r)
	return bytes.NewReader(data)
}

type OllamaEmbeddingResponse struct {
	Embedding []float32 `json:"embedding"`
}

func NewEmbeddingService(baseURL string) *EmbeddingService {
	return &EmbeddingService{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (c *EmbeddingService) ExtractEmbedding(text string) ([]float32, error) {
	request := OllamaEmbeddingRequest{
		Model:  config.Instance.EmbeddingModel,
		Prompt: text,
	}
	resp, err := c.client.Post(fmt.Sprintf("%s/api/embeddings", c.baseURL), "application/json", request.Reader())
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	if resp.Body != nil {
		defer func() {
			_ = resp.Body.Close()
		}()
		data, _ := io.ReadAll(resp.Body)
		response := OllamaEmbeddingResponse{}
		_ = json.Unmarshal(data, &response)
		return response.Embedding, nil

	} else {
		return nil, errors.New("no response body")
	}

}
