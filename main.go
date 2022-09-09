package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/gofiber/template/html"
	_ "github.com/lib/pq"
)

func main() {

	connStr := "postgresql://postgres:test123@localhost:5432/postgres?sslmode=disable" // Conecta ao banco de dados

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/swagger/*", swagger.HandlerDefault) // default

	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL:          "http://example.com/doc.json",
		DeepLinking:  false,
		DocExpansion: "none",
		OAuth: &swagger.OAuthConfig{
			AppName:  "OAuth Provider",
			ClientId: "21bb4edc-05a7-4afc-86f1-2e151e4ba6e2",
		},
		OAuth2RedirectUrl: "http://localhost:8080/swagger/oauth2-redirect.html",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return indexHandler(c, db)
	})

	app.Post("/", func(c *fiber.Ctx) error {
		return postHandler(c, db)
	})

	app.Put("/update", func(c *fiber.Ctx) error {
		return putHandler(c, db)
	})

	app.Delete("/delete", func(c *fiber.Ctx) error {
		return deleteHandler(c, db)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))

	app.Static("/", "./public")
	log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))

}

func indexHandler(c *fiber.Ctx, db *sql.DB) error {
	var res string
	var sm_users []string
	rows, err := db.Query("SELECT * FROM sm_users")

	if err != nil {
		log.Fatalln(err)
		c.JSON("An error occured")
	}
	for rows.Next() {
		rows.Scan(&res)
		sm_users = append(sm_users, res)
	}
	return c.Render("index", fiber.Map{
		"Sm_users": sm_users,
	})
}

type sm_users struct {
	Username string
}

func postHandler(c *fiber.Ctx, db *sql.DB) error {
	newSm_users := sm_users{}
	if err := c.BodyParser(&newSm_users); err != nil {
		log.Printf("An error occured: %v", err)

		return c.SendString(err.Error())
	}
	fmt.Printf("%v", newSm_users)
	if newSm_users.Username != "" {
		_, err := db.Exec("INSERT into sm_users VALUES ($1)", newSm_users.Username)
		if err != nil {
			log.Fatalf("An error occured while executing query: %v", err)
		}
	}

	return c.Redirect("/")
}

func putHandler(c *fiber.Ctx, db *sql.DB) error {
	oldusername := c.Query("oldusername")
	newusername := c.Query("newusername")
	db.Exec("UPDATE sm_users SET username=$1 WHERE username=$2", newusername, oldusername)
	return c.Redirect("/")
}

func deleteHandler(c *fiber.Ctx, db *sql.DB) error {
	sm_usersToDelete := c.Query("username")
	db.Exec("DELETE from sm_users WHERE username=$1", sm_usersToDelete)
	return c.SendString("deleted")
}
