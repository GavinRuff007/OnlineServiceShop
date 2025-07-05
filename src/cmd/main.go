// @title           Product API
// @version         1.0
// @description     Online Shop REST API
// @termsOfService  http://swagger.io/terms/
// @contact.name    Support
// @contact.email   example@example.com
// @license.name    MIT
// @host            202.133.88.175:8090
// @BasePath        /

package main

import (
	httpserver "RestGoTest/src"
	"RestGoTest/src/database"
	"RestGoTest/src/repository"
	"fmt"
	"time"
)

func main() {

	fmt.Println("═══════════════════════════════════════════════")
	fmt.Println("🧩      API Server      ")
	fmt.Println("🚀   OnlineShop REST API   ")
	fmt.Println("═══════════════════════════════════════════════")

	database.InitDatabase()
	defer repository.DB.Close()

	/*Start Point of Service*/
	/*═══════════════════════════════════════════════*/
	a := &httpserver.App{Port: ":9004"}
	a.Init()

	fmt.Printf("✅ Server successfully started on port %s\n", a.Port)
	fmt.Println("🟢 Running... Press Ctrl+C to stop")
	fmt.Printf("📅 Startup time: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	a.Run()
	/*═══════════════════════════════════════════════*/

}
