package telemetry

import (
	"context"
	"fmt"
	"github.com/rhecosystemappeng/crda-cli/pkg/utils"
)

type contextKeyType string
type propertyKeyType string
type contextValue struct {
	props map[propertyKeyType]string
}

func (pt propertyKeyType) ToString() string {
	return fmt.Sprint(pt)
}

const (
	KeyManifest             propertyKeyType = "manifest"
	KeyExitCode             propertyKeyType = "exit-code"
	KeyClient               propertyKeyType = "client"
	KeyTotalVulnerabilities propertyKeyType = "total-vulnerabilities"
	KeyEcosystem            propertyKeyType = "ecosystem"
	KeySnykTokenAssociated  propertyKeyType = "snyk-token-associated"
	KeyJSonOutput           propertyKeyType = "json"
	KeyVerboseOutput        propertyKeyType = "verbose"
	KeySuccess              propertyKeyType = "success"
	KeyPlatform             propertyKeyType = "platform"
	KeyVersion              propertyKeyType = "version"
	KeyDuration             propertyKeyType = "duration"
	KeyError                propertyKeyType = "error"
)

var contextKey = contextKeyType("telemetry")

// GetContext is used for instantiating a value context storing telemetry client and properties
// invoke this functions once to get the context and pass it back when using telemetry
func GetContext(ctx context.Context) context.Context {
	utils.Logger.Debug("initiating context")
	val := contextValue{make(map[propertyKeyType]string)}
	return context.WithValue(ctx, contextKey, val)
}

// SetProperty is used to store a string, a bool, or an int as a telemetry event property
func SetProperty[T string | bool | int](ctx context.Context, key propertyKeyType, value T) {
	ctx.Value(contextKey).(contextValue).props[key] = fmt.Sprint(value)
}

// GetProperty is used to get a stored telemetry event property
func GetProperty(ctx context.Context, key propertyKeyType) (string, bool) {
	val, ok := ctx.Value(contextKey).(contextValue).props[key]
	return val, ok
}
