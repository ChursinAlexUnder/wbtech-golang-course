package router

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/ChursinAlexUnder/wbtech-golang-course/L0/go-app/database"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRouter(pool *pgxpool.Pool, ctx context.Context) *gin.Engine {
	// Режим прода и тестирования соответственно (доп. логирование)
	gin.SetMode(gin.ReleaseMode)
	// gin.SetMode(gin.DebugMode)

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// Редирект на страницу с web-интерфейсом
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/order")
	})

	// Отдача основного html файла страницы
	router.GET("/order/*any", func(c *gin.Context) {
		c.File(filepath.Join("frontend", "index.html"))
	})

	// Отдача статики (css, js)
	router.Static("/static", "frontend")

	// Эндпоинт для получения данных по order_uid
	router.GET("/api/:order_uid", func(c *gin.Context) {
		order_uid := c.Param("order_uid")
		answer, err := database.GetOrderByUid(pool, ctx, order_uid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fmt.Println(answer)

	})
	{

		// api.GET("/universities", func(c *gin.Context) {

		// 	currentRow, err := strconv.Atoi(c.DefaultQuery("currentrow", "0"))
		// 	if err != nil || currentRow < 0 {
		// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 		return
		// 	}
		// 	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
		// 	if err != nil || limit < 1 {
		// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 		return
		// 	}
		// 	query := c.DefaultQuery("search", "")
		// 	location := c.DefaultQuery("location", "")
		// 	min, err := strconv.Atoi(c.DefaultQuery("min", "0"))
		// 	if err != nil || min < 0 {
		// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 		return
		// 	}
		// 	max, err := strconv.Atoi(c.DefaultQuery("max", "0"))
		// 	if err != nil || max < 0 {
		// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 		return
		// 	}
		// 	tmp, err := strconv.ParseFloat(c.DefaultQuery("avg", "0"), 32)
		// 	if err != nil || tmp < 0 {
		// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 		return
		// 	}
		// 	avg := float32(tmp)
		// 	intl, err := strconv.Atoi(c.DefaultQuery("intl", "0"))
		// 	if err != nil || intl < 0 {
		// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 		return
		// 	}
		// 	female, err := strconv.Atoi(c.DefaultQuery("female", "0"))
		// 	if err != nil || female < 0 {
		// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 		return
		// 	}
		// 	male, err := strconv.Atoi(c.DefaultQuery("male", "0"))
		// 	if err != nil || male < 0 {
		// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 		return
		// 	}

		// 	list, hasMore, err := database.GetUniversities(pool, currentRow, limit, min, max, intl, female, male, query, location, avg)
		// 	if err != nil {
		// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 		return
		// 	}
		// 	c.JSON(http.StatusOK, gin.H{
		// 		"data":    list,
		// 		"hasMore": hasMore,
		// 	})
		// })

		// api.GET("/locations", func(c *gin.Context) {
		// 	list, err := database.GetLocations(pool)
		// 	if err != nil {
		// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 		return
		// 	}
		// 	c.JSON(http.StatusOK, gin.H{"countries": list})
		// })

		// api.GET("/studentscount", func(c *gin.Context) {
		// 	min, max, err := database.GetStudentsCount(pool)
		// 	if err != nil {
		// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 		return
		// 	}
		// 	c.JSON(http.StatusOK, gin.H{"min": min, "max": max})
		// })

		// api.GET("/studentsaverage", func(c *gin.Context) {
		// 	average, err := database.GetStudentsAverage(pool)
		// 	if err != nil {
		// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 		return
		// 	}
		// 	c.JSON(http.StatusOK, gin.H{"average": average})
		// })
	}

	return router
}
