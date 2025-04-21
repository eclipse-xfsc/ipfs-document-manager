package main

import (
	"github.com/eclipse-xfsc/ipfs-document-manager/docs"
	vdr "github.com/eclipse-xfsc/ssi-vdr-core"
	"github.com/eclipse-xfsc/ssi-vdr-core/types"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var ipfs types.VerifiableDataRegistry

func init() {
	ipfs = vdr.VerifiableDataRegistryInitializer()
}

type Env interface {
	IsHealthy() bool
	Ipfs() types.VerifiableDataRegistry
	SetSwaggerBasePath(path string)
	SwaggerOptions() []func(config *ginSwagger.Config)
}

type EnvObj struct {
}

func (e *EnvObj) IsHealthy() bool {
	return e.Ipfs().IsAlive()
}

func (e *EnvObj) Ipfs() types.VerifiableDataRegistry {
	return ipfs
}

func (e *EnvObj) SetSwaggerBasePath(path string) {
	docs.SwaggerInfo.BasePath = path + BaseRoute
}
func (e *EnvObj) SwaggerOptions() []func(config *ginSwagger.Config) {
	return make([]func(config *ginSwagger.Config), 0)
}

func GetEnv() Env {
	return &EnvObj{}
}
