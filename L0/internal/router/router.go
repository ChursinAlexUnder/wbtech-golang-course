package router

import (
	"context"

	_ "github.com/ChursinAlexUnder/wbtech-golang-course/L0/docs"
	"github.com/ChursinAlexUnder/wbtech-golang-course/L0/internal/controller"
	"github.com/ChursinAlexUnder/wbtech-golang-course/L0/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(ctx context.Context, pool *pgxpool.Pool, cache *expirable.LRU[uuid.UUID, database.Orders]) *gin.Engine {
	// Режим прода и тестирования соответственно (доп. логирование)
	gin.SetMode(gin.ReleaseMode)
	// gin.SetMode(gin.DebugMode)

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	c := controller.NewController(pool, cache)

	// Редирект на страницу с web-интерфейсом
	router.GET("/", c.RedirectOnMainPage)

	// /swagger/index.html - страница swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Главная страница сервиса
	router.GET("/order/*any", c.GetMainPage)

	// Отдача статики (css, js)
	router.Static("/static", "frontend")

	// Служебный эндпоинт для передачи json о заказе
	router.GET("/api/:order_uid", c.GetOrderByUid)

	return router
}
