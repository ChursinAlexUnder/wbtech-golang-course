package router

import (
	"context"

	"github.com/ChursinAlexUnder/wbtech-golang-course/L0/go-app/controller"
	_ "github.com/ChursinAlexUnder/wbtech-golang-course/L0/go-app/docs"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(ctx context.Context, pool *pgxpool.Pool) *gin.Engine {
	// Режим прода и тестирования соответственно (доп. логирование)
	gin.SetMode(gin.ReleaseMode)
	// gin.SetMode(gin.DebugMode)

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	c := controller.NewController(pool)

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
