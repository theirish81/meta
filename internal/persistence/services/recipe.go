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
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-set/v3"
	"github.com/pgvector/pgvector-go"
	"github.com/theirish81/meta/internal/dto"
	"github.com/theirish81/meta/internal/persistence/connection"
	"github.com/theirish81/meta/internal/persistence/domain"
)

type RecipeService struct {
	conn *connection.Connection
}

func NewRecipeService() *RecipeService {
	return &RecipeService{
		conn: connection.Conn,
	}
}

func (s *RecipeService) InitTables(ctx context.Context) error {
	err := s.conn.AutoMigrate(&domain.Recipe{})
	if err != nil {
		return err
	}
	return nil
}

func (s *RecipeService) Search(ctx context.Context, ownerID string, memory string, tags *[]string, q *string) ([]domain.Recipe, error) {
	res := make([]domain.Recipe, 0)
	query := &domain.Recipe{
		IdentityID: ownerID,
		Memory:     memory,
	}
	tx := s.conn.WithContext(ctx).Model(query).Where(query)
	if tags != nil {
		tagsJSON, _ := json.Marshal(tags)
		tx = tx.Where("tags @> ?::jsonb", tagsJSON)
	}

	if q != nil {
		embedding, err := Services.EmbeddingService.ExtractEmbedding(*q)
		if err != nil {
			return res, err
		}
		tx = tx.Select("*, embedding <=> ? as distance", pgvector.NewVector(embedding))
		tx = tx.Order("distance ASC")
		tx = tx.Limit(5)
	}
	err := tx.Find(&res).Error
	return res, err
}

func (s *RecipeService) Create(ctx context.Context, ownerID string, memory string, meta domain.Recipe) (domain.Recipe, error) {
	meta.ID = uuid.New()
	meta.IdentityID = ownerID
	meta.Memory = memory
	embedding, err := Services.EmbeddingService.ExtractEmbedding(meta.Name + ": " + meta.Description)
	if err != nil {
		return meta, err
	}
	meta.Embedding = pgvector.NewVector(embedding)
	err = s.conn.WithContext(ctx).Create(&meta).Error
	return meta, err
}

func (s *RecipeService) Delete(ctx context.Context, ownerID string, memory string, metaID uuid.UUID) error {
	query := &domain.Recipe{ID: metaID, IdentityID: ownerID, Memory: memory}
	return s.conn.WithContext(ctx).Model(query).Delete(query).Error
}

func (s *RecipeService) Memories(ctx context.Context, ownerID string) (dto.Memories, error) {
	memories := make(dto.Memories)
	mems := make([]string, 0)
	query := &domain.Recipe{IdentityID: ownerID}
	err := s.conn.WithContext(ctx).Model(query).Where(query).Distinct("memory").
		Pluck("memory", &mems).Error
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

func (s *RecipeService) Tags(ctx context.Context, ownerID string, memory string) ([]string, error) {
	res := set.NewTreeSet[string](func(s2 string, s string) int {
		return strings.Compare(s2, s)
	})
	var metas []domain.Recipe
	query := &domain.Recipe{IdentityID: ownerID, Memory: memory}
	err := s.conn.WithContext(ctx).Model(query).Where(query).Find(&metas).Error
	if err != nil {
		return res.Slice(), err
	}
	for _, m := range metas {
		res.InsertSlice(m.Tags)
	}
	return res.Slice(), nil
}

func (s *RecipeService) Show(ctx context.Context, ownerID string, memory string, recipeID uuid.UUID) (domain.Recipe, error) {
	query := &domain.Recipe{ID: recipeID, IdentityID: ownerID, Memory: memory}
	res := domain.Recipe{}
	err := s.conn.WithContext(ctx).Model(query).Where(query).First(&res, "id = ?", recipeID).Error
	return res, err
}

func (s *RecipeService) Update(ctx context.Context, ownerID string, memory string, recipeID uuid.UUID,
	recipe domain.Recipe) (domain.Recipe, error) {
	query := &domain.Recipe{ID: recipeID, IdentityID: ownerID, Memory: memory}
	err := s.conn.WithContext(ctx).Model(query).Where(query).Updates(&recipe).Error
	if err != nil {
		return recipe, err
	}
	return s.Show(ctx, ownerID, memory, recipeID)
}
