package rest

import "github.com/gin-gonic/gin"

const (
	ROOT    = ""
	API     = "api"
	VERSION = "v1"
	DEVICES = "devices"
)

func (r *Rest) Router() *gin.Engine {

	router := gin.Default()

	// Limit request body size to ~6MB
	router.MaxMultipartMemory = 6 << 20

	r.htmlAdjust(router)

	root := router.Group(ROOT)
	r.rootGroup(root)

	api := router.Group(API)
	{
		api.Use(r.MiddlewareHandle())
		r.versionGroup(api, VERSION)
	}

	return router
}

func (r *Rest) htmlAdjust(router *gin.Engine) {
	router.LoadHTMLGlob("./internal/network/web/templates/*.html")
	router.Static("/static", "./internal/network/web/templates/static")
}
