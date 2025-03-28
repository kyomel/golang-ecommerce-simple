package model

import (
	"database/sql"
	"errors"
)

type Product struct {
	ID        string `json:"id" binding:"len=0"`
	Name      string `json:"name"`
	Price     int64  `json:"price"`
	IsDeleted *bool  `json:"is_deleted,omitempty"`
}

var (
	ErrDBNil = errors.New("connection is not available")
)

func SelectProduct(db *sql.DB) ([]Product, error) {
	if db == nil {
		return nil, ErrDBNil
	}

	query := `
		SELECT id, name, price
		FROM products
		WHERE is_deleted = false;
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	products := []Product{}
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil

}

func SelectProductByID(db *sql.DB, id string) (Product, error) {
	if db == nil {
		return Product{}, ErrDBNil
	}

	query := `
		SELECT id, name, price
		FROM products
		WHERE id = $1 AND is_deleted = false;
	`

	var product Product
	row := db.QueryRow(query, id)
	if err := row.Scan(&product.ID, &product.Name, &product.Price); err != nil {
		return Product{}, err
	}

	return product, nil
}

func InsertProduct(db *sql.DB, product Product) error {
	if db == nil {
		return ErrDBNil
	}

	query := `
		INSERT INTO products (id, name, price)
		VALUES ($1, $2, $3);
	`

	_, err := db.Exec(query, product.ID, product.Name, product.Price)
	if err != nil {
		return err
	}

	return nil
}

func UpdateProduct(db *sql.DB, product Product) error {
	if db == nil {
		return ErrDBNil
	}

	query := `
		UPDATE products
		SET name = $2, price = $3
		WHERE id = $1;
	`

	_, err := db.Exec(query, product.ID, product.Name, product.Price)
	if err != nil {
		return err
	}

	return nil
}

func DeleteProduct(db *sql.DB, id string) error {
	if db == nil {
		return ErrDBNil
	}

	query := `
		UPDATE products
		SET is_deleted = true
		WHERE id = $1;
	`

	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
