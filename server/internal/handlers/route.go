package handlers

import (
	"net/http"
	"time"

	"github.com/chrisS41/gobike-server/internal/database"
	"github.com/chrisS41/gobike-server/internal/errors"
	"github.com/chrisS41/gobike-server/internal/logger"
	"github.com/chrisS41/gobike-server/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RouteHandler struct {
	routes *database.Collection
	log    *logger.Log
}

func NewRouteHandler(routes *database.Collection, log *logger.Log) *RouteHandler {
	return &RouteHandler{routes: routes, log: log}
}

func (h *RouteHandler) CreateRoute(c *gin.Context) {
	var route models.Route
	if err := c.ShouldBindJSON(&route); err != nil {
		c.JSON(
			http.StatusBadRequest,
			models.NewErrorResponseWithMessage(errors.ErrInvalidRoute, err.Error()),
		)
		return
	}

	route.CreatedAt = time.Now()
	route.UpdatedAt = time.Now()

	var err error
	route.ID, err = h.routes.Create(route)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			models.NewErrorResponse(errors.ErrFailedToCreateRoute),
		)
		return
	}

	c.JSON(http.StatusCreated, route)
}

func (h *RouteHandler) GetUserRoutes(c *gin.Context) {
	userID, _ := primitive.ObjectIDFromHex(c.Param("userId"))

	cursor, err := h.routes.ReadMany(bson.M{
		"user_id": userID,
	})
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			models.NewErrorResponse(errors.ErrFailedToFetchRoutes),
		)
		return
	}

	var routes []models.Route
	for _, doc := range cursor {
		var route models.Route
		bsonBytes, _ := bson.Marshal(doc)
		bson.Unmarshal(bsonBytes, &route)
		routes = append(routes, route)
	}

	c.JSON(http.StatusOK, routes)
}

func (h *RouteHandler) GetRoute(c *gin.Context) {
	route, err := parseParams[models.Route](c, "id")
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			models.NewErrorResponseWithMessage(
				errors.ErrMissingParams,
				err.Error(),
			),
		)
		return
	}

	// TODO: 라우트 조회 로직 구현
	c.JSON(
		http.StatusOK,
		models.NewSuccessResponse(route.ID),
	)
}

func (h *RouteHandler) UpdateRoute(c *gin.Context) {
	// TODO: 라우트 업데이트 로직 구현
	c.JSON(
		http.StatusOK,
		models.NewSuccessResponse("route updated"),
	)
}

func (h *RouteHandler) DeleteRoute(c *gin.Context) {
	// TODO: 라우트 삭제 로직 구현
	c.JSON(
		http.StatusOK,
		models.NewSuccessResponse("route deleted"),
	)
}
