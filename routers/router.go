package routers

import (
	"mixfrag/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api", controllers.BasicAuthMiddleware)
	{
		api.GET("/profile", controllers.Profile)

		category_note := api.Group("/categorie_notes")
		{
			category_note.GET("", controllers.GetCategoriesNote)
			category_note.POST("", controllers.CreateCategoryNote)
			category_note.GET("/:id", controllers.GetCategoryNote)
			category_note.DELETE("/:id", controllers.DeleteCategoryNote)
		}

		category_parfume := api.Group("/categorie_parfumes")
		{
			category_parfume.GET("", controllers.GetCategorieParfumes)
			category_parfume.POST("", controllers.CreateCategoryParfume)
			category_parfume.GET("/:id", controllers.GetCategoryParfume)
			category_parfume.DELETE("/:id", controllers.DeleteCategoryParfume)
		}

		category_smell := api.Group("/categorie_smells")
		{
			category_smell.GET("", controllers.GetCategorieSmells)
			category_smell.POST("", controllers.CreateCategorySmell)
			category_smell.GET("/:id", controllers.GetCategorySmell)
			category_smell.DELETE("/:id", controllers.DeleteCategorySmell)
		}

		essential_oil := api.Group("/essential_oils")
		{
			essential_oil.GET("", controllers.GetEssentialOils)
			essential_oil.POST("", controllers.CreateEssentialOil)
			essential_oil.GET("/:id", controllers.GetEssentialOil)
			essential_oil.DELETE("/:id", controllers.DeleteEssentialOil)
		}

		parfume := api.Group("/parfumes")
		{
			parfume.GET("", controllers.GetParfumes)
			parfume.POST("", controllers.CreateParfume)
			parfume.GET("/:id", controllers.GetParfumeByID)
			parfume.DELETE("/:id", controllers.DeleteParfume)
		}
	}

	return r
}
