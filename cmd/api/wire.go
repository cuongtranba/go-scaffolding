//go:build wireinject
// +build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	wireproviders "github.com/yourusername/go-scaffolding/internal/wire"
)

// initializeApp initializes the application with all dependencies
func initializeApp(configPath string) (*gin.Engine, func(), error) {
	wire.Build(wireproviders.ProviderSet)
	return nil, nil, nil
}
