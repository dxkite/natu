package server

import (
	"net/http"

	"dxkite.cn/meownest/pkg/httpserver"
	"dxkite.cn/meownest/src/constant"
	"dxkite.cn/meownest/src/service"
	"github.com/gin-gonic/gin"
)

func NewRoute(s service.Route) *Route {
	return &Route{s: s}
}

type Route struct {
	s service.Route
}

// 创建路由
//
// @Summary      创建路由
// @Description  创建路由
// @Tags         路由
// @Accept       json
// @Produce      json
// @Param        body body service.CreateRouteParam true "请求体"
// @Success      201  {object} dto.Route
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /routes [post]
func (s *Route) Create(c *gin.Context) {
	var param service.CreateRouteParam

	if err := c.ShouldBind(&param); err != nil {
		httpserver.ResultErrorBind(c, err)
		return
	}

	rst, err := s.s.Create(c, &param)
	if err != nil {
		httpserver.ResultError(c, err)
		return
	}

	httpserver.Result(c, http.StatusCreated, rst)
}

// 获取路由
//
// @Summary      获取路由
// @Description  获取路由
// @Tags         路由
// @Accept       json
// @Produce      json
// @Param        id path string true "路由ID"
// @Param        expand query []string false "展开数据"
// @Success      200  {object} dto.Route
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /routes/{id} [get]
func (s *Route) Get(c *gin.Context) {
	var param service.GetRouteParam

	if err := c.ShouldBindUri(&param); err != nil {
		httpserver.ResultErrorBind(c, err)
		return
	}

	if err := c.ShouldBindQuery(&param); err != nil {
		httpserver.ResultErrorBind(c, err)
		return
	}

	rst, err := s.s.Get(c, &param)
	if err != nil {
		httpserver.ResultError(c, err)
		return
	}
	httpserver.Result(c, http.StatusOK, rst)
}

// 路由列表
//
// @Summary      路由列表
// @Description  路由列表
// @Tags         路由
// @Accept       json
// @Produce      json
// @Param        name query string false "搜索名称"
// @Param        path query string false "搜索路径"
// @Param		 include_total query bool false "是否包含total"
// @Param        page query int false "页码"
// @Param        pre_page query int false "每页数量"
// @Param        expand query []string false "展开数据"
// @Success      200  {object} service.ListRouteResult
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /routes [get]
func (s *Route) List(c *gin.Context) {
	var param service.ListRouteParam

	if err := c.ShouldBindQuery(&param); err != nil {
		httpserver.ResultErrorBind(c, err)
		return
	}

	rst, err := s.s.List(c, &param)
	if err != nil {
		httpserver.ResultError(c, err)
		return
	}

	httpserver.Result(c, http.StatusOK, rst)
}

// 更新路由
//
// @Summary      更新路由
// @Description  更新路由
// @Tags         路由
// @Accept       json
// @Produce      json
// @Param        id path string true "路由ID"
// @Param        body body service.UpdateRouteParam true "数据"
// @Success      200  {object} dto.Route
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /routes/{id} [post]
func (s *Route) Update(c *gin.Context) {
	var param service.UpdateRouteParam
	param.Id = c.Param("id")

	if err := c.ShouldBind(&param); err != nil {
		httpserver.ResultErrorBind(c, err)
		return
	}

	rst, err := s.s.Update(c, &param)
	if err != nil {
		httpserver.ResultError(c, err)
		return
	}
	httpserver.Result(c, http.StatusOK, rst)
}

// 删除路由
//
// @Summary      删除路由
// @Description  删除路由
// @Tags         路由
// @Accept       json
// @Produce      json
// @Param        id path string true "路由ID"
// @Success      200
// @Failure      400  {object} httpserver.HttpError
// @Failure      500  {object} httpserver.HttpError
// @Router       /routes/{id} [delete]
func (s *Route) Delete(c *gin.Context) {
	var param service.DeleteRouteParam

	if err := c.ShouldBindUri(&param); err != nil {
		httpserver.ResultErrorBind(c, err)
		return
	}
	err := s.s.Delete(c, &param)
	if err != nil {
		httpserver.ResultError(c, err)
		return
	}

	httpserver.ResultEmpty(c, http.StatusOK)
}

func (s *Route) API() httpserver.RouteHandleFunc {
	return func(route gin.IRouter) {
		route.GET("/routes", httpserver.ScopeRequired(constant.ScopeRouteRead), s.List)
		route.POST("/routes", httpserver.ScopeRequired(constant.ScopeRouteWrite), s.Create)
		route.GET("/routes/:id", httpserver.ScopeRequired(constant.ScopeRouteRead), s.Get)
		route.DELETE("/routes/:id", httpserver.ScopeRequired(constant.ScopeRouteWrite), s.Delete)
		route.POST("/routes/:id", httpserver.ScopeRequired(constant.ScopeRouteWrite), s.Update)
	}
}
