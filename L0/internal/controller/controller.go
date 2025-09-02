package controller

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/ChursinAlexUnder/wbtech-golang-course/L0/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Controller struct {
	pool  *pgxpool.Pool
	cache *expirable.LRU[uuid.UUID, database.Orders]
}

func NewController(pool *pgxpool.Pool, cache *expirable.LRU[uuid.UUID, database.Orders]) *Controller {
	return &Controller{pool: pool, cache: cache}
}

type Message struct {
	Message string `json:"message" example:"message"`
}

type HTTPError struct {
	Code    int    `json:"code" example:"500"`
	Message string `json:"message" example:"bad request"`
}

// NewHTTPError — отправляет ошибку клиенту и завершает обработку
func NewHTTPError(ctx *gin.Context, status int, err error) {
	ctx.AbortWithStatusJSON(status, HTTPError{
		Code:    status,
		Message: err.Error(),
	})
}

// @Summary      Главная страница сервиса
// @Description  Показ главной страницы сервиса, где можно вставить order_uid для получения информации интересующего заказа
// @Tags         order
// @Produce      html
// @Success      200  {string}  string "HTML, css и js"
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /order [get]
func (c *Controller) GetMainPage(ctx *gin.Context) {
	ctx.File(filepath.Join("web", "index.html"))
}

// @Summary      Отправка информации о заказе
// @Description  Отправка json файла с сервера на клиент
// @Tags         api
// @Accept       json
// @Produce      json
// @Param        order_uid   path      string  true  "Uid заказа"
// @Success      200  {object}  database.Orders
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /api/{order_uid} [get]
func (c *Controller) GetOrderByUid(ctx *gin.Context) {
	var (
		answer database.Orders
		err    error
	)
	// Обрабатываем order_uid
	order_uid_string := ctx.Param("order_uid")
	order_uid, err := uuid.Parse(order_uid_string)
	if err != nil {
		NewHTTPError(ctx, http.StatusBadRequest, err)
		return
	}

	// Сначала ищем в кеше, иначе берем из бд
	answer, ok := c.cache.Get(order_uid)
	if !ok {
		answer, err = database.SelectOrderByUid(ctx, c.pool, order_uid_string)
		if err != nil {
			NewHTTPError(ctx, http.StatusBadRequest, err)
			return
		}
		log.Printf("Запись с order_uid %s взята из бд!\n", answer.Order_uid)
	} else {
		log.Printf("Запись с order_uid %s успешно взята из кеша!\n", answer.Order_uid)
	}
	ctx.JSON(http.StatusOK, answer)
}

// @Summary      Перенаправление
// @Description  Перенаправление на главную страницу сервиса
// @Tags         order
// @Produce      html
// @Success      200  {string}  string "HTML, css и js"
// @Success      301  {string}  string "HTML, css и js"
// @Success      304  {string}  string "HTML, css и js"
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       / [get]
func (c *Controller) RedirectOnMainPage(ctx *gin.Context) {
	ctx.Redirect(http.StatusMovedPermanently, "/order")
}
