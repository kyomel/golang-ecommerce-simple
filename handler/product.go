package handler

import (
	"database/sql"
	"errors"
	"log"
	"onlinetoko/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ListProducts(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		products, err := model.SelectProduct(db)
		if err != nil {
			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}

		c.JSON(200, products)
	}

}

func GetProducts(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: baca id dari url
		id := c.Param("id")

		// TODO: ambil dari database dengan id
		product, err := model.SelectProductByID(db, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Printf("Terjadi kesalahan pada saat mengambil data dari produk: %v\n", err)
				c.JSON(404, gin.H{"error": "Produk tidak ditemukan"})
				return
			}
			log.Printf("Terjadi kesalahan pada saat mengambil data dari produk: %v\n", err)
			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}

		// TODO: berikan response
		c.JSON(200, product)
	}
}

func CreateProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var product model.Product
		if err := c.Bind(&product); err != nil {
			log.Printf("Terjadi kesalahan pada saat membaca request body: %v\n", err)
			c.JSON(400, gin.H{"error": "Data produk tidak valid"})
			return
		}

		product.ID = uuid.New().String()

		if err := model.InsertProduct(db, product); err != nil {
			log.Printf("Terjadi kesalahan pada saat membuat produk: %v\n", err)
			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}

		c.JSON(201, product)
	}
}

func UpdateProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var product model.Product
		if err := c.Bind(&product); err != nil {
			log.Printf("Terjadi kesalahan pada saat membaca request body: %v\n", err)
			c.JSON(400, gin.H{"error": "Data produk tidak valid"})
			return
		}

		productExisting, err := model.SelectProductByID(db, id)
		if err != nil {
			log.Printf("Terjadi kesalahan pada saat mengambil data dari produk: %v\n", err)
			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}

		if product.Name != "" {
			productExisting.Name = product.Name
		}

		if product.Price != 0 {
			productExisting.Price = product.Price
		}

		if err := model.UpdateProduct(db, productExisting); err != nil {
			log.Printf("Terjadi kesalahan pada saat mengubah produk: %v\n", err)
			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}

		c.JSON(200, productExisting)
	}
}

func DeleteProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: baca id dari url
		id := c.Param("id")

		if err := model.DeleteProduct(db, id); err != nil {
			log.Printf("Terjadi kesalahan pada saat menghapus produk: %v\n", err)
			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}

		// TODO: berikan response
		c.JSON(204, gin.H{"message": "Produk berhasil dihapus"})
	}
}
