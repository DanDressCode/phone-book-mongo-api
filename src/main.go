package main

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"github.com/vertex/phoneBook/mongo_db"
	"github.com/vertex/phoneBook/web_server"
	"log"
	"net/http"
)

type Configuration struct {
	dbAddressConnection string
	dbUser              string
	dbUserPassword      string
	adminName           string
	adminPassword       string
	administrator       []string
}

func main() {
	configuration := readConfiguration()

	initMongoDB(configuration.dbAddressConnection, configuration.dbUser, configuration.dbUserPassword)
	initWebServer(configuration.adminName, configuration.adminPassword)

}

func readConfiguration() Configuration {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	return Configuration{
		dbAddressConnection: viper.GetString("dbAddressConnection"),
		dbUser:              viper.GetString("dbUser"),
		dbUserPassword:      viper.GetString("dbUserPassword"),
		adminName:           viper.GetString("adminName"),
		adminPassword:       viper.GetString("adminPassword"),
	}
}

func initMongoDB(address string, user string, password string) {
	if mongo_db.CreateConnection(address, user, password) {
		fmt.Println("Connected")
	} else {
		errors.New("db connection failed")
	}

}

func initWebServer(adminUser string, adminPassword string) {

	mux := http.NewServeMux()
	mux.HandleFunc("/getAll", web_server.GetAllUsers)
	//mux.HandleFunc("/api/auth", web_server.AddBasicAuth)
	mux.HandleFunc("/api/create", web_server.CreateUser)
	mux.HandleFunc("/api/update", web_server.UpdateUser)
	mux.HandleFunc("/api/delete", web_server.DeleteUser)

	log.Println("Запуск веб-сервера на http://127.0.0.1:62961")
	err := http.ListenAndServe(":62961", mux)
	log.Fatal(err)
}
