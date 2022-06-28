package models

import (
	"api-desafio-kvr/proto"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

const (
	UpdateOnly = "UPDATE"
	UpVote     = "UPVOTE"
	DownVote   = "DOWNVOTE"
)

type CryptoCurrencies struct {
	CryptoCurrencies []CryptoCurrency `json:"cryptos" bson:"cryptos"`
}

type CryptoCurrency struct {
	Id         primitive.ObjectID `json:"id" bson:"_id"`
	Name       string             `json:"name" bson:"name"`
	AssetId    string             `json:"asset_id" bson:"asset_id"`
	PriceUsd   float64            `json:"price_usd" bson:"price_usd"`
	Votes      int32              `json:"votes" bson:"votes"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
	UpdateType string             `json:"-" bson:"-"` // Not insert in db
}

// Converter model crypto para proto crypto ?
func (c *CryptoCurrency) ToProtoCrypto() proto.CryptoCurrency {
	return proto.CryptoCurrency{
		Id:        c.Id.Hex(),
		Name:      c.Name,
		AssetId:   c.AssetId,
		PriceUsd:  c.PriceUsd,
		Votes:     int32(c.Votes),
		CreatedAt: c.CreatedAt.Format("2006-01-02T15:04:05.999Z"),
		UpdatedAt: c.UpdatedAt.Format("2006-01-02T15:04:05.999Z"),
	}
}

func (c *CryptoCurrency) PrepateToInsert() {
	c.Id = primitive.NewObjectID()
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
}

func (c *CryptoCurrency) RevertPrepateToInsert() {
	c.Id = primitive.NilObjectID
	c.CreatedAt = time.Time{}
	c.UpdatedAt = time.Time{}
}

func (c CryptoCurrency) FieldsToUpdate() bson.M {
	return bson.M{
		"name":       c.Name,
		"asset_id":   c.AssetId,
		"price_usd":  c.PriceUsd,
		"updated_at": time.Now(),
	}
}
