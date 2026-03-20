package category

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const CategoryTableCreationQuery = `
CREATE TABLE IF NOT EXISTS categories (
	category_id SERIAL PRIMARY KEY,
	name VARCHAR(100) NOT NULL UNIQUE,
	image_url TEXT,
	is_active BOOLEAN DEFAULT FALSE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

const GetAllCategoriesQuery = `
SELECT category_id, name, image_url, created_at, updated_at FROM categories;
`

const CreateCategoryQuery = `
INSERT INTO categories (name, image_url)
VALUES ($1, $2)
RETURNING category_id;
`

const UpdateCategoryQuery = `
UPDATE categories
SET name = $1, image_url = $2, updated_at = CURRENT_TIMESTAMP
WHERE category_id = $3;
`

const DeleteCategoryQuery = `
DELETE FROM categories
WHERE category_id = $1 RETURNING image_url;
`

func RunCategoryTableCreationQuery(conn *pgxpool.Pool) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := conn.Exec(ctx, CategoryTableCreationQuery)
	if err != nil {
		return err
	}
	return nil
}

func GetAllCategoriesFromDB(conn *pgxpool.Pool) ([]Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := conn.Query(ctx, GetAllCategoriesQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category

	for rows.Next() {
		var c Category
		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.ImageURL,
			&c.CreatedAt,
			&c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}

func CreateCategoryInDB(conn *pgxpool.Pool, category Category, image_url string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var newID int

	err := conn.QueryRow(ctx, CreateCategoryQuery,
		category.Name,
		image_url,
	).Scan(&newID)
	if err != nil {
		return 0, err
	}
	return newID, nil
}

func UpdateCategoryInDB(conn *pgxpool.Pool, category Category) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := conn.Exec(ctx, UpdateCategoryQuery,
		category.Name,
		category.ImageURL,
		category.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func DeleteCategoryInDBAndGetImageURL(conn *pgxpool.Pool, categoryID int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var imageURL string
	err := conn.QueryRow(ctx, DeleteCategoryQuery, categoryID).Scan(&imageURL)
	if err != nil {
		return "", err
	}
	return imageURL, nil
}
