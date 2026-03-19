package product

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const ProductTableCreationQuery = `
CREATE TABLE IF NOT EXISTS products (
    product_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL CHECK (price >= 0),
    stock_quantity INTEGER DEFAULT 0,
    category_id INTEGER NOT NULL, 
	image_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_category 
        FOREIGN KEY(category_id) 
        REFERENCES categories(category_id) 
        ON DELETE CASCADE
);
`

const GetAllProductsQuery = `
SELECT 
    p.product_id, 
    p.name, 
    COALESCE(p.description, '') AS description, 
    p.price, 
    p.stock_quantity, 
    p.image_url,
    p.created_at, 
    p.updated_at,
    c.category_id, 
    c.name, 
	c.image_url,
    c.created_at, 
    c.updated_at
FROM products p
LEFT JOIN categories c ON p.category_id = c.category_id;
`

const GetProductByIDQuery = `
SELECT 
	p.product_id,
	p.name,
	p.description,
	p.price,
	p.stock_quantity,
	p.image_url,
	p.created_at,
	p.updated_at,
	c.category_id, 
	c.name, 
	c.image_url,
	c.created_at, 
	c.updated_at
FROM products p
LEFT JOIN categories c ON p.category_id = c.category_id
WHERE p.product_id = $1;
`

const CreateProductQuery = `
INSERT INTO products (name, description, price, stock_quantity, image_url, category_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING product_id;
`

const UpdateProductQuery = `
UPDATE products
SET name = $1, description = $2, price = $3, stock_quantity = $4, image_url = $5, category_id = $6, updated_at = CURRENT_TIMESTAMP
WHERE product_id = $7;
`

const DeleteProductQuery = `DELETE FROM products WHERE product_id = $1 RETURNING image_url;`

func RunProductTableCreationQuery(conn *pgxpool.Pool) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := conn.Exec(ctx, ProductTableCreationQuery)
	if err != nil {
		return err
	}
	return nil
}

func GetAllProductsFromDB(conn *pgxpool.Pool) ([]ProductView, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := conn.Query(ctx, GetAllProductsQuery)
	if err != nil {
		log.Printf("Query Execution Error: %v", err)
		return nil, err
	}
	defer rows.Close()

	products := make([]ProductView, 0) // Ensures we return [] instead of null

	for rows.Next() {
		var p ProductView
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.Price,
			&p.Stock,
			&p.ImageURL,
			&p.CreatedAt,
			&p.UpdatedAt,
			// Nested Category fields start here
			&p.Category.ID,
			&p.Category.Name,
			&p.Category.ImageURL,
			&p.Category.CreatedAt,
			&p.Category.UpdatedAt,
		)
		if err != nil {
			log.Printf("Row Scan Error: %v", err) // Find out which column is the problem
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func GetProductByIDFromDB(conn *pgxpool.Pool, productID int) (ProductView, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var p ProductView
	err := conn.QueryRow(ctx, GetProductByIDQuery, productID).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.Stock,
		&p.ImageURL,
		&p.CreatedAt,
		&p.UpdatedAt,
		// Nested Category fields start here
		&p.Category.ID,
		&p.Category.Name,
		&p.Category.ImageURL,
		&p.Category.CreatedAt,
		&p.Category.UpdatedAt,
	)
	if err != nil {
		return ProductView{}, err
	}

	return p, nil
}

func CreateProductInDB(conn *pgxpool.Pool, product ProductCreate, image_url string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var newID int
	err := conn.QueryRow(ctx, CreateProductQuery,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		image_url,
		product.Category,
	).Scan(&newID)

	if err != nil {
		return 0, fmt.Errorf("failed to insert product: %v", err)
	}

	return newID, nil
}

func UpdateProductInDB(conn *pgxpool.Pool, product ProductCreate) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := conn.Exec(ctx, UpdateProductQuery,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.ImageURL,
		product.Category,
		product.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update product: %v", err)
	}
	return nil
}

func DeleteProductInDBAndGetImageURL(conn *pgxpool.Pool, productID int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var imageURL string
	err := conn.QueryRow(ctx, DeleteProductQuery, productID).Scan(&imageURL)
	if err != nil {
		return "", fmt.Errorf("failed to delete product: %v", err)
	}
	return imageURL, nil
}
