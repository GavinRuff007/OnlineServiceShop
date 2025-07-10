// @title           Product API
// @version         1.0
// @description     Online Shop REST API
// @termsOfService  http://swagger.io/terms/
// @contact.name    Support
// @contact.email   example@example.com
// @license.name    MIT
// @host            localhost:8090
// @BasePath        /

package main

import (
	httpserver "RestGoTest/src"
	"RestGoTest/src/cache"
	"RestGoTest/src/config"
	db "RestGoTest/src/database"
	"RestGoTest/src/database/migrations"
	"RestGoTest/src/pkg/logging"
	"fmt"
	"time"
)

func main() {

	cfg := config.GetConfig()
	logger := logging.NewLogger(cfg)

	fmt.Println("═══════════════════════════════════════════════")
	fmt.Println("🧩      API Server      ")
	fmt.Println("🚀   OnlineShop REST API   ")
	fmt.Println("═══════════════════════════════════════════════")

	InternalPort := fmt.Sprintf(":%s", cfg.Server.InternalPort)

	fmt.Printf("✅ Server successfully started on port %s\n", InternalPort)
	fmt.Println("🟢 Running... Press Ctrl+C to stop")
	fmt.Printf("📅 Startup time: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	/*Start Point of Service*/
	/*═══════════════════════════════════════════════*/

	err := cache.InitRedis(cfg)
	if err != nil {
		logger.Fatal(logging.Redis, logging.Startup, err.Error(), nil)
	}
	defer cache.CloseRedis()

	err = db.InitDb(cfg)
	migrations.InitMigrations()
	if err != nil {
		logger.Fatal(logging.Mysql, logging.Startup, err.Error(), nil)
	}
	defer db.CloseDb()

	a := &httpserver.App{Port: InternalPort}
	a.Init(cfg)
	a.Run()
	/*═══════════════════════════════════════════════*/

}
