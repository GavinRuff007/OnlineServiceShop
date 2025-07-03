package main

import (
	"RestGoTest/httpserver"
	"fmt"
	"time"
)

func printBanner() {
	fmt.Println("═══════════════════════════════════════════════")
	fmt.Println("🧩   RestGoTest API Server")
	fmt.Println("🚀   Golang + SQLite + REST API")
	fmt.Println("═══════════════════════════════════════════════")
}

func main() {

	printBanner()

	a := &httpserver.App{Port: ":9004"}
	a.Init()

	fmt.Printf("✅ Server successfully started on port %s\n", a.Port)
	fmt.Println("🟢 Running... Press Ctrl+C to stop")
	fmt.Printf("📅 Startup time: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	a.Run()
}
