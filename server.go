package main

import (
	"github.com/bagus-aulia/pondok-sinau-golang/config"
	"github.com/bagus-aulia/pondok-sinau-golang/middleware"
	"github.com/bagus-aulia/pondok-sinau-golang/routes"
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
		api.GET("/auth/:provider", routes.RedirectHandler)          //done
		api.GET("/auth/:provider/callback", routes.CallbackHandler) //done

		api.POST("/check-username/:role", routes.CheckUsername) //done

		member := api.Group("/member")
		{
			member.GET("/", routes.MemberIndex)
			member.GET("/profile/:id", routes.MemberProfile)

			member.POST("/", middleware.IsAuth(), routes.MemberIndex)

			member.PUT("/update/:id", routes.MemberProfile)

			member.DELETE("/delete/:id", middleware.IsAdmin(), routes.MemberIndex)
		}

		admin := api.Group("/admin")
		{
			admin.GET("/", routes.AdminIndex)              //done
			admin.GET("/profile/:id", routes.AdminProfile) //minus show transactions too

			admin.PUT("/update/:id", middleware.IsAuth(), routes.AdminUpdate)           //done
			admin.PUT("/toggle-role/:id", middleware.IsAdmin(), routes.AdminToggleRole) //done
		}

		book := api.Group("/book")
		{
			book.GET("/", routes.BookIndex)              //done
			book.GET("/detail/:code", routes.BookDetail) //done

			book.POST("/", middleware.IsAuth(), routes.BookCreate) // minus upload cover

			book.PUT("/update/:id", middleware.IsAuth(), routes.BookUpdate) //minus update cover

			book.DELETE("/delete/:id", middleware.IsAuth(), routes.BookDelete) //done
		}

		borrow := api.Group("/borrow")
		{
			borrow.GET("/", middleware.IsAuth(), routes.BorrowIndex)              //minus detail info
			borrow.GET("/detail/:code", middleware.IsAuth(), routes.BorrowDetail) //minus detail info

			borrow.POST("/", middleware.IsAuth(), routes.BorrowCreate) //minus save detail

			borrow.PUT("/update/:id", middleware.IsAuth(), routes.BorrowUpdate) //minus update detail

			borrow.DELETE("/delete/:id", middleware.IsAuth(), routes.BorrowDelete) //minus delete detail
		}

		returns := api.Group("/return")
		{
			returns.GET("/get-borrow/:code", routes.ReturnGet) //belum dites

			returns.PUT("/update/:id", middleware.IsAuth(), routes.MemberIndex) //minus update detail

			returns.PUT("/rollback/:id", middleware.IsAuth(), routes.ReturnGet) //concept
		}
	}

	route.Run()
}
