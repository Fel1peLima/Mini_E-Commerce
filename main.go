package main

import (
	"TESTE_API_GO/rest"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fmt.Println("Start mini-loja...")
	// Execute queries aqui...
	//Checkout
	http.HandleFunc("/checkout", rest.Checkout)
	http.HandleFunc("/", rest.HelloHandler)

	//CRUD = CREATE / READ/ UPDATE/ DELETE
	//User

	http.HandleFunc("/create_user", rest.CreateUser)
	http.HandleFunc("/get_users", rest.GetUsers)
	http.HandleFunc("/get_user/{id}", rest.GetUserByID)
	http.HandleFunc("/update_user/{id}", rest.UpdateUser)
	http.HandleFunc("/delete_user/{id}", rest.DeleteUser)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
