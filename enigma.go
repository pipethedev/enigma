package enigma

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
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

func encrypt(data []byte, passphrase string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(createHash(passphrase)))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
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

func (e *Enigmas) Add(email, appKey string) {
	ctx := context.Background()

	*e = append(*e, data{
		Email:     email,
		AppKey:    appKey,
		CreatedAt: time.Now(),
	})

	encryptedValue, _ := encrypt([]byte(appKey), email)

	fmt.Println(string(encryptedValue))

	err := rdb.Set(ctx, email, encryptedValue, 0).Err()

	if err != nil {
		panic(err)
	}
}

func getValue(key string) {
	ctx := context.Background()

	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	println(val)
}
