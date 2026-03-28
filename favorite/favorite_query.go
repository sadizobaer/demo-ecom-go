package favorite

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const FavoriteTableCreationQuery = `
CREATE TABLE IF NOT EXISTS favorites (
	id SERIAL PRIMARY KEY,	
	user_id INT NOT NULL,
	product_id INT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
	FOREIGN KEY (product_id) REFERENCES products(product_id) ON DELETE CASCADE
);
`

const AddToFavoriteQuery = `
INSERT INTO favorites (user_id, product_id)
VALUES ($1, $2)
RETURNING id;
`

const GetFavoriteByUserIDQuery = `
SELECT 
    f.id, 
    p.product_id, 
    p.name, 
    p.description, 
    p.price, 
    p.image_url 
FROM favorites f
JOIN products p ON f.product_id = p.product_id
WHERE f.user_id = $1;
`

const RemoveFromFavoriteQuery = `
DELETE FROM favorites
WHERE user_id = $1 AND product_id = $2;
`

func RunFavoriteTableCreationQuery(conn *pgxpool.Pool) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := conn.Exec(ctx, FavoriteTableCreationQuery)
	if err != nil {
		return err
	}
	return nil
}

func AddToFavoriteInDB(conn *pgxpool.Pool, userID int, productID int) (int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var favoriteID int
	err := conn.QueryRow(ctx, AddToFavoriteQuery, userID, productID).Scan(&favoriteID)
	if err != nil {
		return 0, err
	}
	return favoriteID, nil
}

func GetFavoriteByUserIDFromDB(conn *pgxpool.Pool, userID int) ([]FavoriteItemView, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := conn.Query(ctx, GetFavoriteByUserIDQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var favorites []FavoriteItemView
	for rows.Next() {
		var item FavoriteItemView
		err := rows.Scan(
			&item.ID,
			&item.Product.ID,
			&item.Product.Name,
			&item.Product.Description,
			&item.Product.Price,
			&item.Product.ImageURL,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err) // Helpful for debugging
		}
		favorites = append(favorites, item)
	}

	return favorites, nil
}

func RemoveFromFavoriteInDB(conn *pgxpool.Pool, userID int, productID int) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := conn.Exec(ctx, RemoveFromFavoriteQuery, userID, productID)
	if err != nil {
		return err
	}
	return nil
}
