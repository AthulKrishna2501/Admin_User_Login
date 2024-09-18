package main

import (
	db "admin_user_login/DB"
	handlers "admin_user_login/Handlers"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	db.InitDatabase()

	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	//User
	app.Get("/", handlers.Login)
	app.Post("/", handlers.LoginPost)
	app.Get("/signup", handlers.Signup)
	app.Post("/signup", handlers.SignupPost)
	app.Get("/home", handlers.Home)
	app.Get("/logout", handlers.Logout)

	// //Admin
	app.Get("/admin", handlers.AdminHome)
	app.Post("/admin", handlers.AdminAddUser)
	app.Get("adminupdate", handlers.AdminUpdate)
	app.Post("adminupdatepost", AdminUpdatePost)
	app.Get("/admindelete", AdminDelete)
	app.Get("/adminlogou", handlers.AdminLogout)
	fmt.Println("Server started at :8080")
	app.Listen(":8080")
}
