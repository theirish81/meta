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

package cmd

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"slices"
	"time"

	jwt2 "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/kataras/iris/v12/middleware/jwt"
	"github.com/oasdiff/yaml"
	"github.com/spf13/cobra"
	"github.com/theirish81/meta/internal/auth"
)

var RootCmd = &cobra.Command{
	Use:   "meta",
	Short: "Meta CLI for administration commands",
}

var (
	subjectParam     string
	emailParam       string
	permissionsParam string
)

var key = &cobra.Command{
	Use:   "key",
	Short: "Generate a user key",
	Run: func(cmd *cobra.Command, args []string) {
		if subjectParam == "" {
			subjectParam = uuid.NewString()
		}
		if !slices.Contains([]string{"read", "write"}, permissionsParam) {
			cmd.PrintErrln("Permissions must be either read or write")
			return
		}
		data, err := os.ReadFile("etc/keys/private.pem")
		if err != nil {
			cmd.PrintErrln(err)
			return
		}
		pk, err := jwt.ParsePrivateKeyRSA(data)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}
		exp := jwt2.NumericDate{Time: time.Now().Add(24 * 30 * 12 * 10 * time.Hour)}
		claims := auth.MetaClaims{
			RegisteredClaims: jwt2.RegisteredClaims{
				Subject:   subjectParam,
				ExpiresAt: &exp,
			},
			Email:       emailParam,
			Nonce:       generateNonce(16),
			Permissions: permissionsParam,
		}
		token, _ := jwt.Sign(jwt.RS512, pk, claims)
		fmt.Println("-----BEGIN TOKEN-----\n" + string(token) + "\n-----END TOKEN-----")
		claimsData, _ := yaml.Marshal(claims)
		fmt.Println("-----BEGIN DATA-----\n" + string(claimsData) + "\n-----END DATA-----")
	},
}

func generateNonce(length int) string {
	bytes := make([]byte, length)
	_, _ = rand.Read(bytes)
	return base64.RawURLEncoding.EncodeToString(bytes)
}
func init() {
	RootCmd.AddCommand(key)
	key.Flags().StringVarP(&subjectParam, "subject", "s", "", "The subject of the key")
	key.Flags().StringVarP(&emailParam, "email", "e", "", "The email of the user")
	key.Flags().StringVarP(&permissionsParam, "permissions", "p", "", "The permissions of the user (read or write)")
	_ = key.MarkFlagRequired("email")
	_ = key.MarkFlagRequired("permissions")
}
