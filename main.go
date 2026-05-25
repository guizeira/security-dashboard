package main

import (
	"log"
	"time"

	"intel-dashboard/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New(fiber.Config{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 5 * time.Minute,
		IdleTimeout:  30 * time.Second,
	})

	// ======================================================
	// STATIC FILES
	// ======================================================

	app.Static("/static", "./static")

	// ======================================================
	// ROUTES
	// ======================================================

	app.Get("/", handlers.Index)

	// ==========================================
	// SSE REALTIME SCAN
	// ==========================================

	app.Get(
		"/api/scan-stream",
		handlers.StreamScan,
	)

	// ==========================================
	// API
	// ==========================================

	app.Post("/api/dns", handlers.LookupDNS)

	app.Post("/api/whois", handlers.LookupWhois)

	app.Post("/api/ssl", handlers.LookupSSL)

	// ======================================================
	// START SERVER
	// ======================================================

	log.Println(
		"security-dashboard listening on :8080",
	)

	log.Fatal(app.Listen(":8080"))
}
