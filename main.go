package main

import (
	"TESTE_API_GO/configs"
	"TESTE_API_GO/models"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := configs.Connect()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("testedb")

	users := make([]models.User, 0)
	err = db.Raw("select * from users").Scan(&users).Error

	if err != nil {
		fmt.Println(err)
		return
	}

	/* fmt.Println("Nome: ", users[0].Name)
	fmt.Println("e-mail: ", users[0].Email)
	fmt.Println("Sexo: ", users[0].Sexo)
	fmt.Println("Quantia: ", users[0].Amount) */

	/*fmt.Println("Nascimento: ", users[0].Birth)*/
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
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Verificar se o user está apto a comprar o produto
	var user models.User
	err = db.Where("userid = ?", checkout.UserID).First(&user).Error
	if err != nil {
		fmt.Println(err)
		return
	}

	if user.Amount < product.Price {
		fmt.Println(err)
		return
	}

	// Atualize o saldo do usuário
	err = db.Exec("UPDATE users SET quantia=quantia-? WHERE user_id=?", product.Price, user.ID).Error
	if err != nil {
		return
	}

	// Descontar a quantidade do produto no estoque
	err = db.Exec("UPDATE products SET quantidade = quantidade - 1 WHERE products_id = ?", product.ID).Error
	if err != nil {
		db.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Retornar uma resposta de sucesso
	w.WriteHeader(http.StatusOK)

}
