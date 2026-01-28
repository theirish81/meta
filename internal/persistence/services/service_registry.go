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
	"errors"

	"github.com/theirish81/meta/internal/config"
	"github.com/theirish81/meta/internal/persistence/connection"
	"gorm.io/gorm/logger"
)

type ServiceRegistry struct {
	EmbeddingService     EmbeddingService
	RecipeService        *RecipeService
	KnowledgeBaseService *KnowledgeBaseService
	ObjectService        *ObjectService
}

var Services ServiceRegistry

func Init() error {
	if err := connection.InitConnection(config.Instance.DatabaseURL, logger.Info); err != nil {
		return err
	}
	var err error
	switch config.Instance.EmbeddingService {
	case "ollama":
		Services.EmbeddingService = NewOllamaService(config.Instance.OllamaBaseURL)
	case "gemini":
		Services.EmbeddingService, err = NewGeminiService()
	default:
		err = errors.New("embedding service not selected")
	}
	if err != nil {
		return err
	}

	metaService := NewRecipeService()
	if err := metaService.InitTables(context.Background()); err != nil {
		return err
	}
	Services.RecipeService = metaService

	knowledgeService := NewKnowledgeBaseService()
	if err := knowledgeService.InitTables(context.Background()); err != nil {
		return err
	}
	Services.KnowledgeBaseService = knowledgeService

	objectService := NewObjectService()
	if err := objectService.InitTables(context.Background()); err != nil {
		return err
	}
	Services.ObjectService = objectService
	return nil
}
