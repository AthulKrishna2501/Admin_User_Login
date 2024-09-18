package handlers

import (
	db "admin_user_login/DB"
	helpers "admin_user_login/Helpers"
	middleware "admin_user_login/Middleware"
	models "admin_user_login/Models"
	"fmt"
	"regexp"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *fiber.Ctx) error {

	if middleware.ValidateCookie(c) {

		return c.Redirect("/", fiber.StatusFound)
	}

	fmt.Println("Rendering Signup Page")

	return c.Render("signup", fiber.Map{
		"Error": nil,
	})
}

func SignupPost(c *fiber.Ctx) error {
	fmt.Println("Signup form submitted")

	var err models.InvalidErr
	username := c.FormValue("Name")
	email := c.FormValue("Email")
	password := c.FormValue("Password")
	confirmpassword := c.FormValue("ConfirmPassword")

	pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(pattern)
	if !regex.MatchString(email) {
		err.EmailError = "Email not in proper format"
		return c.Render("signup", fiber.Map{
			"EmailError":    err.EmailError,
			"PasswordError": err.PasswordError,
		})
	}

	if password != confirmpassword {
		err.PasswordError = "Passwords do not match"
		return c.Render("signup", fiber.Map{
			"EmailError":    err.EmailError,
			"PasswordError": err.PasswordError,
		})
	}

	var count int
	if dbErr := db.Db.Raw("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count).Error; dbErr != nil {
		fmt.Println("Error querying user count:", dbErr)
		return c.Render("signup", fiber.Map{
			"Error": "Error checking user existence",
		})
	}
	if count > 0 {
		err.EmailError = "User already exists"
		return c.Render("signup", fiber.Map{
			"EmailError": err.EmailError,
		})
	}

	hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if hashErr != nil {
		fmt.Println("Error hashing password:", hashErr)
		return c.Render("signup", fiber.Map{
			"Error": "Error hashing password",
		})
	}

	if insertErr := db.Db.Exec("INSERT INTO users(user_name, email, password) VALUES(?, ?, ?)", username, email, hashedPassword).Error; insertErr != nil {
		fmt.Println("Error inserting user:", insertErr)
		return c.Render("signup", fiber.Map{
			"Error": "Error inserting data",
		})
	}

	return c.Redirect("/", fiber.StatusFound)
}

func Login(c *fiber.Ctx) error {
	c.Set("Cache-Control", "no-store")
	c.Set("Pragma", "no-cache")
	c.Set("Expires", "0")

	if middleware.ValidateCookie(c) {
		role, _, _ := middleware.FindRole(c)
		if role == "user" {
			return c.Redirect("/home", fiber.StatusFound)
		} else if role == "admin" {
			return c.Redirect("/admin", fiber.StatusFound)
		}
	}

	return c.Render("login", fiber.Map{
		"EmailError":    nil,
		"PasswordError": nil,
	})
}

func LoginPost(c *fiber.Ctx) error {
	var err models.InvalidErr

	email := c.FormValue("Email")
	password := c.FormValue("Password")

	var compare models.Compare

	if dbErr := db.Db.Raw("SELECT password, role, user_name FROM users WHERE email = ?", email).Scan(&compare).Error; dbErr != nil {
		fmt.Println("Error querying user:", dbErr)
		return c.Render("login", fiber.Map{
			"Error": "Invalid email or password",
		})
	}

	if bcrypt.CompareHashAndPassword([]byte(compare.Password), []byte(password)) != nil {
		err.PasswordError = "Invalid Password"
		return c.Render("login", fiber.Map{
			"PasswordError": err.PasswordError,
		})
	}

	user := models.User{
		Role:     compare.Role,
		UserName: compare.UserName,
	}
	fmt.Println("Creating token for user:", user.UserName)

	tokenErr := helpers.CreateToken(user, c)
	if tokenErr != nil {
		fmt.Println("Error creating token:", tokenErr)
		return c.Render("login", fiber.Map{
			"Error": "Error creating token",
		})
	}

	if compare.Role == "user" {
		return c.Render("home", fiber.Map{
			"EmailError": nil,
		})
	} else if compare.Role == "admin" {
		return c.Redirect("/admin", fiber.StatusFound)
	}

	return c.Redirect("/home", fiber.StatusFound)
}

func Home(c *fiber.Ctx) error {

	c.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Set("Pragma", "no-cache")
	c.Set("Expires", "0")

	ok := middleware.ValidateCookie(c)
	role, user, _ := middleware.FindRole(c)
	if !ok || role != "user" {
		return c.Redirect("/", fiber.StatusFound)
	}

	return c.Render("home", fiber.Map{
		"UserName": user,
	})
}

func Logout(c *fiber.Ctx) error {

	middleware.DeleteCookie(c)

	c.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Set("Pragma", "no-cache")
	c.Set("Expires", "0")

	return c.Redirect("/", fiber.StatusFound)
}
