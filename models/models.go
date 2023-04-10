package models

type User struct {
	ID    int    `gorm:"column:user_id"`
	Name  string `gorm:"column:nome"`
	Email string `gorm:"column:email"`
	/* Birth  time.Time `gorm:"column:nascimento"`*/
	Sexo   string  ` gorm:"column:sexo"`
	Amount float64 `gorm:"column:quantia"`
}

func (User) TableName() string {
	return "users"
}

type Category struct {
	ID   int    `gorm:"column:categories_Id"`
	Name string `gorm:"column:nome"`
}

type Product struct {
	ID         int     `gorm:"column:products_id"`
	Name       string  `gorm:"column:nome"`
	Price      float64 `gorm:"column:preco"`
	Quantity   int     `gorm:"column:quantidade"`
	CategoryID int     `gorm:"column:idcategories"`
}

type Payment struct {
	UserID    int
	ProductID int
	Amount    float64
}

type Checkout struct {
	UserID    int `json:"user_id" gorm:"column:user_id"`
	ProductID int `json:"products_id" gorm:"column:products_id"`
}
