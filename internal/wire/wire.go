//go:build wireinject

package wire

import (
	"iivineri/internal/container"

	"github.com/google/wire"
)

func InitializeContainer() (*container.Container, error) {
	wire.Build(ProviderSet)
	return nil, nil
}