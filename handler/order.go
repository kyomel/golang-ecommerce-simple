package handler

import (
	"database/sql"
	"log"
	"math/rand"
	"onlinetoko/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CheckoutOrder(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: ambil data pesanan dari request
		var checkoutOrder model.Checkout
		if err := c.BindJSON(&checkoutOrder); err != nil {
			log.Printf("Terjadi kesalahan saat membaca request body: %v\n", err)
			c.JSON(400, gin.H{"error": "Data pesanan tidak valid"})
			return
		}

		ids := []string{}
		orderQty := make(map[string]int32)
		for _, o := range checkoutOrder.Products {
			ids = append(ids, o.ID)
			orderQty[o.ID] = int32(o.Quantity)
		}

		// TODO: ambil produk data dari database
		products, err := model.SelectProductIn(db, ids)
		if err != nil {
			log.Printf("Terjadi kesalahan saat mengambil produk: %v\n", err)
			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}

		// TODO: buat kata sandi
		passcode := generatePasscode(5)

		// TODO: hash kata sandi
		hashcode, err := bcrypt.GenerateFromPassword([]byte(passcode), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Terjadi kesalahan saat menghash kata sandi: %v\n", err)
			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}

		hashcodeString := string(hashcode)

		// TODO: buat order & detail
		order := model.Order{
			ID:         uuid.New().String(),
			Email:      checkoutOrder.Email,
			Address:    checkoutOrder.Address,
			Passcode:   &hashcodeString,
			GrandTotal: 0,
		}

		details := []model.OrderDetail{}

		for _, p := range products {
			total := p.Price * int64(orderQty[p.ID])

			detail := model.OrderDetail{
				ID:        uuid.New().String(),
				OrderID:   order.ID,
				ProductID: p.ID,
				Quantity:  orderQty[p.ID],
				Price:     p.Price,
				Total:     total,
			}

			details = append(details, detail)

			order.GrandTotal += total
		}

		model.CreateOrder(db, order, details)

		orderWithDetail := model.OrderWithDetail{
			Order:   order,
			Details: details,
		}

		orderWithDetail.Order.Passcode = &passcode

		c.JSON(200, orderWithDetail)
	}
}

func generatePasscode(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

	randomGenerator := rand.New(rand.NewSource(time.Now().UnixNano()))

	code := make([]byte, length)
	for i := range code {
		code[i] = charset[randomGenerator.Intn(len(charset))]
	}

	return string(code)
}

func ConfirmOrder(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: ambil id dari param
		id := c.Param("id")

		// TODO: baca request body
		var confirmReq model.Confirm
		if err := c.BindJSON(&confirmReq); err != nil {
			log.Printf("Terjadi kesalahan saat membaca request body: %v\n", err)
			c.JSON(400, gin.H{"error": "Data konfirmasi tidak valid"})
			return
		}

		// TODO: ambil data order dari database
		order, err := model.SelectOrderByID(db, id)
		if err != nil {
			log.Printf("Terjadi kesalahan pada saat membaca data pesanan: %v\n", err)
			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}

		// TODO: cocokkan kata sandi dari pesanan
		if order.Passcode == nil {
			log.Printf("Pesanan tidak memiliki kata sandi\n")
			c.JSON(400, gin.H{"error": "Pesanan tidak memiliki kata sandi"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(*order.Passcode), []byte(confirmReq.Passcode)); err != nil {
			log.Printf("Terjadi kesalahan saat cocokkan kata sandi: %v\n", err)
			c.JSON(401, gin.H{"error": "Tidak diizinkan mengakses pesanan"})
			return
		}

		// TODO: pastikan pesanan belum dibayar
		if order.PaidAt != nil {
			log.Printf("Pesanan sudah dibayar\n")
			c.JSON(400, gin.H{"error": "Pesanan sudah dibayar"})
			return
		}

		// TODO: cocokkan jumlah pembayaran
		if order.GrandTotal != confirmReq.Amount {
			log.Printf("Jumlah pembayaran tidak sesuai: %d\n", confirmReq.Amount)
			c.JSON(400, gin.H{"error": "Jumlah pembayaran tidak sesuai"})
			return
		}

		// TODO update informasi pesanan
		current := time.Now()
		if err = model.UpdateOrderByID(db, id, confirmReq, current); err != nil {
			log.Printf("Terjadi kesalahan saat mengubah data pesanan: %v\n", err)
			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}

		// Password digunakan satu kali lalu dihapus
		order.Passcode = nil

		order.PaidAt = &current
		order.PaidBank = &confirmReq.Bank
		order.PaidAccountNumber = &confirmReq.AccountNumber

		c.JSON(200, order)
	}
}

func GetOrder(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: ambil id dari param
		id := c.Param("id")

		// TODO: ambil dari request param
		passcode := c.Query("passcode")

		// TODO: ambil data order dari database
		order, err := model.SelectOrderByID(db, id)
		if err != nil {
			log.Printf("Terjadi kesalahan pada saat membaca data pesanan: %v\n", err)
			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}

		if order.Passcode == nil {
			log.Printf("Pesanan tidak memiliki kata sandi\n")
			c.JSON(400, gin.H{"error": "Pesanan tidak memiliki kata sandi"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(*order.Passcode), []byte(passcode)); err != nil {
			log.Printf("Terjadi kesalahan saat cocokkan kata sandi: %v\n", err)
			c.JSON(401, gin.H{"error": "Tidak diizinkan mengakses pesanan"})
			return
		}

		order.Passcode = nil

		c.JSON(200, order)
	}
}
