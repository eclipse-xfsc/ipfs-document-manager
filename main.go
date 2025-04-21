package main

import (
	"os"

	"github.com/eclipse-xfsc/microservice-core-go/pkg/logr"
	"github.com/eclipse-xfsc/microservice-core-go/pkg/server"
	"github.com/gin-gonic/gin"
)

const BaseRoute = "/api/ipfs"

var serviceConfiguration Config
var logger logr.Logger

var env Env

func init() {
	serviceConfiguration = getLoadedConfig()
	logger = createLogger(serviceConfiguration.LogLevel, serviceConfiguration.IsDev)
	env = GetEnv()
}

// @title			IPFS Document Manager API
// @version		1.0
// @description	A service wrapping Interplanetary Filesystem which can be used to store and retrieve files there.
// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
// @host			localhost:8000
func main() {
	engine := server.New(env, serviceConfiguration.ServerMode)
	addEndpoints(engine)
	err := engine.Run(serviceConfiguration.ListenPort)
	if err != nil {
		logger.Error(err, "failed to run server")
		os.Exit(1)
	}
}

func addEndpoints(engine *server.GinServer) {
	engine.Add(func(base *gin.RouterGroup) {
		base = base.Group(BaseRoute)
		base.GET("/:id", WrapHandler(GetDocument, env))
		base.GET("/list", GetDocuments(env))
		base.POST("/create", WrapHandler(CreateDocument, env))
		base.PUT(":id/update", WrapHandler(UpdateDocument, env))
		base.DELETE("/:id", WrapHandler(DeleteDocument, env))
	})
}
