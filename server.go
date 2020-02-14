package main

import (
	"github.com/bagus-aulia/pondok-lentera/config"
	"github.com/bagus-aulia/pondok-lentera/middleware"
	"github.com/bagus-aulia/pondok-lentera/routes"
	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
)

func main() {
	gotenv.Load()
	config.InitDB()
	defer config.DB.Close()

	route := gin.Default()

	api := route.Group("/api/v1/")
	{
		api.GET("/", routes.Home)
		api.GET("/dashboard", routes.Home)
		api.GET("/auth/:provider", routes.RedirectHandler)
		api.GET("/auth/:provider/callback", routes.CallbackHandler)

		member := api.Group("/member")
		{
			member.GET("/", routes.MemberIndex)
			member.GET("/profile/:id", routes.MemberProfile)

			member.PUT("/update/:id", routes.MemberProfile)
		}

		admin := api.Group("/admin")
		{
			admin.GET("/", routes.MemberIndex)
			admin.GET("/profile/:id", routes.MemberIndex)

			admin.PUT("/update/:id", routes.MemberIndex)
		}

		book := api.Group("/book")
		{
			book.GET("/", routes.MemberIndex)
			book.GET("/detail/:code", routes.MemberIndex)

			book.POST("/", middleware.IsAuth(), routes.MemberIndex)

			book.PUT("/update/:id", routes.MemberIndex)

			book.DELETE("/delete/:id", routes.MemberIndex)
		}

		borrow := api.Group("/borrow")
		{
			borrow.GET("/", routes.MemberIndex)
			borrow.GET("/expired", routes.MemberIndex)
			borrow.GET("/create", routes.MemberIndex)
			borrow.GET("/add-book/:code", routes.MemberIndex)
			borrow.GET("/detail/:code", routes.MemberIndex)

			borrow.POST("/", routes.MemberIndex)

			borrow.PUT("/update/:id", routes.MemberIndex)
		}

		returns := api.Group("/return")
		{
			returns.GET("/create", routes.MemberIndex)
			returns.GET("/get-borrow/:code", routes.MemberIndex)

			returns.POST("/", routes.MemberIndex)

			returns.PUT("/update/:id", routes.MemberIndex)
		}
	}

	route.Run()
}
