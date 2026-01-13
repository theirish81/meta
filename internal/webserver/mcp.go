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
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/labstack/echo/v4"
	"github.com/modelcontextprotocol/go-sdk/auth"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/theirish81/edjson"
	auth2 "github.com/theirish81/meta/internal/auth"
	"github.com/theirish81/meta/internal/dto"
	"github.com/theirish81/meta/internal/persistence/services"
)

func (s Server) initMCP() {
	mcpServer := mcp.NewServer(&mcp.Implementation{Name: "META", Version: "v1.0.0"}, nil)
	mcp.AddTool(mcpServer, toolRecipesMemories,
		func(ctx context.Context, request *mcp.CallToolRequest, args any) (*mcp.CallToolResult, any, error) {
			defer func() {
				if e := recover(); e != nil {
					log.Println(e)
				}
			}()
			claims := getMetaClaims(request.GetExtra().TokenInfo.Extra)
			res, err := services.Services.RecipeService.Memories(ctx, claims.Subject)
			return toCallResult(res), nil, err
		})
	mcp.AddTool(mcpServer, toolRecipeSearch,
		func(ctx context.Context, request *mcp.CallToolRequest, args recipeParams) (*mcp.CallToolResult, any, error) {
			defer func() {
				if e := recover(); e != nil {
					log.Println(e)
				}
			}()
			claims := getMetaClaims(request.GetExtra().TokenInfo.Extra)
			res, err := services.Services.RecipeService.Search(ctx, claims.Subject, args.Memory, &args.Tag, &args.Q)
			return toCallResult(edjson.MustCopy[dto.Recipes](res)), nil, err
		})

	mcp.AddTool(mcpServer, toolKnowledgeSearch,
		func(ctx context.Context, request *mcp.CallToolRequest, args kbParams) (*mcp.CallToolResult, any, error) {
			defer func() {
				if e := recover(); e != nil {
					log.Println(e)
				}
			}()
			claims := getMetaClaims(request.GetExtra().TokenInfo.Extra)
			res, err := services.Services.KnowledgeBaseService.Search(ctx, claims.Subject, args.Memory, &args.Tag, args.Q)
			return toCallResult(edjson.MustCopy[dto.KnowledgeChunks](res)), nil, err
		})
	mcp.AddTool(mcpServer, toolKnowledgeMemories,
		func(ctx context.Context, request *mcp.CallToolRequest, input any) (*mcp.CallToolResult, any, error) {
			defer func() {
				if e := recover(); e != nil {
					log.Println(e)
				}
			}()
			claims := getMetaClaims(request.GetExtra().TokenInfo.Extra)
			res, err := services.Services.KnowledgeBaseService.Memories(ctx, claims.Subject)
			return toCallResult(res), nil, err
		})
	method := mcp.NewStreamableHTTPHandler(func(request *http.Request) *mcp.Server {
		return mcpServer
	}, nil)
	authMiddleware := auth.RequireBearerToken(func(ctx context.Context, token string, req *http.Request) (*auth.TokenInfo, error) {
		tokenObject, err := jwt.ParseWithClaims(token, &auth2.MetaClaims{}, func(token *jwt.Token) (any, error) {
			return loadPublicKey()
		})
		claims := tokenObject.Claims.(*auth2.MetaClaims)
		return &auth.TokenInfo{UserID: claims.Subject, Extra: map[string]any{"claims": claims}, Expiration: claims.ExpiresAt.Time}, err
	}, nil)
	s.E.Any("/mcp", echo.WrapHandler(authMiddleware(method)))
}

func toCallResult(data any) *mcp.CallToolResult {
	output := map[string]any{"result": data}
	content, _ := json.Marshal(output)
	return &mcp.CallToolResult{StructuredContent: output, Content: []mcp.Content{
		&mcp.TextContent{
			Text: string(content),
		},
	}}
}

func getMetaClaims(m map[string]any) *auth2.MetaClaims {
	return m["claims"].(*auth2.MetaClaims)
}

type kbParams struct {
	Memory string   `json:"memory"`
	Tag    []string `json:"tag"`
	Q      string   `json:"q"`
}

type recipeParams struct {
	Memory string   `json:"memory"`
	Tag    []string `json:"tag"`
	Q      string   `json:"q"`
}

var toolKnowledgeSearch = &mcp.Tool{
	Name:        "meta_search_knowledge",
	Description: "searches  knowledge. Knowledge records contain knowledge that is useful as-is to the user. Call this for most topics",
	InputSchema: &jsonschema.Schema{
		Type:     "object",
		Required: []string{"memory", "tag", "q"},
		Properties: map[string]*jsonschema.Schema{
			"memory": {
				Type:        "string",
				Description: "the memory slot to use. This is mandatory. Use only memories that have either been provided by the user themselves, or returned by the meta_list_knowledge_memories call",
			},
			"tag": {
				Type:        "array",
				Description: "tags to identify the topic to be searched. Use only tags that have either been provided by the user themselves, or returned by the meta_list_knowledge_memories call",
				Items: &jsonschema.Schema{
					Type: "string",
				},
			},
			"q": {
				Type:        "string",
				Description: "the user prompt that led to this tool execution",
			},
		},
	},
}
var toolKnowledgeMemories = &mcp.Tool{
	Name:        "meta_list_knowledge_memories",
	Description: "lists all memory slots and their tags. Call this first to get the list of memories and tags to use in the meta_search_knowledge tool.",
	InputSchema: &jsonschema.Schema{
		Type:       "object",
		Properties: map[string]*jsonschema.Schema{},
	},
}

var toolRecipeSearch = &mcp.Tool{
	Name:        "meta_search_recipes",
	Description: "searches for recipes. Each recipe contains a manual to help the assistant perform complex tasks",
	InputSchema: &jsonschema.Schema{
		Type:     "object",
		Required: []string{"memory", "tag", "q"},
		Properties: map[string]*jsonschema.Schema{
			"memory": {
				Type:        "string",
				Description: "the memory slot to use. Use only memories that have either been provided by the user themselves, or returned by the meta_list_recipes_memories call",
			},
			"tag": {
				Type:        "array",
				Description: "tags to identify the task to be accomplished. Use only tags that have either been provided by the user themselves, or returned by the meta_list_recipes_memories call",
				Items: &jsonschema.Schema{
					Type: "string",
				},
			},
			"q": {
				Type:        "string",
				Description: "the user prompt that led to this tool execution",
			},
		},
	},
}
var toolRecipesMemories = &mcp.Tool{
	Name:        "meta_list_recipes_memories",
	Description: "lists all recipes memories and tags. Call this first to get the list of memories  and tags to use in the meta_search_recipes tool.",
	InputSchema: &jsonschema.Schema{
		Type:       "object",
		Properties: map[string]*jsonschema.Schema{},
	},
}
