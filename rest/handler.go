package rest

import (
	"TESTE_API_GO/configs"
	"TESTE_API_GO/models"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, friend!"))
}

func Checkout(w http.ResponseWriter, r *http.Request) {
	// Obter os dados do produto a partir da solicitação
	var checkout models.Checkout
	err := json.NewDecoder(r.Body).Decode(&checkout)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db, err := configs.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Verificar se o produto está disponível no estoque
	var product models.Product
	err = db.Where("products_id = ?", checkout.ProductID).First(&product).Error
	if err != nil {
		fmt.Println(err)
		return
	}

	if product.Quantity <= 0 {
		fmt.Println("Produto sem estoque")
		return
	}

	// Verificar se o user está apto a comprar o produto
	var user models.User
	err = db.Where("user_id = ?", checkout.UserID).First(&user).Error
	if err != nil {
		fmt.Println(err)
		return
	}

	if user.Amount < product.Price {
		fmt.Println("Saldo do cliente é insuficiente")
		return
	}

	// Atualize o saldo do usuário
	err = db.Exec("UPDATE users SET quantia=quantia-? WHERE user_id=?", product.Price, user.ID).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Descontar a quantidade do produto no estoque
	err = db.Exec("UPDATE products SET quantidade = quantidade - 1 WHERE products_id = ?", product.ID).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Retornar uma resposta de sucesso
	w.WriteHeader(http.StatusOK)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	//Criar uma variável do tipo user
	userToCreate := &models.User{}

	//Decodar de Json para o tipo User
	err := json.NewDecoder(r.Body).Decode(userToCreate)
	if err != nil {
		log.Println(err.Error())
		return
	}

	//Conectar com o DB
	db, err := configs.Connect()
	if err != nil {
		log.Println(err.Error())
		return
	}

	//db.Exec("INSERT INTO ...")

	if err = db.Create(userToCreate).Error; err != nil {
		log.Println(err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	//Conectar com o DB
	db, err := configs.Connect()
	if err != nil {
		log.Println(err.Error())
		return
	}

	users := make([]models.User, 0)

	if err = db.Find(&users).Error; err != nil {
		log.Println(err.Error())
		return
	}

	jsonUsers, err := json.Marshal(users)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonUsers)
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	//Pegando o id de forma dinâmica na url+
	userIDString := r.URL.Query().Get("id")

	//Transformando de string para int
	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		log.Println(err.Error())
		return
	}

	//Conectar com o DB
	db, err := configs.Connect()
	if err != nil {
		log.Println(err.Error())
		return
	}

	userToFind := &models.User{}

	// SELECT * FROM users as u WHERE u.user_id = 15

	if err = db.Where("user_id = ?", userID).
		First(userToFind).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(204)
			w.Write([]byte(fmt.Sprintf("este usuario nao existe!")))
			return
		}
	}

	jsonUser, err := json.Marshal(userToFind)
	if err != nil {
		log.Println(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonUser)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	//Pegando o id de forma dinâmica na url+
	userIDString := r.URL.Query().Get("id")
	// "10"

	//Transformando de string para int
	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		log.Println(err.Error())
		return
	}
	// 10

	userToUpdate := &models.User{}

	if err := json.NewDecoder(r.Body).Decode(userToUpdate); err != nil {
		log.Println(err.Error())
		return
	}
	// Nome: Paulo; Amount: 25.000

	userToUpdate.ID = userID

	// ID: 10; Nome: Paulo; Amount: 25.000

	//Conectar com o DB
	db, err := configs.Connect()
	if err != nil {
		log.Println(err.Error())
		return
	}

	if err = db.Save(userToUpdate).Error; err != nil {
		log.Println(err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	//Pegando o id de forma dinâmica na url+
	userIDString := r.URL.Query().Get("id")
	// "10"

	//Transformando de string para int
	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		log.Println(err.Error())
		return
	}
	// 10

	//Conectar com o DB
	db, err := configs.Connect()
	if err != nil {
		log.Println(err.Error())
		return
	}

	if err = db.Delete(&models.User{}, "user_id = ?", userID).Error; err != nil {
		fmt.Println(err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
