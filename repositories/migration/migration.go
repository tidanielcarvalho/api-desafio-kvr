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
var nameLog = "MIGRATION"

func CreateInitialCryptosBulk(collection mongodb.IMCollection) {
	countDoc, err := mongodb.CountDocuments(collection)
	if err != nil {
		logger.Error(nameLog, "Error in migration CountDocuments: "+err.Error())
		return
	}

	// if exists documents in collection
	if countDoc > 1 {
		logger.Warn("", "Already exists "+strconv.FormatInt(countDoc, 10)+" cryptos in collection "+mongodb.NameCollection())
		return
	}

	// else, then import
	cryptos := GetFileToImport()

	amountCryptos := len(cryptos)
	logger.Info("", "Importing "+strconv.Itoa(amountCryptos)+" cryptos in collection "+mongodb.NameCollection())

	for i := 0; i < amountCryptos; i++ {
		cryptos[i].Id = primitive.NewObjectID()
		cryptos[i].CreatedAt = time.Now()
		cryptos[i].UpdatedAt = time.Now()

		_, err := mongodb.InsertCryptos(collection, cryptos[i])
		if err != nil {
			logger.Error("", "Error in import "+err.Error())
		}
		logger.Debug("", "Crypto "+cryptos[i].Name+" imported")
	}
}

func GetFileToImport() []models.CryptoCurrency {
	var cryptos []models.CryptoCurrency
	path := "repositories/migration/dataInitial.json"

	jsonFile, err := os.Open(path)
	if err != nil {
		logger.Error(nameLog, err.Error())
		return cryptos
	}
	logger.Info(nameLog, "Successful to open file json "+path)

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		logger.Error(nameLog, err.Error())
		return cryptos
	}

	err = json.Unmarshal(byteValue, &cryptos)
	if err != nil {
		logger.Error(nameLog, err.Error())
	}

	return cryptos
}
