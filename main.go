package main

import (
    "card-event-service/handlers"
    "github.com/gofiber/fiber/v2"
)

func main() {
    app := fiber.New()

    app.Post("/cards", handlers.CreateCard)
    app.Get("/cards", handlers.GetCards)

    app.Listen(":3000")
}
