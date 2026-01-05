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
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-set/v3"
	"github.com/pgvector/pgvector-go"
	"github.com/samber/lo"
	"github.com/theirish81/meta/internal/config"
	"github.com/theirish81/meta/internal/dto"
	"github.com/theirish81/meta/internal/persistence/connection"
	"github.com/theirish81/meta/internal/persistence/domain"
	"github.com/tmc/langchaingo/textsplitter"
	"gorm.io/datatypes"
)

type KnowledgeBaseService struct {
	conn *connection.Connection
}

func NewKnowledgeBaseService() *KnowledgeBaseService {
	return &KnowledgeBaseService{
		conn: connection.Conn,
	}
}

func (s *KnowledgeBaseService) InitTables(ctx context.Context) error {
	return s.conn.WithContext(ctx).AutoMigrate(&domain.KnowledgeChunk{})
}

func (s *KnowledgeBaseService) Search(ctx context.Context, ownerID string, memory string, tags *[]string, q string) ([]domain.KnowledgeChunk, error) {
	res := make([]domain.KnowledgeChunk, 0)
	query := &domain.KnowledgeChunk{
		IdentityID: ownerID,
		Memory:     memory,
	}
	tx := s.conn.WithContext(ctx).Model(query).Where(query)
	if tags != nil {
		tagsJSON, _ := json.Marshal(tags)
		tx = tx.Where("tags @> ?::jsonb", tagsJSON)
	}
	embedding, err := Services.EmbeddingService.ExtractEmbedding(q)
	if err != nil {
		return res, err
	}
	tx = tx.Select("*, embedding <=> ? as distance", pgvector.NewVector(embedding))
	tx = tx.Order("distance ASC")
	tx = tx.Limit(15)
	if err = tx.Find(&res).Error; err != nil {
		return res, err
	}
	res = lo.Filter(res, func(item domain.KnowledgeChunk, index int) bool {
		return item.Distance < config.Instance.KbDistanceThreshold
	})
	return res, nil
}

func (s *KnowledgeBaseService) RecordDocument(ctx context.Context, ownerID string, memory string, document string,
	tags []string, content string) error {
	var splitter textsplitter.TextSplitter
	if filepath.Ext(document) == ".md" {
		splitter = textsplitter.NewMarkdownTextSplitter(
			textsplitter.WithChunkSize(500),
			textsplitter.WithHeadingHierarchy(true),
			textsplitter.WithChunkOverlap(50))
	} else {
		splitter = textsplitter.NewRecursiveCharacter(
			textsplitter.WithChunkSize(500),
			textsplitter.WithChunkOverlap(50))
	}

	chunks, err := splitter.SplitText(content)
	if err != nil {
		return err
	}

	for _, chunk := range chunks {
		kc := domain.KnowledgeChunk{
			Memory:     memory,
			Document:   document,
			Tags:       tags,
			Chunk:      chunk,
			IdentityID: ownerID,
		}
		embedding, err := Services.EmbeddingService.ExtractEmbedding(kc.Chunk)
		if err != nil {
			return err
		}
		kc.Embedding = pgvector.NewVector(embedding)
		if err := s.conn.WithContext(ctx).Create(&kc).Error; err != nil {
			return err
		}
	}
	return nil
}

func (s *KnowledgeBaseService) DeleteDocument(ctx context.Context, ownerID string, memory string, document string) error {
	return s.conn.WithContext(ctx).Delete(&domain.KnowledgeChunk{}, "identity_id = ? AND memory = ? AND document = ?", ownerID, memory, document).Error
}

func (s *KnowledgeBaseService) ListDocuments(ctx context.Context, ownerID string, memory string) ([]string, error) {
	var documents []string
	query := &domain.KnowledgeChunk{
		IdentityID: ownerID,
		Memory:     memory,
	}
	err := s.conn.WithContext(ctx).Model(&query).Distinct("document").Where(query).
		Pluck("document", &documents).Error
	return documents, err
}

func (s *KnowledgeBaseService) Tags(ctx context.Context, ownerID string, memory string) ([]string, error) {
	res := set.NewTreeSet[string](func(s2 string, s string) int {
		return strings.Compare(s2, s)
	})
	tags := make([]datatypes.JSONSlice[string], 0)
	query := &domain.KnowledgeChunk{
		IdentityID: ownerID,
		Memory:     memory,
	}
	err := s.conn.WithContext(ctx).Model(query).Where(query).Pluck("tags", &tags).Error
	if err != nil {
		return res.Slice(), err
	}
	for _, m := range tags {
		res.InsertSlice(m)
	}
	return res.Slice(), nil
}

func (s *KnowledgeBaseService) Memories(ctx context.Context, ownerID string) (dto.Memories, error) {
	memories := make(dto.Memories)
	mems := make([]string, 0)
	query := &domain.KnowledgeChunk{
		IdentityID: ownerID,
	}
	err := s.conn.WithContext(ctx).Model(&query).Where(query).Distinct("memory").Pluck("memory", &mems).Error
	if err != nil {
		return memories, err
	}
	for _, m := range mems {
		tags, err := s.Tags(ctx, ownerID, m)
		if err != nil {
			return memories, err
		}
		memories[m] = dto.Memory{
			AvailableTags: tags,
		}
	}
	return memories, nil
}
