package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"order-stream/internal/entity"
	"order-stream/internal/usecase"
	"order-stream/pkg/logger"
)

type orderRoutes struct {
	uc usecase.Order
	l  logger.Interface
}

func newOrderRoutes(handler *gin.RouterGroup, uc usecase.Order, l logger.Interface) {
	r := &orderRoutes{uc, l}
	{
		handler.GET("/orders/:id", r.getOrderById)
	}
}

// @Summary     Get order by id
// @Description Get order by id
// @ID          getOrderById
// @Tags  	    order
// @Accept      json
// @Produce     json
// @Success     200 {object} entity.Order
// @Failure     500 {object} response
// @Router      /orders/:id [get]
func (r *orderRoutes) getOrderById(c *gin.Context) {
	id := c.Param("id")
	order, err := r.uc.GetOrder(c.Request.Context(), entity.OrderUID(id))
	if err != nil {
		r.l.Error(err, "http - v1 - getOrderById")
		errorResponse(c, http.StatusInternalServerError, "database problems")
		return
	}

	c.JSON(http.StatusOK, order)
}
