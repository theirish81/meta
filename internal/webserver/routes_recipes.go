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
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/theirish81/edjson"
	"github.com/theirish81/meta/internal/dto"
	"github.com/theirish81/meta/internal/persistence/domain"
)

func (s Server) SearchRecipes(ctx echo.Context, memory string, params dto.SearchRecipesParams) error {
	meta, err := s.Services.RecipeService.Search(ctx.Request().Context(), MustGetUser(ctx).Subject, memory, params.Tag, params.Q)
	return edjson.JSON[dto.Recipes](ctx, http.StatusOK, meta, err)
}

func (s Server) ListRecipesMemories(ctx echo.Context) error {
	memories, err := s.Services.RecipeService.Memories(ctx.Request().Context(), MustGetUser(ctx).Subject)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, memories)
}

func (s Server) CreateRecipe(ctx echo.Context, memory string) error {
	identity := MustGetUser(ctx)
	if !identity.CanWrite() {
		return echo.NewHTTPError(http.StatusForbidden)
	}
	body := dto.RecipeRequest{}
	if err := ctx.Bind(&body); err != nil {
		return err
	}
	meta := edjson.MustCopy[domain.Recipe](body)
	meta, err := s.Services.RecipeService.Create(ctx.Request().Context(), identity.Subject, memory, meta)
	return edjson.JSON[dto.Recipe](ctx, http.StatusCreated, meta, err)
}

func (s Server) DeleteRecipe(ctx echo.Context, memory string, recipeID openapi_types.UUID) error {
	identity := MustGetUser(ctx)
	if !identity.CanWrite() {
		return echo.NewHTTPError(http.StatusForbidden)
	}
	err := s.Services.RecipeService.Delete(ctx.Request().Context(), identity.Subject, memory, recipeID)
	if err != nil {
		return err
	}
	return ctx.NoContent(http.StatusNoContent)
}

func (s Server) UpdateRecipe(ctx echo.Context, memory string, recipeId openapi_types.UUID) error {
	body := dto.RecipeRequest{}
	if err := ctx.Bind(&body); err != nil {
		return err
	}
	recipe := edjson.MustCopy[domain.Recipe](body)
	recipe, err := s.Services.RecipeService.Update(ctx.Request().Context(), MustGetUser(ctx).Subject, memory, recipeId, recipe)
	return edjson.JSON[dto.Recipe](ctx, http.StatusOK, recipe, err)

}
