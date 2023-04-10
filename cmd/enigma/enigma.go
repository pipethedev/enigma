package enigma

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/go-redis/redis/v8"
	evp "github.com/walkert/go-evp"
)

type data struct {
	Email     string
	AppKey    string
	CreatedAt time.Time
}

type Enigmas []data

var rdb *redis.Client

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func Aes256Encode(plaintext string, key string, iv string, blockSize int) string {
	bKey := []byte(key)
	bIV := []byte(iv)
	bPlaintext := PKCS5Padding([]byte(plaintext), blockSize, len(plaintext))
	block, err := aes.NewCipher(bKey)
	if err != nil {
		panic(err)
	}
	ciphertext := make([]byte, len(bPlaintext))
	mode := cipher.NewCBCEncrypter(block, bIV)
	mode.CryptBlocks(ciphertext, bPlaintext)
	return hex.EncodeToString(ciphertext)
}

func PKCS5Padding(ciphertext []byte, blockSize int, _ int) []byte {
	padding := (blockSize - len(ciphertext)%blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func encrypt(rawKey string, plainText []byte) string {
	salt := []byte("ABCDEFGH")

	key, iv := evp.BytesToKeyAES256CBCMD5([]byte(salt), []byte(rawKey))

	block, err := aes.NewCipher(key)
	if err != nil {
		return err.Error()
	}

	cipherText := make([]byte, len(plainText))

	encryptStream := cipher.NewCTR(block, iv)
	encryptStream.XORKeyStream(cipherText, plainText)

	ivHex := hex.EncodeToString(iv)
	encryptedDataHex := hex.EncodeToString([]byte("Salted__")) + hex.EncodeToString(salt) + hex.EncodeToString(cipherText)
	return ivHex + ":" + encryptedDataHex
}

func init() {
	ctx := context.Background()
	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
}

func (e *Enigmas) Add(email, appKey string) string {
	encryptedValue := ""

	ctx := context.Background()

	//Validate email address

	if isValidEmail(email) {
		if len(appKey) < 10 {
			log.Fatal("Error with app key:", "AppKey must be greater than 10 characters long")
		}
		*e = append(*e, data{
			Email:     email,
			AppKey:    appKey,
			CreatedAt: time.Now(),
		})

		encryptedValue = encrypt(appKey, []byte(email+"_"+appKey+"_"+time.Now().String()))

		err := rdb.Set(ctx, email, encryptedValue, 0).Err()

		if err != nil {
			panic(err)
		}
	}
	return encryptedValue
}

func (e *Enigmas) Get(email string) string {
	ctx := context.Background()

	val, err := rdb.Get(ctx, email).Result()

	if err != nil {
		//panic(err)
		log.Fatal("Error with email:", "Unable to fetch email")
	}

	if len(val) == 0 {
		log.Fatal("Error with email:", "Email not found")
	}
	return val
}

func isValidEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	r, err := regexp.Compile(regex)
	if err != nil {
		log.Fatal("Error compiling regular expression:", err)
		return false
	}

	if !r.MatchString(email) {
		log.Fatal("Invalid email address")
	}
	return r.MatchString(email)
}
