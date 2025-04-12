package routes

import(
	"webroutes/controllers"
	"webroutes/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	authController := controllers.NewAuthController(db)
	userController := controllers.NewUserController(db)
	profileController := controllers.NewProfileController(db)
	postController := controllers.NewPostController(db)


	api := r.Group("/api") 
	{
		auth := api.Group("/auth") 
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}

		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware()) 
		{
			protected.GET("/users", userController.GetUsers)	
			protected.GET("/get-without-db", controllers.GetUserWithoutDB)	
			protected.POST("/create-without-db", controllers.CreateUserWithoutDB)	
			protected.POST("/create-with-db", userController.CreateUser)	
			protected.POST("/create-profile", profileController.CreateProfile)	

			protected.POST("/tag", postController.CreateTag)	
			protected.POST("/post", postController.CreatePost)	
			protected.PUT("/post/:id", postController.UpdatePost)	
			protected.DELETE("/post/:id", postController.DeletePost)	
		}
	}
}