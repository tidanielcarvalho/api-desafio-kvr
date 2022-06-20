package migration

import (
	"api-desafio-kvr/helpers"
	"api-desafio-kvr/models"
	"api-desafio-kvr/repositories/mongodb"
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var logger = &helpers.Log{}

func CreateInitialCryptosBulk(collection mongodb.IMCollection) {
	countDoc, _ := mongodb.CountDocuments(collection)

	// if exists documents in collection
	if countDoc > 1 {
		logger.Warn("", "Already exists "+strconv.FormatInt(countDoc, 10)+" cryptos in collection "+mongodb.NameCollection())
		return
	}

	// else, then import
	cryptos := GetFileToImport()

	amountCryptos := len(cryptos.CryptoCurrencies)
	logger.Info("", "Importing "+strconv.Itoa(amountCryptos)+" cryptos in collection "+mongodb.NameCollection())

	for i := 0; i < amountCryptos; i++ {
		cryptos.CryptoCurrencies[i].Id = primitive.NewObjectID()
		cryptos.CryptoCurrencies[i].CreatedAt = time.Now()
		cryptos.CryptoCurrencies[i].UpdatedAt = time.Now()

		_, err := mongodb.InsertCryptos(collection, cryptos.CryptoCurrencies[i])
		if err != nil {
			logger.Error("", "Error in import "+err.Error())
		}
		logger.Debug("", "Crypto "+cryptos.CryptoCurrencies[i].Name+" imported")
	}
}

func GetFileToImport() models.CryptoCurrencies {
	path := "repositories/migration/dataInitial.json"
	jsonFile, err := os.Open(path)
	if err != nil {
		logger.Fatal("", err.Error(), err)
	}
	logger.Info("", "Successful to open file json "+path)

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var cryptos models.CryptoCurrencies
	json.Unmarshal(byteValue, &cryptos)

	return cryptos
}
