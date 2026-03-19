package product

import (
	"ecommerce/utilities"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetAllProducts(w http.ResponseWriter, r *http.Request, conn *pgxpool.Pool) {
	isValidMethod := utilities.IsMethodValid(r.Method, []string{"GET"})
	if !isValidMethod {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err := RunProductTableCreationQuery(conn)
	if err != nil {
		print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data, err := GetAllProductsFromDB(conn)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if data == nil {
		data = []ProductView{}
	}

	dataMap := map[string]interface{}{"products": data}
	json.NewEncoder(w).Encode(dataMap)
}

func GetProductByID(w http.ResponseWriter, r *http.Request, conn *pgxpool.Pool) {
	defer r.Body.Close()

	isValidMethod := utilities.IsMethodValid(r.Method, []string{"GET"})
	if !isValidMethod {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.PathValue("id")
	productID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	data, err := GetProductByIDFromDB(conn, productID)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)

}

func CreateProduct(w http.ResponseWriter, r *http.Request, conn *pgxpool.Pool) {
	defer r.Body.Close()

	isValidMethod := utilities.IsMethodValid(r.Method, []string{"POST"})
	if !isValidMethod {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if utilities.IsFileSizeLimitExceeded(w, r, 5*1024*1024) {
		http.Error(w, "File size limit exceeded or parse error", http.StatusBadRequest)
		return
	}
	defer r.MultipartForm.RemoveAll()

	var newProduct ProductCreate

	newProduct.Name = r.FormValue("name")
	newProduct.Description = r.FormValue("description")
	newProduct.Price, _ = strconv.ParseFloat(r.FormValue("price"), 64)
	newProduct.Stock, _ = strconv.Atoi(r.FormValue("stock"))
	categoryID, err := strconv.Atoi(r.FormValue("category"))
	if err != nil {
		http.Error(w, "Invalid category_id format", http.StatusBadRequest)
		return
	}
	newProduct.Category = categoryID

	uploadedFile, header, err := r.FormFile("image")
	var fileUrlPath string
	if err == nil {
		defer uploadedFile.Close()
		fileUrl, err := utilities.SaveUploadedFile(uploadedFile, header, "products")
		if err != nil {
			http.Error(w, "Error saving uploaded file", http.StatusInternalServerError)
			return
		}
		fileUrlPath = fileUrl
	}
	newProduct.ImageURL = fileUrlPath

	if newProduct.Name == "" {
		http.Error(w, "Product name is required", http.StatusBadRequest)
		return
	} else if newProduct.Price <= 0 {
		http.Error(w, "Price must be greater than zero", http.StatusBadRequest)
		return
	} else if newProduct.Category <= 0 {
		http.Error(w, "Valid category_id is required", http.StatusBadRequest)
		return
	} else if newProduct.ImageURL == "" {
		http.Error(w, "Product image is required", http.StatusBadRequest)
		return
	}

	error := RunProductTableCreationQuery(conn)
	if error != nil {
		print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data, err := CreateProductInDB(conn, newProduct, fileUrlPath)

	if err != nil {
		http.Error(w, "Internal Server Error "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)

}

func UpdateProduct(w http.ResponseWriter, r *http.Request, conn *pgxpool.Pool) {
	defer r.Body.Close()

	isValidMethod := utilities.IsMethodValid(r.Method, []string{"PUT"})
	if !isValidMethod {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.PathValue("id")
	ProductID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var product ProductCreate
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	product.ID = ProductID

	err = UpdateProductInDB(conn, product)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Product updated successfully!"})

}

func DeleteProduct(w http.ResponseWriter, r *http.Request, conn *pgxpool.Pool) {
	defer r.Body.Close()

	isValidMethod := utilities.IsMethodValid(r.Method, []string{"DELETE"})
	if !isValidMethod {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.PathValue("id")
	ProductID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	imageURL, err := DeleteProductInDBAndGetImageURL(conn, ProductID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if imageURL != "" {
		// os.Remove handles the deletion.
		// We often ignore the error here because if the file is already gone,
		// we still want the API to succeed since the DB record is deleted.
		err = os.Remove(imageURL)
		if err != nil {
			fmt.Printf("Warning: Could not delete file %s: %v\n", imageURL, err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Product deleted successfully!"})

}
