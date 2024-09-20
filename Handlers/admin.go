package handlers

import (
	db "admin_user_login/DB"
	middleware "admin_user_login/Middleware"
	models "admin_user_login/Models"
	"fmt"
	"log"
	"regexp"

	"github.com/gofiber/fiber/v2"
)

type AdminResponse struct {
	Name    string
	Users   []models.UserDetails
	Invalid models.InvalidErr
}

var errors models.InvalidErr

func AdminHome(c *fiber.Ctx) error {
	c.Set("Cache-Control", "no-cache,no-store,must-revalidate")
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
	c.Set("Cache-Control", "no-cache, no-store, must-revalidate")
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

	pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(pattern)
	if !regex.MatchString(userEmail) {
		errors.EmailError = "Email not in the correct format"
		return c.Redirect("/admin", fiber.StatusFound)
	}

	var count int
	if err := db.Db.Raw("SELECT COUNT(*) FROM users WHERE email=$1", userEmail).Scan(&count).Error; err != nil {
		log.Fatal(err)
		return c.Redirect("/admin", fiber.StatusFound)
	}
	if count > 0 {
		errors.Err = "User already exists"
		return c.Redirect("/admin", fiber.StatusFound)
	}

	var userRole string
	if c.FormValue("checkbox") == "on" {
		userRole = "admin"
	} else {
		userRole = "user"
	}

	if err := db.Db.Exec("INSERT INTO users (user_name, email, password, user_role) VALUES($1, $2, $3, $4)", userName, userEmail, userPassword, userRole).Error; err != nil {
		log.Fatal(err)
		return c.Redirect("/admin", fiber.StatusFound)
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
	email := c.Query("Email")
	userName := c.FormValue("Name")
	if err := db.Db.Exec("UPDATE users SET user_name = $1 where email = $2", userName, email).Error; err != nil {
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
