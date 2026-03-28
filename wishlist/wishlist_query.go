package wishlist

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const WishlistTableCreationQuery = `
CREATE TABLE IF NOT EXISTS wishlists (
	id SERIAL PRIMARY KEY,	
	user_id INT NOT NULL,
	product_id INT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
	FOREIGN KEY (product_id) REFERENCES products(product_id) ON DELETE CASCADE
);
`

const AddToWishlistQuery = `
INSERT INTO wishlists (user_id, product_id)
VALUES ($1, $2)
RETURNING id;
`

const GetWishlistByUserIDQuery = `
SELECT 
    w.id, 
    p.product_id, 
    p.name, 
    p.description, 
    p.price, 
    p.image_url 
FROM wishlists w
JOIN products p ON w.product_id = p.product_id
WHERE w.user_id = $1;
`

const RemoveFromWishlistQuery = `
DELETE FROM wishlists
WHERE user_id = $1 AND product_id = $2;
`

func RunWishlistTableCreationQuery(conn *pgxpool.Pool) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := conn.Exec(ctx, WishlistTableCreationQuery)
	if err != nil {
		return err
	}
	return nil
}

func AddToWishlistInDB(conn *pgxpool.Pool, userID int, productID int) (int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var wishlistID int
	err := conn.QueryRow(ctx, AddToWishlistQuery, userID, productID).Scan(&wishlistID)
	if err != nil {
		return 0, err
	}
	return wishlistID, nil
}

func GetWishlistByUserIDFromDB(conn *pgxpool.Pool, userID int) ([]WishlistItemView, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := conn.Query(ctx, GetWishlistByUserIDQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wishlist []WishlistItemView
	for rows.Next() {
		var item WishlistItemView
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
		wishlist = append(wishlist, item)
	}

	return wishlist, nil
}

func RemoveFromWishlistInDB(conn *pgxpool.Pool, userID int, productID int) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := conn.Exec(ctx, RemoveFromWishlistQuery, userID, productID)
	if err != nil {
		return err
	}
	return nil
}
