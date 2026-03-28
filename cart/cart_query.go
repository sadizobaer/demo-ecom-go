package cart

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// --- SQL Queries ---

const CartTableCreationQuery = `
CREATE TABLE IF NOT EXISTS carts (
    cart_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    is_ordered BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);`

const CartPartialIndexQuery = `
CREATE UNIQUE INDEX IF NOT EXISTS idx_one_active_cart_per_user 
ON carts (user_id) 
WHERE (is_ordered = FALSE);`

const CartItemsTableCreationQuery = `
CREATE TABLE IF NOT EXISTS cart_items (
    item_id SERIAL PRIMARY KEY,
    cart_id INTEGER REFERENCES carts(cart_id) ON DELETE CASCADE,
    product_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    UNIQUE(cart_id, product_id)
);`

const GetOrCreateCartQuery = `
INSERT INTO carts (user_id, is_ordered)
VALUES ($1, FALSE)
ON CONFLICT (user_id) WHERE (is_ordered = FALSE) 
DO UPDATE SET updated_at = CURRENT_TIMESTAMP
RETURNING cart_id;
`

const UpsertCartItemQuery = `
INSERT INTO cart_items (cart_id, product_id, quantity)
VALUES ($1, $2, $3)
ON CONFLICT (cart_id, product_id) 
DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity;`

const GetCartItemsByUserIDQuery = `
SELECT 
    c.cart_id, 
    c.user_id, 
    p.product_id, 
    p.name, 
    p.description, 
    p.price, 
    p.image_url, 
    ci.quantity
FROM carts c
JOIN cart_items ci ON c.cart_id = ci.cart_id
JOIN products p ON ci.product_id = p.product_id
WHERE c.user_id = $1 AND c.is_ordered = FALSE;`

const UpdateCartItemQuantityQuery = `
UPDATE cart_items 
SET quantity = $3 
WHERE cart_id = $1 AND product_id = $2;
`

const RemoveFromCartQuery = `
DELETE FROM cart_items 
WHERE product_id = $1 AND cart_id = $2;`

// --- Database Logic Functions ---

func RunCartTableCreationQuery(conn *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Create Carts Table
	if _, err := conn.Exec(ctx, CartTableCreationQuery); err != nil {
		return err
	}

	// 2. IMPORTANT: Remove the old strict constraint if it exists from previous code
	// This prevents the "only one order allowed" bug
	_, _ = conn.Exec(ctx, "ALTER TABLE carts DROP CONSTRAINT IF EXISTS unique_active_cart;")

	// 3. Create Partial Index
	if _, err := conn.Exec(ctx, CartPartialIndexQuery); err != nil {
		return err
	}

	// 4. Create Cart Items Table
	if _, err := conn.Exec(ctx, CartItemsTableCreationQuery); err != nil {
		return err
	}

	return nil
}

func AddProductToCartInDB(conn *pgxpool.Pool, cartItem CartItemAdd) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Use a transaction to ensure both Cart and Item are handled together
	tx, err := conn.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	var cartID int
	// 1. Get existing active cart ID or create a new one
	err = tx.QueryRow(ctx, GetOrCreateCartQuery, cartItem.UserID).Scan(&cartID)
	if err != nil {
		return 0, fmt.Errorf("failed to get/create cart: %v", err)
	}

	// 2. Add product to the cart (or update quantity if already exists)
	_, err = tx.Exec(ctx, UpsertCartItemQuery, cartID, cartItem.ProductID, cartItem.Quantity)
	if err != nil {
		return 0, fmt.Errorf("failed to upsert cart item: %v", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return cartID, nil
}

func GetCartItemsByUserIDFromDB(conn *pgxpool.Pool, userID int) (CartView, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := conn.Query(ctx, GetCartItemsByUserIDQuery, userID)
	if err != nil {
		return CartView{}, err
	}
	defer rows.Close()

	var cart CartView
	cart.UserID = userID
	cart.Items = []CartItem{} // Initialize as empty slice for valid JSON []

	for rows.Next() {
		var item CartItem
		// Scan matches the 8 columns in GetCartItemsByUserIDQuery
		err := rows.Scan(
			&cart.ID,
			&cart.UserID,
			&item.Product.ID,
			&item.Product.Name,
			&item.Product.Description,
			&item.Product.Price,
			&item.Product.ImageURL,
			&item.Quantity,
		)
		if err != nil {
			return cart, err
		}
		cart.Items = append(cart.Items, item)
	}

	return cart, nil
}

func UpdateCartItemQuantityInDB(conn *pgxpool.Pool, cartID int, productID int, newQuantity int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// If quantity is 0 or less, we should probably just remove the item
	if newQuantity <= 0 {
		return RemoveProductFromCartInDB(conn, cartID, productID)
	}

	result, err := conn.Exec(ctx, UpdateCartItemQuantityQuery, cartID, productID, newQuantity)
	if err != nil {
		return err
	}

	// Check if any row was actually updated
	if result.RowsAffected() == 0 {
		return fmt.Errorf("item not found in cart")
	}

	return nil
}

func RemoveProductFromCartInDB(conn *pgxpool.Pool, cartID int, productID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := conn.Exec(ctx, RemoveFromCartQuery, productID, cartID)
	return err
}
