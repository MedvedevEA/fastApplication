package controller

import (
	"database/sql"
	"errors"
	"fastApplication/internal/model"
	"fmt"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service Service
}
type Service interface {
	ExecuteQuery(path string, req *model.Params) (*any, error)
	GetListQueries(schemaName string) ([]*model.Query, error)
	SetQueryRoutes(queryRoutes map[string]string)
}

func RegisterRoutes(router *gin.Engine, service Service, schemaName string) error {
	controller := &Controller{
		service: service,
	}

	listQuery, err := controller.service.GetListQueries(schemaName)
	if err != nil {
		return err
	}
	queryRoutes := map[string]string{}

	for _, value := range listQuery {
		var (
			method string
			path   string
		)

		switch value.BaseQueryName {
		case "add":
			method = "POST"
			path = fmt.Sprintf("/%s", value.TableName)
			router.POST(path, controller.add)
		case "get":
			method = "GET"
			path = fmt.Sprintf("/%s/:id", value.TableName)
			router.GET(path, controller.get)
		case "list":
			method = "GET"
			path = fmt.Sprintf("/%s", value.TableName)
			router.GET(path, controller.list)
		case "update":
			method = "PUT"
			path = fmt.Sprintf("/%s/:id", value.TableName)
			router.PUT(path, controller.update)
		case "remove":
			method = "DELETE"
			path = fmt.Sprintf("/%s/:id", value.TableName)
			router.DELETE(path, controller.remove)
		}
		queryRoutes[method+path] = value.Query

	}
	controller.service.SetQueryRoutes(queryRoutes)

	return nil
}

func (controller *Controller) add(ctx *gin.Context) {

	var body any
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.Status(400)
		return
	}
	path := ctx.Request.Method + ctx.FullPath()
	res, err := controller.service.ExecuteQuery(path, &model.Params{
		Body: &body,
	})
	if err != nil {
		ctx.Status(500)
		return
	}
	ctx.JSON(200, res)
}
func (controller *Controller) get(ctx *gin.Context) {
	var uriParam any = ctx.Param("id")
	path := ctx.Request.Method + ctx.FullPath()
	res, err := controller.service.ExecuteQuery(path, &model.Params{
		Uri: &uriParam,
	})
	if errors.Is(err, sql.ErrNoRows) {
		ctx.Status(404)
		return
	}
	if err != nil {
		ctx.Status(500)
		return
	}
	ctx.JSON(200, res)
}
func (controller *Controller) list(ctx *gin.Context) {
	var queryParam any
	if err := ctx.ShouldBindQuery(&queryParam); err != nil {
		ctx.Status(400)
		return
	}
	path := ctx.Request.Method + ctx.FullPath()
	res, err := controller.service.ExecuteQuery(path, &model.Params{
		Query: &queryParam,
	})
	if err != nil {
		ctx.Status(500)
		return
	}
	ctx.JSON(200, res)
}
func (controller *Controller) update(ctx *gin.Context) {
	var uriParam any = ctx.Param("id")
	var body any
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.Status(400)
		return
	}
	path := ctx.Request.Method + ctx.FullPath()
	res, err := controller.service.ExecuteQuery(path, &model.Params{
		Uri:  &uriParam,
		Body: &body,
	})
	if errors.Is(err, sql.ErrNoRows) {
		ctx.Status(404)
		return
	}
	if err != nil {
		ctx.Status(500)
		return
	}
	ctx.JSON(200, res)
}
func (controller *Controller) remove(ctx *gin.Context) {
	var uriParam any = ctx.Param("id")
	path := ctx.Request.Method + ctx.FullPath()
	_, err := controller.service.ExecuteQuery(path, &model.Params{
		Uri: &uriParam,
	})
	if errors.Is(err, sql.ErrNoRows) {
		ctx.Status(404)
		return
	}
	if err != nil {
		ctx.Status(500)
		return
	}
	ctx.Status(204)

}
