package handlers

import (
	db "admin_user_login/DB"
	middleware "admin_user_login/Middleware"
	models "admin_user_login/Models"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type AdminResponse struct {
	Name    string
	Users   []models.UserDetails
	Invalid models.InvalidErr
}

type AdminSearch struct {
	UserList    []models.User
	SearchError string
}

var errors models.InvalidErr

func AdminHome(c *fiber.Ctx) error {
	c.Set("Cache-Control", "no-cache, no-store")
	c.Set("Expires", "0")

	ok := middleware.ValidateCookie(c)
	if !ok {
		return c.Render("login", fiber.Map{
			"EmailError":    nil,
			"PasswordError": nil,
		})
	}

	role, name, err := middleware.FindRole(c)
	if err != nil {
		log.Fatal(err)
	}
	if role != "admin" {
		return c.Redirect("/", fiber.StatusFound)
	}

	var Collect []models.UserDetails

	if err := db.Db.Raw("SELECT user_name,email from users").Scan(&Collect).Error; err != nil {
		log.Fatal("Error in fecthing users", err)
	}
	result := AdminResponse{
		Name:    name,
		Users:   Collect,
		Invalid: errors,
	}
	fmt.Println("Rendering admin")
	err = c.Render("admin", fiber.Map{
		"title": result,
	})
	if err != nil {
		log.Fatal(err)
	}
	return nil

}

func AdminAddUser(c *fiber.Ctx) error {
	c.Set("Cache-Control", "no-cache, no-store")
	c.Set("Pragma", "no-cache")
	c.Set("Expires", "0")
	return c.Render("admin", fiber.Map{
		"error": nil,
	})
}

func AdminAddUserPost(c *fiber.Ctx) error {
	c.Set("Cache-Control", "no-cache, no-store")
	c.Set("Pragma", "no-cache")
	c.Set("Expires", "0")

	ok := middleware.ValidateCookie(c)
	role, _, _ := middleware.FindRole(c)
	if !ok || role != "admin" {
		return c.Redirect("/", fiber.StatusOK)
	}

	userName := c.FormValue("Name")
	userEmail := c.FormValue("Email")
	userPassword := c.FormValue("Password")

	
	var count int
	if err := db.Db.Raw("SELECT COUNT(*) FROM users WHERE email = $1", userEmail).Scan(&count).Error; err != nil {
		log.Fatal(err)
		return c.Redirect("/admin", fiber.StatusFound)
	}

	
	if count > 0 {
		errors.EmailError = "User already exists"
		return c.Render("admin", fiber.Map{
			"title": AdminResponse{
				Name:    "",
				Users:   nil,
				Invalid: errors,
			},
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
		return c.Redirect("/adminAddUser", fiber.StatusFound)
	}

	if err := db.Db.Exec("INSERT INTO users (user_name, email, password) VALUES ($1, $2, $3)", userName, userEmail, string(hashedPassword)).Error; err != nil {
		log.Fatal(err)
		return c.Redirect("/adminAddUser", fiber.StatusFound)
	}

	return c.Redirect("/admin", fiber.StatusFound)
}

func AdminUpdate(c *fiber.Ctx) error {
	c.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Set("Pragma", "no-cache")
	c.Set("Expires", "0")
	ok := middleware.ValidateCookie(c)
	if !ok {
		return c.Redirect("/", fiber.StatusOK)
	}
	username := c.Query("Username")
	email := c.Query("Email")
	log.Println("Received Username:", username)
	log.Println("Received Email:", email)
	return c.Render("updateuser", fiber.Map{
		"UserName": username,
		"Email":    email,
	})

}

func AdminUpdatePost(c *fiber.Ctx) error {
	c.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Set("Pragma", "no-cache")
	c.Set("Expires", "0")

	ok := middleware.ValidateCookie(c)
	if !ok {
		return c.Redirect("/", fiber.StatusOK)

	}
	fmt.Println("HIII")
	email := c.Query("Email")
	userName := c.FormValue("Name")
	log.Println("Received email:", email, "and new username:", userName)
	if err := db.Db.Exec("UPDATE users SET user_name = $1 WHERE email = $2", userName, email).Error; err != nil {
		log.Fatal(err)

	}
	return c.Redirect("/admin", fiber.StatusFound)
}

func AdminDelete(c *fiber.Ctx) error {
	c.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Set("Pragma", "no-cache")
	c.Set("Expires", "0")

	ok := middleware.ValidateCookie(c)
	role, _, _ := middleware.FindRole(c)
	if !ok || role != "admin" {
		return c.Redirect("/", fiber.StatusOK)
	}
	email := c.Query("Email")

	if err := db.Db.Exec("DELETE FROM users WHERE email = $1", email).Error; err != nil {
		log.Fatal("Could not fetch details", err)
	}
	return c.Redirect("/admin", fiber.StatusFound)
}

func AdminLogout(c *fiber.Ctx) error {
	middleware.DeleteCookie(c)
	return c.Redirect("/", fiber.StatusFound)
}
