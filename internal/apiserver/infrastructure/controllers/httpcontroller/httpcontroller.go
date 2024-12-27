package httpcontroller

import (
	"github.com/gin-gonic/gin"
)

type HttpController struct {
	router *gin.Engine
}

func Init(router *gin.Engine) {

	httpControllers := &HttpController{
		router: router,
	}
	router.POST("add", httpControllers.add)
	router.GET("get", httpControllers.get)
	router.GET("list", httpControllers.list)
	router.PUT("update", httpControllers.update)
	router.DELETE("remove", httpControllers.remove)

}

func (httpControllers *HttpController) add(ctx *gin.Context) {

}
func (httpControllers *HttpController) get(ctx *gin.Context) {

}
func (httpControllers *HttpController) list(ctx *gin.Context) {

}
func (httpControllers *HttpController) update(ctx *gin.Context) {

}

func (httpControllers *HttpController) remove(ctx *gin.Context) {

}
