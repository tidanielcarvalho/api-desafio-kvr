package redis

import (
	"api-desafio-kvr/helpers"
	"api-desafio-kvr/models"
	"encoding/json"

	"github.com/go-redis/redis"
)

var logger = &helpers.Log{}
var YesDeleteAll = true
var NoDeleteAll = false
var PrefixDeleteAll = "ListAll"
var nameLog = "REDIS"

// Cache são criados em todas as operações do controller (exceto exclusao de crypto)
// Toda vez que é realizada uma operação de criação/edição,
// o cache de ListAll é apagado, evitando assim um cache desatualizado.

func Connect() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return rdb
}

func Get(key string) []byte {
	logger.Debug(nameLog, "Getting cache for key: "+key)

	client := Connect()
	result, err := client.Get(key).Result()
	if err != nil {
		logger.Error(key, err.Error())
	}

	if result == "" {
		return nil
	}

	return []byte(result)
}

func Set(key string, crypto models.CryptoCurrency, deleteAll bool) error {
	logger.Debug(nameLog, "Setting cache for key: "+key)

	byteValue, err := json.Marshal(crypto)
	if err != nil {
		logger.Error(crypto.Id.Hex(), "Error in response: "+err.Error())
		return err
	}

	client := Connect()
	err = client.Set(key, string(byteValue), 0).Err()

	if deleteAll {
		err = DeleteAll()
	}

	return err
}

func SetByByte(key string, value string, deleteAll bool) error {
	logger.Debug(nameLog, "Setting cache for key: "+key)
	client := Connect()
	err := client.Set(key, value, 0).Err()

	if deleteAll {
		err = DeleteAll()
	}

	return err
}

func Del(key string) error {
	logger.Debug(nameLog, "Deleting cache for key: "+key)

	client := Connect()
	err := client.Del(key).Err()

	return err
}

func DeleteAll() error {
	logger.Debug(nameLog, "Deleting cache for "+PrefixDeleteAll)

	client := Connect()
	iter := client.Scan(0, PrefixDeleteAll+"*", 0).Iterator()
	for iter.Next() {
		err := client.Del(iter.Val()).Err()
		if err != nil {
			logger.Error(nameLog, err.Error())
		}
	}

	err := iter.Err()

	return err
}
