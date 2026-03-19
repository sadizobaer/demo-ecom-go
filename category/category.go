package category

import (
	"ecommerce/utilities"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetAllCategories(w http.ResponseWriter, r *http.Request, conn *pgxpool.Pool) {
	isValidMethod := utilities.IsMethodValid(r.Method, []string{"GET"})
	if !isValidMethod {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err := RunCategoryTableCreationQuery(conn)
	if err != nil {
		print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data, err := GetAllCategoriesFromDB(conn)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if data == nil {
		data = []Category{}
	}

	dataMap := map[string]interface{}{"categories": data}
	json.NewEncoder(w).Encode(dataMap)

}

func CreateCategory(w http.ResponseWriter, r *http.Request, conn *pgxpool.Pool) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if utilities.IsFileSizeLimitExceeded(w, r, 5*1024*1024) {
		http.Error(w, "File size limit exceeded or parse error", http.StatusBadRequest)
		return
	}
	defer r.MultipartForm.RemoveAll()

	var newCategory Category
	newCategory.Name = r.FormValue("name")
	if newCategory.Name == "" {
		http.Error(w, "Category name is required", http.StatusBadRequest)
		return
	}

	uploadedFile, header, err := r.FormFile("image")
	var fileUrlPath string
	if err == nil {
		defer uploadedFile.Close()
		fileUrl, err := utilities.SaveUploadedFile(uploadedFile, header, "categories")
		if err != nil {
			http.Error(w, "Error saving uploaded file", http.StatusInternalServerError)
			return
		}
		fileUrlPath = fileUrl
	}
	newCategory.ImageURL = fileUrlPath

	if newCategory.Name == "" {
		http.Error(w, "Category name is required", http.StatusBadRequest)
		return
	} else if newCategory.ImageURL == "" {
		http.Error(w, "Category image is required", http.StatusBadRequest)
		return
	}

	err = RunCategoryTableCreationQuery(conn)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	_, err = CreateCategoryInDB(conn, newCategory, fileUrlPath)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Category created successfully!"})
}

func UpdateCategory(w http.ResponseWriter, r *http.Request, conn *pgxpool.Pool) {
	defer r.Body.Close()

	isValidMethod := utilities.IsMethodValid(r.Method, []string{"PUT"})

	if !isValidMethod {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.PathValue("id")
	categoryID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var category Category
	err = json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	category.ID = categoryID
	err = RunCategoryTableCreationQuery(conn)
	if err != nil {
		print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = UpdateCategoryInDB(conn, category)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Category updated successfully!"})

}

func DeleteCategory(w http.ResponseWriter, r *http.Request, conn *pgxpool.Pool) {
	defer r.Body.Close()

	isValidMethod := utilities.IsMethodValid(r.Method, []string{"DELETE"})
	if !isValidMethod {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.PathValue("id")
	categoryID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	err = RunCategoryTableCreationQuery(conn)
	if err != nil {
		print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	imageURL, err := DeleteCategoryInDBAndGetImageURL(conn, categoryID)
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
	json.NewEncoder(w).Encode(map[string]string{"message": "Category deleted successfully!"})
}
