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
	"math"
	"net/http"

	"github.com/theirish81/meta/internal/config"
)

type OllamaService struct {
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

func NewOllamaService(baseURL string) *OllamaService {
	return &OllamaService{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (c *OllamaService) ExtractEmbeddings(input []string) ([]Embedding, error) {
	embeddings := make([]Embedding, 0)
	var mainErr error
	for _, text := range input {
		func() {
			request := OllamaEmbeddingRequest{
				Model:  config.Instance.EmbeddingModel,
				Prompt: text,
			}
			resp, err := c.client.Post(fmt.Sprintf("%s/api/embeddings", c.baseURL), "application/json", request.Reader())
			if err != nil {
				mainErr = err
				return
			}
			if resp.Body != nil {
				defer func() {
					_ = resp.Body.Close()
				}()
				data, _ := io.ReadAll(resp.Body)
				response := OllamaEmbeddingResponse{}
				_ = json.Unmarshal(data, &response)
				embedding := response.Embedding
				c.Normalize(embedding)
				embedding = c.PadQwenToGemini(embedding)
				embeddings = append(embeddings, Embedding{Text: text, Vector: embedding})
			} else {
				mainErr = errors.New("no response body")
				return
			}
		}()
	}
	return embeddings, mainErr
}

func (c *OllamaService) Normalize(v []float32) {
	var norm float32
	for _, x := range v {
		norm += x * x
	}
	norm = float32(math.Sqrt(float64(norm)))
	if norm == 0 {
		return
	}
	for i := range v {
		v[i] /= norm
	}
}

func (c *OllamaService) PadQwenToGemini(src []float32) []float32 {
	dst := make([]float32, 3072)
	copy(dst, src)
	return dst
}
