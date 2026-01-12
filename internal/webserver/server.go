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
	"crypto/rsa"
	"encoding/base64"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/theirish81/echosec"
	"github.com/theirish81/meta/internal/auth"
	"github.com/theirish81/meta/internal/persistence/services"
)

var _ ServerInterface = (*Server)(nil)

type Server struct {
	Services services.ServiceRegistry
	E        *echo.Echo
}

func NewServer() (*Server, error) {
	server := Server{
		Services: services.Services,
		E:        echo.New(),
	}
	publicKey, err := loadPublicKey()
	if err != nil {
		return &server, err
	}
	server.E.Static("/web", "/web")
	grp := server.E.Group("/api/v1")
	grp.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:    publicKey,
		SigningMethod: "RS512",
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return &auth.MetaClaims{}
		},
	}))
	sw, _ := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	cfg, _ := echosec.NewOApiConfig(sw, map[string]echosec.OApiValidationFunc{
		"can_read": func(c echo.Context, params []string) error {
			return nil
		},
		"can_write": func(c echo.Context, params []string) error {
			if !MustGetUser(c).CanWrite() {
				return echo.NewHTTPError(403, "not authorized")
			}
			return nil
		},
	}, true)
	grp.Use(echosec.WithOpenApiConfig(cfg))
	server.E.HTTPErrorHandler = func(err error, c echo.Context) {
		_ = c.JSON(500, map[string]string{"error": err.Error()})
	}
	RegisterHandlers(grp, &server)
	server.initMCP()
	return &server, nil
}

func (s Server) Run() error {
	return s.E.Start(":8080")
}

func loadPublicKey() (*rsa.PublicKey, error) {
	publicKey, err := os.ReadFile("etc/keys/public.pem")
	if err != nil {
		panic(err)
	}
	return jwt.ParseRSAPublicKeyFromPEM(publicKey)
}
