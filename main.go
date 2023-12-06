package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "modernc.org/sqlite"
)

var db *gorm.DB
var err error

// Product struct representa um produto na loja
type Product struct {
	ID       string  `json:"id" gorm:"primary_key"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Code     string  `json:"code"`
	Category string  `json:"category"`
	Promotionalcode string `json:"promotionalcode"`
}

func main() {
	// Configurar conexão com o banco de dados SQLite
	db, err = gorm.Open("sqlite", "store.db")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()

	// Se não existe, ela é criada
	db.AutoMigrate(&Product{})

	// Inicializar o roteador do gin
	r := gin.Default()

	// Definir os endpoints
	r.GET("/products", GetProducts)
	r.GET("/products/:id", GetProduct)
	r.POST("/products", CreateProduct)
	r.PUT("/products/:id", UpdateProduct)
	r.DELETE("/products/:id", DeleteProduct)

	// Configurar um canal para capturar sinais do sistema operacional
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// Executar o servidor em uma goroutine
	go func() {
		err := r.Run(":8080")
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Aguarde sinais para encerrar o programa
	<-stopChan

	fmt.Println("Encerrando o programa...")
}

// Obter todos os produtos
func GetProducts(c *gin.Context) {
	var products []Product
	if err := db.Find(&products).Error; err != nil {
		c.AbortWithStatus(500)
		fmt.Println(err)
	} else {
		c.JSON(200, products)
	}
}

// Obter um produto por ID
func GetProduct(c *gin.Context) {
	id := c.Params.ByName("id")
	var product Product
	if err := db.Where("id = ?", id).First(&product).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, product)
	}
}

// Criar produto
func CreateProduct(c *gin.Context) {
	var product Product
	c.BindJSON(&product)

	// Gerar UUID
	product.ID = uuid.New().String()

	db.Create(&product)
	c.JSON(200, product)
}

// Atualizar produto
func UpdateProduct(c *gin.Context) {
	id := c.Params.ByName("id")
	var product Product
	if err := db.Where("id = ?", id).First(&product).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
		return
	}
	c.BindJSON(&product)
	db.Save(&product)
	c.JSON(200, product)
}

// Excluir produto
func DeleteProduct(c *gin.Context) {
	id := c.Params.ByName("id")
	var product Product
	d := db.Where("id = ?", id).Delete(&product)
	fmt.Println(d)
	c.JSON(200, gin.H{"id #" + id: "deleted"})
}
