package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v3"
)

type Config struct {
	Title string `json:"title"`
    Port int `json:"port"`
}

func main() {
    // Open and read the GlobalConfig file
    configFile, err := os.Open("GlobalConfig.json")
    if err != nil {
        fmt.Println("Error opening GlobalConfig file:", err)
        return
    }
    defer configFile.Close()

    // Decode JSON into Config struct
    var config Config
    decoder := json.NewDecoder(configFile)
    err = decoder.Decode(&config)
    if err != nil {
        fmt.Println("Error decoding config JSON:", err)
        return
    }

    // Create a new Fiber instance
    app := fiber.New()

    // Serve the home page with the form
    app.Get("/", func(c fiber.Ctx) error {
        // Render the head partial with the title from config
        return c.Render("./partial/head.html", fiber.Map{
            "Title": config.Title,
        })
    })

    // Start the Fiber app using the port from the config
    app.Listen(fmt.Sprintf(":%d", config.Port))
}