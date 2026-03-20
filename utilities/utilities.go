package utilities

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func IsMethodValid(method string, validMethods []string) bool {
	for _, m := range validMethods {
		if method == m {
			return true
		}
	}
	return false
}

func IsValidEmail(email string) bool {
	var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return emailRegex.MatchString(email)
}

func HashPassword(s string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func IsFileSizeLimitExceeded(w http.ResponseWriter, r *http.Request, maxSize int64) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxSize)
	err := r.ParseMultipartForm(maxSize)
	return err != nil
}

func SaveUploadedFile(file multipart.File, header *multipart.FileHeader, subPath string) (string, error) {
	filePath := "./uploads/" + subPath + "/" + header.Filename

	err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		return "Error creating directory", err
	}

	savedFile, err := os.Create(filePath)
	if err != nil {
		return "Error saving file", err
	}

	_, err = io.Copy(savedFile, file)
	if err != nil {
		return "Error copying file", err
	}

	defer file.Close()
	defer savedFile.Close()

	return filePath, nil
}

func GenerateTokens(username string, email string) (string, string, error) {

	accessToken, err := GenarateToken(username, email, "access")
	if err != nil {
		return "", "", err
	}

	refreshToken, err := GenarateToken(username, email, "refresh")
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil

}

func GenarateToken(username string, email string, tokenType string) (string, error) {

	payload := jwt.MapClaims{
		"username": username,
		"email":    email,
		"type":     tokenType,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	if tokenType == "refresh" {
		payload["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte("your_secret_key"))

}
