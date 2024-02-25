package user

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/smtp"
	"os"
	"path"
	"time"

	"github.com/clfdrive/server/internal/rest"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	table     = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))

	rsaPrivateKeyRaw, _ = os.ReadFile("secrets/refresh.key")
	rsaPrivateKey, _    = jwt.ParseRSAPrivateKeyFromPEM(rsaPrivateKeyRaw)
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	return string(bytes), err
}

func comparePasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func parseRefreshToken(refreshToken string) (int, error) {
	token, err := jwt.ParseWithClaims(refreshToken, new(rest.AccessTokenClaims), func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodRS256 {
			return nil, errors.New("incorrect signing method")
		}

		return rsaPrivateKey.Public(), nil
	})
	if err != nil {
		return -1, err
	}

	fmt.Println(token.Claims.(*rest.AccessTokenClaims).UserId)

	return token.Claims.(*rest.AccessTokenClaims).UserId, nil
}

func GenerateTokenPair(userId int) (string, string, error) {
	// Create token
	claims := &rest.AccessTokenClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, &rest.AccessTokenClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 256)),
		},
	})
	rt, err := refreshToken.SignedString(rsaPrivateKey)
	if err != nil {
		return "", "", err
	}

	return t, rt, nil
}

func sendEmail(email, code string) error {
	from := os.Getenv("EMAIL")
	password := os.Getenv("EMAIL_PASSWORD")

	// Receiver email address.
	to := []string{
		email,
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	t, err := template.ParseFiles(path.Join("static", "index.html"))
	if err != nil {
		log.Println(err)
		return err
	}

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: CLFD Account Verification \n%s\n\n", mimeHeaders)))

	t.Execute(&body, struct {
		Code string
	}{
		Code: code,
	})

	// Sending email.
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func genVerifCode() string {
	b := make([]byte, 6)
	n, err := io.ReadAtLeast(rand.Reader, b, 6)
	if n != 6 {
		panic(err)
	}

	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}

	return string(b)
}
