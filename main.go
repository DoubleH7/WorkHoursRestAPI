package main

import (
	"fmt"
	"os"

	"github.com/DoubleH7/presenceHoursLog/database"
	"github.com/DoubleH7/presenceHoursLog/webservice"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	fmt.Println("initializing")

	godotenv.Load(".env")
	port := os.Getenv("PORT")

	fmt.Println("Connecting to database ...")

	client, err := database.ConnectDB()
	defer database.DisconnectDB(client)

	if err != nil {
		panic(err)
	}

	// sampleUser := database.User{
	// 	FullName: "Hesamoddin Haddadadel",
	// 	Username: "TheDoubleH",
	// 	Password: "Hesam123",
	// 	Access:   3,
	// 	Age:      24,
	// }
	// err = database.InsertUser(client, sampleUser)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	fmt.Println("Setting up server routes...")

	e := echo.New()

	authenticated := e.Group("/user")

	// adding basic authentication to various endpoints
	authenticated.Use(middleware.BasicAuth(
		webservice.UserpassCheck(client)),
	)

	//testing login system
	authenticated.GET("/myinfo", webservice.MyInfo)

	// Setting up admin logs
	fmt.Print("admin log setup...")
	file, err := os.OpenFile("./admin_access.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("failed")
		panic(err)
	}
	fmt.Println("")

	adminControl := authenticated.Group("/admin")
	adminControl.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}]   ${Status} from ${remote_ip}  ${method}${path}  ${latency_human}` + "\n\t" + `login_info: ${header:Authorization}` + "\n",
		Output: file,
	}))

	//hooking up the server check handler
	e.GET("/", webservice.ServerAlive)

	// // hooking up the admin handlers
	adminControl.GET("/all", webservice.GetUsers(client))
	// adminGroup.GET("/users/id/:id", webService.GetUserbyid(client))
	// adminGroup.GET("/users/name/:name", webService.GetUserbyname(client))
	adminControl.POST("/users", webservice.CreateUser(client))

	// // hooking up the user handlers
	authenticated.POST("/start", webservice.StartSession(client))
	authenticated.POST("/stop", webservice.StopSession(client))

	//starting server
	fmt.Println("initializing...")
	err = e.Start(fmt.Sprintf(":%s", port))
	if err != nil {
		fmt.Println("Failed")
		panic(err)
	}
}
