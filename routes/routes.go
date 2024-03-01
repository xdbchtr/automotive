package routes

import (
	"automotive/controllers"
	"automotive/middlewares"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Initialize the validator instance
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
	})

	r.POST("/login", controllers.Login)

	r.GET("/roles", controllers.GetAllRoles)
	r.GET("/roles-with-users", controllers.GetAllRolesWithUsers)
	roleRoute := r.Group("/role")
	roleRoute.Use(middlewares.JWTAuthMiddleware())
	roleRoute.POST("/", controllers.CreateRole)
	roleRoute.PUT("/:id", controllers.UpdateRole)
	roleRoute.DELETE("/:id", controllers.DeleteRole)

	r.GET("/users", controllers.GetAllUser)
	userRoute := r.Group("/user")
	userRoute.Use(middlewares.JWTAuthMiddleware())
	userRoute.POST("/", controllers.CreateUser)
	userRoute.PUT("/:id", controllers.UpdateUser)
	userRoute.DELETE("/:id", controllers.DeleteUser)

	r.GET("/brands", controllers.GetAllBrand)
	brandRoute := r.Group("/brand")
	brandRoute.Use(middlewares.JWTAuthMiddleware())
	brandRoute.POST("/", controllers.CreateBrand)
	brandRoute.PUT("/:id", controllers.UpdateBrand)

	r.GET("/cars", controllers.GetAllCars)
	carRoute := r.Group("/car")
	carRoute.Use(middlewares.JWTAuthMiddleware())
	carRoute.POST("/", controllers.CreateCar)
	carRoute.PUT("/:id", controllers.UpdateCar)
	carRoute.DELETE("/:id", controllers.DeleteCar)

	r.GET("/types", controllers.GetAllTypes)
	typeRoute := r.Group("/type")
	typeRoute.Use(middlewares.JWTAuthMiddleware())
	typeRoute.POST("/", controllers.CreateType)
	typeRoute.PUT("/:id", controllers.UpdateType)
	typeRoute.DELETE("/:id", controllers.DeleteType)

	cartRoute := r.Group("/cart")
	cartRoute.Use(middlewares.JWTAuthMiddleware())
	cartRoute.PUT("/update-cart", controllers.UpsertCart)
	cartRoute.GET("/user-cart", controllers.GetUserCart)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
