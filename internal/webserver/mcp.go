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
	"errors"
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
	"github.com/theirish81/meta/internal/persistence/domain"
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
			return toCallResult(res, "recipes_memories"), nil, err
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
			return toCallResult(edjson.MustCopy[dto.Recipes](res), "recipes"), nil, err
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
			return toCallResult(edjson.MustCopy[dto.KnowledgeChunks](res), "knowledge"), nil, err
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
			return toCallResult(res, "knowledge_memories"), nil, err
		})
	mcp.AddTool(mcpServer, toolObjectMemories,
		func(ctx context.Context, request *mcp.CallToolRequest, input objectParams) (*mcp.CallToolResult, any, error) {
			defer func() {
				if e := recover(); e != nil {
					log.Println(e)
				}
			}()
			claims := getMetaClaims(request.GetExtra().TokenInfo.Extra)
			res, err := services.Services.ObjectService.Memories(ctx, claims.Subject)
			return toCallResult(res, "memories"), nil, err
		})
	mcp.AddTool(mcpServer, toolObjectCreate,
		func(ctx context.Context, request *mcp.CallToolRequest, input objectParams) (*mcp.CallToolResult, any, error) {
			defer func() {
				if e := recover(); e != nil {
					log.Println(e)
				}
			}()
			claims := getMetaClaims(request.GetExtra().TokenInfo.Extra)
			err := services.Services.ObjectService.Upsert(ctx, claims.Subject, input.Memory, domain.Object{
				Name:        input.Name,
				Content:     input.Content,
				ContentType: input.ContentType,
			})
			if err != nil {
				return toCallResult("could not create object", "result"), nil, err
			}
			return toCallResult("object created", "result"), nil, nil
		})
	mcp.AddTool(mcpServer, toolObjectGetByName,
		func(ctx context.Context, request *mcp.CallToolRequest, input objectParams) (*mcp.CallToolResult, any, error) {
			defer func() {
				if e := recover(); e != nil {
					log.Println(e)
				}
			}()
			claims := getMetaClaims(request.GetExtra().TokenInfo.Extra)
			res, err := services.Services.ObjectService.GetByName(ctx, claims.Subject, input.Memory, input.Name)
			if err != nil {
				return toCallResult("could not create object", "result"), nil, err
			}
			return toCallResult(res, "object"), nil, nil
		})
	method := mcp.NewStreamableHTTPHandler(func(request *http.Request) *mcp.Server {
		return mcpServer
	}, nil)
	authMiddleware := auth.RequireBearerToken(func(ctx context.Context, token string, req *http.Request) (*auth.TokenInfo, error) {
		tokenObject, err := jwt.ParseWithClaims(token, &auth2.MetaClaims{}, func(token *jwt.Token) (any, error) {
			return loadPublicKey()
		})
		if err != nil {
			return nil, err
		}
		claims, ok := tokenObject.Claims.(*auth2.MetaClaims)
		if !ok {
			return nil, errors.New("invalid claims")
		}
		return &auth.TokenInfo{UserID: claims.Subject, Extra: map[string]any{"claims": claims}, Expiration: claims.ExpiresAt.Time}, err
	}, nil)
	s.E.Any("/mcp", echo.WrapHandler(authMiddleware(method)))
}

func toCallResult(data any, rootObjectName string) *mcp.CallToolResult {
	output := map[string]any{rootObjectName: data}
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

type objectParams struct {
	Memory      string `json:"memory"`
	Name        string `json:"name"`
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
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

var toolObjectMemories = &mcp.Tool{
	Name:        "meta_list_objects_memories",
	Description: "lists all objects memories and tags. Call this first to get the list of memories and tags to use in the meta_get_object_by_name tool.",
	InputSchema: &jsonschema.Schema{
		Type:       "object",
		Properties: map[string]*jsonschema.Schema{},
	},
}

var toolObjectCreate = &mcp.Tool{
	Name:        "meta_create_object",
	Description: "creates or updates a generic text object for later use, with a specific name into a specific memory slot.",
	InputSchema: &jsonschema.Schema{
		Type:     "object",
		Required: []string{"memory", "name", "content", "content_type"},
		Properties: map[string]*jsonschema.Schema{
			"memory": {
				Type: "string",
			},
			"name": {
				Type: "string",
			},
			"content": {
				Type: "string",
			},
			"content_type": {
				Type: "string",
				Enum: []any{
					"text/markdown",
					"text/plain",
					"application/json",
				},
			},
		},
	},
}

var toolObjectGetByName = &mcp.Tool{
	Name:        "meta_get_object_by_name",
	Description: "gets a generic text object by its name and memory slot.",
	InputSchema: &jsonschema.Schema{
		Type:     "object",
		Required: []string{"memory", "name"},
		Properties: map[string]*jsonschema.Schema{
			"memory": {
				Type: "string",
			},
			"name": {
				Type: "string",
			},
		},
	},
}
