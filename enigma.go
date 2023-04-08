package enigma

import (
	"context"
	"crypto/aes"
	"crypto/md5"
	"encoding/hex"
	"log"
	"regexp"
	"time"

	"github.com/go-redis/redis/v8"
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

func encrypt(plaintext string, key []byte) string {
	// create cipher
	c, err := aes.NewCipher(key)

	if err != nil {
		panic(err)
	}

	// allocate space for ciphered data
	out := make([]byte, len(plaintext))

	// encrypt
	c.Encrypt(out, []byte(plaintext))
	// return hex string
	return hex.EncodeToString(out)
}

func init() {
	ctx := context.Background()
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
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
		if len(appKey) < 32 {
			log.Fatal("Error with app key:", "AppKey must be 32 characters long")
		}
		*e = append(*e, data{
			Email:     email,
			AppKey:    appKey,
			CreatedAt: time.Now(),
		})

		encryptedValue = encrypt(email, []byte(appKey))

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
