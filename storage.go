package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	db *redis.Client
}

func NewRedisStore() *RedisStore {
	connStr := os.Getenv("CONN_STR")
	if connStr == "" {
		log.Fatal("CONN_STR not found in environment!")
	}

	opt, err := redis.ParseURL(connStr)
	if err != nil {
		log.Fatal("Unable to connect to database! Error:", err)
	}

	db := redis.NewClient(opt)

	err = db.Set(context.Background(), "foo", "bar", 0).Err()
	if err != nil {
		log.Fatal("An error occurred, error: ", err)
	}

	_, err = db.Get(context.Background(), "foo").Result()
	if err != nil {
		log.Fatalf("Unable to retrieve value to provided key %s, error: %v", "foo", err)
	}
	log.Printf("Connected to Redis server succesfully!")
	return &RedisStore{db}
}

// saves the key value pair, hash : fullURL(value)
func (r *RedisStore) mapURL(ctx context.Context, hash, fullURL string) error {
	err := r.db.Set(ctx, hash, fullURL, time.Hour*24*15).Err()
	return err
}

// retrieves the mapped fullURL for the hash
func (r *RedisStore) getFullURL(ctx context.Context, hash string) (string, error) {
	val, err := r.db.Get(ctx, hash).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (r *RedisStore) EncodeURL(url string) (string, error) {
	hasher := sha256.New()
	key, err := r.generateKey(hasher, []byte(url))
	if err != nil {
		return "", err
	}

	_, err = r.getFullURL(context.Background(), key)
	// check if that key is already mapped to another url
	if err == nil {
		key, err = r.EncodeURL(url)
	}
	return key, nil
}

func (r *RedisStore) generateKey(hasher hash.Hash, url []byte) (string, error) {

	randomByte := make([]byte, 1)
	if _, err := rand.Read(randomByte); err != nil {
		return "", err
	}
	// writing url bytes
	if _, err := hasher.Write(url); err != nil {
		return "", err
	}
	// writing random byte
	if _, err := hasher.Write(randomByte); err != nil {
		return "", err
	}

	hash := hasher.Sum(nil)
	// use only the first 6 bytes
	key := hex.EncodeToString(hash[:6])
	return key, nil
}
