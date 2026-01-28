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

package domain

import (
	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
	"gorm.io/datatypes"
)

// Ollama vector size 2560

type KnowledgeChunk struct {
	ID         uuid.UUID                   `gorm:"primary_key;type:uuid;default:gen_random_uuid();<-:create"`
	Memory     string                      `gorm:"not null"`
	Document   string                      `gorm:"not null"`
	Tags       datatypes.JSONSlice[string] `gorm:"not null"`
	Chunk      string                      `gorm:"not null"`
	Embedding  pgvector.Vector             `gorm:"type:vector(3072); not null"`
	IdentityID string                      `gorm:"not null"`
	Distance   float64                     `gorm:"column:distance;<-:false;-:migration"`
}
