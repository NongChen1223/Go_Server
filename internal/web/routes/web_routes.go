package routers

//
//import (
//	"github.com/gin-gonic/gin"
//	"go_server/internal/web/controllers"
//)
//
//// SetupWebRoutes 设置官网端路由
//func SetupWebRoutes(router *gin.RouterGroup) {
//	// 官网路由组
//	web := router.Group("/web")
//	{
//		// 首页相关
//		web.GET("/home", controllers.GetHomeData)           // 获取首页数据
//		web.GET("/featured", controllers.GetFeaturedGames)  // 获取推荐游戏
//
//		// 游戏展示
//		gameGroup := web.Group("/games")
//		{
//			gameGroup.GET("", controllers.GetWebGames)              // 获取游戏列表（官网展示）
//			gameGroup.GET("/:id", controllers.GetWebGameDetail)     // 获取游戏详情（官网展示）
//			gameGroup.GET("/categories", controllers.GetGameCategories) // 获取游戏分类
//			gameGroup.GET("/hot", controllers.GetHotGames)          // 获取热门游戏
//			gameGroup.GET("/new", controllers.GetNewGames)          // 获取最新游戏
//		}
//
//		// 新闻资讯
//		newsGroup := web.Group("/news")
//		{
//			newsGroup.GET("", controllers.GetNewsList)         // 获取新闻列表
//			newsGroup.GET("/:id", controllers.GetNewsDetail)   // 获取新闻详情
//			newsGroup.GET("/hot", controllers.GetHotNews)      // 获取热门新闻
//		}
//
//		// 厂商信息
//		publisherGroup := web.Group("/publishers")
//		{
//			publisherGroup.GET("", controllers.GetPublishers)           // 获取厂商列表
//			publisherGroup.GET("/:id", controllers.GetPublisherDetail)  // 获取厂商详情
//			publisherGroup.GET("/:id/games", controllers.GetPublisherGames) // 获取厂商游戏
//		}
//
//		// 搜索功能
//		searchGroup := web.Group("/search")
//		{
//			searchGroup.GET("/games", controllers.SearchGames)      // 搜索游戏
//			searchGroup.GET("/suggest", controllers.GetSearchSuggest) // 搜索建议
//		}
//
//		// 统计信息（公开）
//		statsGroup := web.Group("/stats")
//		{
//			statsGroup.GET("/overview", controllers.GetPublicStats) // 获取公开统计信息
//		}
//	}
//}
