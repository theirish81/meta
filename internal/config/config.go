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

package config

import (
	"errors"
	"io/fs"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	KbDistanceThreshold   float64 `mapstructure:"KB_DISTANCE_THRESHOLD" validate:"required,numeric,min=0,max=1"`
	MetaDistanceThreshold float64 `mapstructure:"META_DISTANCE_THRESHOLD" validate:"required,numeric,min=0,max=1"`
	DatabaseURL           string  `mapstructure:"DATABASE_URL" validate:"required"`
	EmbeddingModel        string  `mapstructure:"EMBEDDING_MODEL" validate:"required"`
	OllamaBaseURL         string  `mapstructure:"OLLAMA_BASE_URL" validate:"required,url"`
}

var Instance Config

func Init() error {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.SetDefault("EMBEDDING_MODEL", "")
	viper.SetDefault("KB_DISTANCE_THRESHOLD", "")
	viper.SetDefault("META_DISTANCE_THRESHOLD", "")
	viper.SetDefault("DATABASE_URL", "")
	viper.SetDefault("OLLAMA_BASE_URL", "")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError *fs.PathError
		if !errors.As(err, &configFileNotFoundError) {
			return err
		}
	}
	if err := viper.Unmarshal(&Instance); err != nil {
		return err
	}
	return validator.New().Struct(Instance)
}
