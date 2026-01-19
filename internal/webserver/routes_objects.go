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

package webserver

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/theirish81/edjson"
	"github.com/theirish81/meta/internal/dto"
	"github.com/theirish81/meta/internal/persistence/domain"
)

func (s Server) CreateObject(ctx echo.Context, memory string) error {
	body := dto.DataObject{}
	if err := ctx.Bind(&body); err != nil {
		return err
	}
	object := edjson.MustCopy[domain.Object](body)
	err := s.Services.ObjectService.Upsert(ctx.Request().Context(), MustGetUser(ctx).Subject, memory, object)
	if err != nil {
		return err
	}
	return ctx.NoContent(http.StatusCreated)
}

func (s Server) GetObjectByName(ctx echo.Context, memory string, params dto.GetObjectByNameParams) error {
	item, err := s.Services.ObjectService.GetByName(ctx.Request().Context(), MustGetUser(ctx).Subject, memory, params.Name)
	return edjson.JSON[dto.DataObject](ctx, http.StatusOK, item, err)
}

func (s Server) ListObjectsMemories(ctx echo.Context) error {
	memories, err := s.Services.ObjectService.Memories(ctx.Request().Context(), MustGetUser(ctx).Subject)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, memories)
}

func (s Server) ListObjects(ctx echo.Context, memory string) error {
	objects, err := s.Services.ObjectService.List(ctx.Request().Context(), MustGetUser(ctx).Subject, memory)
	return edjson.JSON[[]dto.DataObject](ctx, http.StatusOK, objects, err)
}

func (s Server) DeleteObjectByName(ctx echo.Context, memory string, params dto.DeleteObjectByNameParams) error {
	err := s.Services.ObjectService.DeleteByName(ctx.Request().Context(), MustGetUser(ctx).Subject, memory, params.Name)
	if err != nil {
		return err
	}
	return ctx.NoContent(http.StatusNoContent)
}
