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
)

func (s Server) SearchKb(ctx echo.Context, memory string, params dto.SearchKbParams) error {
	kbs, err := s.Services.KnowledgeBaseService.Search(ctx.Request().Context(), MustGetUser(ctx).Subject, memory,
		params.Tag, params.Q)
	return edjson.JSON[dto.KnowledgeChunks](ctx, http.StatusOK, kbs, err)
}

func (s Server) SubmitDocument(ctx echo.Context, memory string, document string) error {
	dx := dto.Document{}
	if err := ctx.Bind(&dx); err != nil {
		return err
	}
	err := s.Services.KnowledgeBaseService.RecordDocument(ctx.Request().Context(), MustGetUser(ctx).Subject, memory,
		document, dx.Tags, dx.Content)
	if err != nil {
		return err
	}
	return ctx.NoContent(http.StatusCreated)
}

func (s Server) DeleteDocument(ctx echo.Context, memory string, document string) error {
	err := s.Services.KnowledgeBaseService.DeleteDocument(ctx.Request().Context(), MustGetUser(ctx).Subject, memory, document)
	if err != nil {
		return err
	}
	return ctx.NoContent(http.StatusNoContent)
}

func (s Server) ListDocuments(ctx echo.Context, memory string) error {
	docs, err := s.Services.KnowledgeBaseService.ListDocuments(ctx.Request().Context(), MustGetUser(ctx).Subject, memory)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, docs)
}

func (s Server) ListKbMemories(ctx echo.Context) error {
	memories, err := s.Services.KnowledgeBaseService.Memories(ctx.Request().Context(), MustGetUser(ctx).Subject)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, memories)
}
