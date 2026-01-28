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

	"github.com/google/uuid"
	"github.com/theirish81/meta/internal/persistence/connection"
	"github.com/theirish81/meta/internal/persistence/domain"
)

type ObjectService struct {
	conn *connection.Connection
}

func NewObjectService() *ObjectService {
	return &ObjectService{
		conn: connection.Conn,
	}
}

func (s *ObjectService) InitTables(ctx context.Context) error {
	return s.conn.WithContext(ctx).AutoMigrate(&domain.Object{})
}

func (s *ObjectService) GetByName(ctx context.Context, ownerID string, memory string, name string) (domain.Object, error) {
	query := domain.Object{IdentityID: ownerID, Memory: memory, Name: name}
	object := domain.Object{}
	err := s.conn.WithContext(ctx).Model(&query).First(&object, query).Error
	return object, err
}

func (s *ObjectService) Upsert(ctx context.Context, ownerID string, memory string, obj domain.Object) error {
	obj.ID = uuid.New()
	obj.IdentityID = ownerID
	obj.Memory = memory
	if existingItem, err := s.GetByName(ctx, ownerID, memory, obj.Name); err != nil {
		return s.conn.WithContext(ctx).Create(&obj).Error
	} else {
		existingItem.Content = obj.Content
		return s.conn.WithContext(ctx).Model(&existingItem).Updates(&existingItem).Error
	}
}

func (s *ObjectService) Memories(ctx context.Context, ownerID string) ([]string, error) {
	memories := make([]string, 0)
	query := &domain.Object{IdentityID: ownerID}
	err := s.conn.WithContext(ctx).Model(&query).Where(&query).Distinct("memory").Pluck("memory", &memories).Error
	return memories, err
}

func (s *ObjectService) Delete(ctx context.Context, ownerID string, memory string, name string) error {
	return s.conn.WithContext(ctx).Delete(&domain.Object{}, "identity_id = ? AND memory = ? AND name = ?", ownerID, memory, name).Error
}

func (s *ObjectService) List(ctx context.Context, ownerID string, memory string) ([]domain.Object, error) {
	query := &domain.Object{IdentityID: ownerID, Memory: memory}
	objects := make([]domain.Object, 0)
	err := s.conn.WithContext(ctx).Model(&query).Where(&query).Find(&objects).Error
	return objects, err
}

func (s *ObjectService) DeleteByName(ctx context.Context, ownerID string, memory string, name string) error {
	query := domain.Object{IdentityID: ownerID, Memory: memory, Name: name}
	return s.conn.WithContext(ctx).Delete(&domain.Object{}, query).Error
}
