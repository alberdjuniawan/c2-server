package middleware

import (
	"c2/config"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// ValidateJWT memverifikasi token JWT pada setiap request
func ValidateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ambil token dari header Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		// Token harus dimulai dengan "Bearer "
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			http.Error(w, "Authorization token is missing", http.StatusUnauthorized)
			return
		}

		// Verifikasi token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verifikasi metode signing key
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid signing method: %v", token.Method)
			}
			// Kembalikan secret key dari konfigurasi
			return []byte(config.JWTSecret), nil
		})

		// Periksa error token atau validitas token
		if err != nil {
			log.Println("Error parsing token:", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Periksa apakah token valid
		if !token.Valid {
			log.Println("Invalid token: token is not valid")
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Jika token kadaluarsa, beri tahu pengguna
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && claims["exp"] != nil {
			expirationTime := int64(claims["exp"].(float64))
			if expirationTime < int64(time.Now().Unix()) {
				log.Println("Token expired")
				http.Error(w, "Token expired, please login again", http.StatusUnauthorized)
				return
			}
		}

		// Token valid, lanjutkan request ke handler berikutnya
		next.ServeHTTP(w, r)
	})
}
