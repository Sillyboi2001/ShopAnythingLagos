package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Product struct {
  ID string `json:"id"`
  NAME string `json:"name"`
  DESCRIPTION string `json:"description"`
  PRICE float64 `json:"price"`
  MERCHANTID string `json:"merchantId"`
  DATE time.Time `json:"date"`
}

var products =  []Product{}

func getProducts(c *gin.Context) {
  merchantId := c.Query("merchantId")
  if merchantId == "" {
    c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Merchant ID is required"})
    return
  }
  var merchantProducts []Product
  for _, prod := range products {
    if prod.MERCHANTID == merchantId {
      merchantProducts = append(merchantProducts, prod)
    }
}

  if len(merchantProducts) == 0 {
    c.IndentedJSON(http.StatusNotFound, gin.H{"message": "No products found for this merchant"})
    return
  }
  c.IndentedJSON(http.StatusOK, merchantProducts)
}

func createProduct(c *gin.Context) {
  var newProduct Product
  merchantId := c.Query("merchantId")
  println(merchantId)

  if err := c.BindJSON(&newProduct); err != nil {
    c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
    return
  }
  if merchantId == "" {
    c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Merchant ID is required"})
    return
  }
  newProduct.ID = uuid.NewString()
  newProduct.MERCHANTID = merchantId
  newProduct.DATE = time.Now()
  products = append(products, newProduct)
  c.IndentedJSON(http.StatusCreated, newProduct)
}

func getProductById(c *gin.Context) {
  id := c.Param("id")
  merchantId := c.Query("merchantId")
  for _, prod := range products {
    if prod.ID == id {
      if prod.MERCHANTID != merchantId {
        c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
        return
      }
      c.IndentedJSON(http.StatusOK, prod)
      return
    }
  }
  c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Product not found"})
}

func updateProduct(c *gin.Context) {
  id := c.Param("id")
  merchantId := c.Query("merchantId")
  var updateProduct Product

    if err := c.BindJSON(&updateProduct); err != nil {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
        return
    }

    for i, prod := range products {
      if prod.ID == id {
        if prod.MERCHANTID != merchantId {
          c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
          return
        }
        if updateProduct.NAME != "" {
            products[i].NAME = updateProduct.NAME
        }
        if updateProduct.DESCRIPTION != "" {
            products[i].DESCRIPTION = updateProduct.DESCRIPTION
        }
        if updateProduct.PRICE != 0 {
            products[i].PRICE = updateProduct.PRICE
        }
        c.IndentedJSON(http.StatusOK, products[i])
        return
      }
    }
    c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Product not found"})
}

func deleteProduct(c *gin.Context) {
  id := c.Param("id")
  merchantId := c.Query("merchantId")

  for i, prod := range products {
      if prod.ID == id {
        if prod.MERCHANTID != merchantId {
          c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
          return
        }
          products = append(products[:i], products[i+1:]...)
          c.IndentedJSON(http.StatusOK, gin.H{"message": "Product deleted"})
          return
      }
  }
  c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Product not found"})
}

func main() {
  router := gin.Default()
  router.Use(cors.Default())

  //API documentation routes
  router.Static("/api/docs", "./assets")
  //API requests routes
  router.GET("products", getProducts)
  router.GET("/products/:id", getProductById)
  router.POST("/products", createProduct)
  router.PATCH("/products/:id", updateProduct)
  router.DELETE("/products/:id", deleteProduct)

  router.Run("localhost:8080")
}