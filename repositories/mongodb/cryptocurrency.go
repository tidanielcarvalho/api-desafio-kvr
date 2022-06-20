package mongodb

import (
	"api-desafio-kvr/models"
	"api-desafio-kvr/repositories"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

const (
	DATABASE   = "kvrDb"
	COLLECTION = "cryptos"
)

func GetDataBase(client *mongo.Client) *mongo.Collection {
	return client.Database(DATABASE).Collection(COLLECTION)
}

func NameCollection() string {
	return COLLECTION
}

var InsertCryptos = func(coll IMCollection, crypto models.CryptoCurrency) (models.CryptoCurrency, error) {
	crypto.PrepateToInsert()

	result, err := coll.InsertOne(context.Background(), crypto)

	if result.InsertedID == nil {
		crypto.RevertPrepateToInsert()
		return crypto, errors.New("some error to insert")
	}

	logger.Debug(crypto.Id.Hex(), "Crypto inserted...")
	return crypto, err
}

var GetById = func(coll IMCollection, id primitive.ObjectID) (crypto models.CryptoCurrency, err error) {
	err = coll.FindOne(context.Background(), bson.M{"_id": id}).Decode(&crypto)
	logger.Debug(id.Hex(), "Crypto found...")
	return crypto, err
}

var ListAll = func(coll IMCollection, sort repositories.SortParams) (result []models.CryptoCurrency, err error) {
	field, order := OrderBy(sort)
	cursor, err := coll.Find(context.Background(), bson.M{}, options.Find().SetSort(bson.M{field: order}))
	if err != nil {
		logger.Error("", "Error in find ListAll: "+err.Error())
		return result, err
	}

	defer cursor.Close(context.Background())

	err = cursor.All(context.Background(), &result)

	logger.Debug("", "Returning cryptos...")
	return result, err
}

var UpdateCrypto = func(coll IMCollection, crypto models.CryptoCurrency) (models.CryptoCurrency, int64, error) {
	var matchedCount int64
	// SetUpsert(false) = if not exists then not insert
	opts := options.Update().SetUpsert(false)
	filter, update, err := QueryToUpdate(crypto)
	if err != nil {
		return crypto, matchedCount, err
	}

	result, err := coll.UpdateOne(context.TODO(), filter, update, opts)

	logger.Debug(crypto.Id.Hex(), "Updated crypto...")
	return crypto, result.MatchedCount, err
}

var CountDocuments = func(coll IMCollection) (count int64, err error) {
	count, err = coll.CountDocuments(context.Background(), bson.M{})

	logger.Debug("", "Count cryptos "+strconv.FormatInt(count, 10)+" ...")
	return count, err
}

var DeleteAll = func(coll IMCollection) {
	cryptos, err := ListAll(coll, repositories.SortDefault())
	if err != nil {
		logger.Error("", "Error in DeleteAll "+err.Error())
	}

	deleted := []string{}

	for i := 0; i < len(cryptos); i++ {
		logger.Debug(cryptos[i].Id.Hex(), "Deleting crypto in delete all")
		deleted = append(deleted, cryptos[i].Id.Hex())

		_, err := DeleteById(coll, cryptos[i].Id)
		if err != nil {
			logger.Error(cryptos[i].Id.Hex(), "Error in delete all")
		}
	}

	if len(deleted) < 1 {
		logger.Debug("", "Documents not deleted "+fmt.Sprint(deleted))
		return
	}

	logger.Debug("", "Deleted all documents with id "+fmt.Sprint(deleted))
}

var DeleteById = func(coll IMCollection, id primitive.ObjectID) (primitive.ObjectID, error) {
	var deletedDocument bson.M

	err := coll.FindOneAndDelete(context.Background(), bson.M{"_id": id}).Decode(&deletedDocument)

	logger.Debug(id.Hex(), "Document deleted...")
	return id, err
}

var QueryToUpdate = func(crypto models.CryptoCurrency) (where bson.M, update bson.M, err error) {
	if crypto.UpdateType == "" {
		err = errors.New("updateType is empty")
		logger.Error(crypto.Id.Hex(), err.Error())
		return bson.M{}, bson.M{}, err
	}

	switch crypto.UpdateType {
	case models.UpVote: // Increment votes
		where = bson.M{"_id": bson.M{"$eq": crypto.Id}}
		update = bson.M{"$inc": bson.M{"votes": 1}, "$set": bson.M{"updated_at": time.Now()}}

	case models.DownVote: // Decrement votes only if votes > 0
		where = bson.M{"_id": bson.M{"$eq": crypto.Id}, "votes": bson.M{"$gt": 0}}
		update = bson.M{"$inc": bson.M{"votes": -1}, "$set": bson.M{"updated_at": time.Now()}}

	default: // Trazer o UpdateOnly como default
		where = bson.M{"_id": bson.M{"$eq": crypto.Id}}
		update = bson.M{"$set": crypto.FieldsToUpdate()}
	}

	// Help to log
	whereLog, _ := json.Marshal(where)
	updateLog, _ := json.Marshal(update)
	logger.Debug(crypto.Id.Hex(), "Query to updated selected - where: "+string(whereLog)+" - update: "+string(updateLog))

	return where, update, nil
}

var OrderBy = func(sort repositories.SortParams) (string, int) {
	field := selectField(sort.Field)
	// default desc
	orderBy := -1
	if sort.Asc {
		orderBy = 1
	}
	logger.Debug("", "Filter with field: "+field+" and orderBy: "+strconv.Itoa(orderBy))
	return field, orderBy
}

var selectField = func(field string) string {
	switch field {
	case "votes", "price_usd":
		return field
	default:
		return "name"
	}
}
