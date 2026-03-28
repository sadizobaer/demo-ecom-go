package order

import (
	"context"
	"ecommerce/cart"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const CreateOrderTableQuery = `
CREATE TABLE IF NOT EXISTS orders (
	order_id SERIAL PRIMARY KEY,
	user_id INT NOT NULL,
	cart_id INT NOT NULL,
	total_amount NUMERIC(10, 2) NOT NULL,
	status VARCHAR(20) DEFAULT 'pending',
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
	FOREIGN KEY (cart_id) REFERENCES carts(cart_id) ON DELETE CASCADE
);
`

const CreateOrderQuery = `
INSERT INTO orders (user_id, cart_id, total_amount)
VALUES ($1, $2, $3)
RETURNING order_id;
`

const GetCartItemsByUserIDQuery = `
SELECT 
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

const UpdateCartToOrderedQuery = `
UPDATE carts
SET is_ordered = TRUE, updated_at = CURRENT_TIMESTAMP
WHERE cart_id = $1;
`

const GetOrdersByUserIDQuery = `
SELECT 
	order_id, 
	user_id, 
	cart_id, 
	total_amount, 
	status
FROM orders
WHERE user_id = $1;
`

func RunOrderTableCreationQuery(conn *pgxpool.Pool) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := conn.Exec(ctx, CreateOrderTableQuery)
	if err != nil {
		return err
	}
	return nil
}

func GetCartItemsByUserID(conn *pgxpool.Pool, userID int) ([]cart.CartItem, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := conn.Query(ctx, GetCartItemsByUserIDQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cartItems []cart.CartItem
	for rows.Next() {
		var item cart.CartItem
		// Scan matches the 8 columns in GetCartItemsByUserIDQuery
		err := rows.Scan(
			&item.Product.ID,
			&item.Product.Name,
			&item.Product.Description,
			&item.Product.Price,
			&item.Product.ImageURL,
			&item.Quantity,
		)
		if err != nil {
			return nil, err
		}
		cartItems = append(cartItems, item)
	}

	return cartItems, nil
}

func CreateOrderInDB(conn *pgxpool.Pool, userID int, cartID int) (int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cartItems, err := GetCartItemsByUserID(conn, userID)
	if err != nil {
		return 0, err
	}

	var totalAmount float64
	for _, item := range cartItems {
		totalAmount += item.Product.Price * float64(item.Quantity)
	}

	var orderID int
	error := conn.QueryRow(ctx, CreateOrderQuery, userID, cartID, totalAmount).Scan(&orderID)
	if error != nil {
		return 0, error
	}

	_, err = conn.Exec(ctx, UpdateCartToOrderedQuery, cartID)
	if err != nil {
		return 0, err
	}

	return orderID, nil
}

func GetOrdersByUserIDFromDB(conn *pgxpool.Pool, userID int) ([]OrderItemView, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Fetch the basic order details
	// Ensure your SQL selects cart_id so we can look up the items
	rows, err := conn.Query(ctx, GetOrdersByUserIDQuery, userID)
	if err != nil {
		return []OrderItemView{}, err
	}
	defer rows.Close()

	var orders []OrderItemView

	for rows.Next() {
		var order OrderItemView
		var cartID int // We need this to fetch the items

		// Match these exactly to your SELECT statement in GetOrdersByUserIDQuery
		err := rows.Scan(
			&order.ID,          // order_id
			&order.UserID,      // user_id
			&cartID,            // cart_id (Crucial!)
			&order.TotalAmount, // total_amount
			&order.Status,      // status
		)
		if err != nil {
			return nil, fmt.Errorf("scan order error: %v", err)
		}

		// 2. Fetch the products for THIS specific cart/order
		// We reuse your GetCartItemsByUserID function but modify it slightly
		// to accept a CartID instead of a UserID if possible.
		items, err := GetCartItemsByUserID(conn, cartID)
		if err != nil {
			return nil, fmt.Errorf("fetch items error for cart %d: %v", cartID, err)
		}

		order.Products = items
		orders = append(orders, order)
	}

	return orders, nil
}
