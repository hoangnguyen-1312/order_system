package models

import (
	"context"
	"fmt"
	_"os"
	_"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

var (
	tokenSecret = []byte("tokensecret")
)

type User struct {
	ID              uuid.UUID `json:"id"`
	Email           string    `json:"email"`
	PasswordHash    string    `json:"-"`
	Password        string    `json:"password"`
	CreatedAt       time.Time `json:"_"`
	UpdatedAt       time.Time `json:"_"`
}

func (u *User) Register(conn *pgx.Conn) error {

	row := conn.QueryRow(context.Background(), "SELECT id from user_account WHERE email = $1", u.Email)
	userLookup := User{}
	err := row.Scan(&userLookup)
	if err != pgx.ErrNoRows {
		fmt.Println("found user")
		fmt.Println(userLookup.Email)
		return fmt.Errorf("already exists")
	}

	pwdHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	if err != nil {
		return fmt.Errorf("There was an error creating your account.")
	}
	u.PasswordHash = string(pwdHash)

	now := time.Now()
	_, err = conn.Exec(context.Background(), "INSERT INTO user_account (created_at, updated_at, email, password_hashing) VALUES($1, $2, $3, $4)", now, now, u.Email, u.PasswordHash)

	return err
}


func (u *User) GetAuthToken() (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = u.ID
	claims["exp"] = time.Now().Add(time.Hour * 24)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	authToken, err := token.SignedString(tokenSecret)
	return authToken, err
}

func (u *User) IsAuthenticated(conn *pgx.Conn) error {
	row := conn.QueryRow(context.Background(), "SELECT id, password_hashing from user_account WHERE email = $1", u.Email)
	err := row.Scan(&u.ID, &u.PasswordHash)
	if err == pgx.ErrNoRows {
		fmt.Println("User with email not found")
		return fmt.Errorf("Invalid login credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(u.Password))
	if err != nil {
		return fmt.Errorf("Invalid login credentials")
	}

	return nil
}

func IsTokenValid(tokenString string) (bool, string) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		fmt.Printf("Parsing: %v \n", token)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok == false {
			return nil, fmt.Errorf("Token signing method is not valid")
		}

		return tokenSecret, nil
	})

	if err != nil {
		fmt.Printf("Err %v \n", err)
		return false, ""
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
		userID := claims["user_id"]
		return true, userID.(string)
	} else {
		fmt.Println(err)
		return false, "uuid.UUID{}"
	}
}
