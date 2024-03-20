package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"task/mongo"
	"log"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"github.com/joho/godotenv"
)

func GenerateTokensHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Обработка параметров запроса
	userID := r.URL.Query().Get("user_id")
	w.Write([]byte("got userID"))

	// Генерация Access токена
	token := jwt.New(jwt.SigningMethodHS512)
	accessToken, _ := token.SignedString([]byte("your_secret_key"))
	w.Write([]byte("created AccessToken"))

	// Генерация Refresh токена и его хеширование
	refreshRaw := "your_refresh_token"
	refreshHash, err := bcrypt.GenerateFromPassword([]byte(refreshRaw), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	w.Write([]byte("Create RefreshToken"))

	// Сохранение refreshHash в MongoDB
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variablse: %s", err.Error())
	}
	client := mongo.GetMongoClient()
	var DB_NAME string = os.Getenv("DB_NAME")
	var DB_COLLECTIONS string = os.Getenv("DB_COLLECTIONS")

	collection := client.Database(DB_NAME).Collection(DB_COLLECTIONS)

	ctx := context.TODO()
	_, err = collection.InsertOne(ctx, bson.D{
		{Key: "userId", Value: userID},
		{Key: "refreshTokenHash", Value: refreshHash},
	})
	if err != nil {
		panic(err)
	}
	w.Write([]byte("Put in DB"))

	// Возвращение Access и Refresh токенов в формате JSON
	tokens := map[string]string{"access_token": accessToken, "refresh_token_hash": string(refreshHash)}
	json.NewEncoder(w).Encode(tokens)
}

func RefreshTokensHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Распаковка данных из тела запроса
	var tokens map[string]string
	err := json.NewDecoder(r.Body).Decode(&tokens)
	if err != nil {
		panic(err)
	}

	// Получение Refresh токена из запроса
	refreshToken := tokens["refresh_token"]

	// Поиск хеша Refresh токена в базе данных MongoDB
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variablse: %s", err.Error())
	}
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variablse: %s", err.Error())
	}
	client := mongo.GetMongoClient()
	var DB_NAME string = os.Getenv("DB_NAME")
	var DB_COLLECTIONS string = os.Getenv("DB_COLLECTIONS")

	collection := client.Database(DB_NAME).Collection(DB_COLLECTIONS)
	var result struct {
		RefreshTokenHash string
	}
	filter := bson.M{"refreshTokenHash": refreshToken}
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		panic(err)
	}

	// Проверка соответствия Refresh токена хешу в базе
	err = bcrypt.CompareHashAndPassword([]byte(result.RefreshTokenHash), []byte(refreshToken))
	if err != nil {
		panic(err)
	}

	// Если Refresh токен верный, обновление Access токена
	token := jwt.New(jwt.SigningMethodHS512)
	accessToken, _ := token.SignedString([]byte("your_secret_key"))

	// Возврат Access токена клиенту
	response := map[string]string{"access_token": accessToken}
	json.NewEncoder(w).Encode(response)
}
