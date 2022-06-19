package helpers

import (
	"api-desafio-kvr/proto"
	"errors"
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var OnlyLetter = regexp.MustCompile(`^[a-z A-Z]+$`).MatchString

func ValidatorInCreateCrypto(req *proto.CreateCryptoReq) (err error) {
	err = NameValidator(req.GetName())
	if err != nil {
		return err
	}

	err = AssetValidator(req.GetAssetId())
	if err != nil {
		return err
	}

	err = PriceValidator(req.GetPriceUsd())
	if err != nil {
		return err
	}

	return nil
}

func ValidatorInEditCrypto(req *proto.EditCryptoReq) (err error) {
	err = IdValidator(req.GetId())
	if err != nil {
		return err
	}

	err = NameValidator(req.GetName())
	if err != nil {
		return err
	}

	err = AssetValidator(req.GetAssetId())
	if err != nil {
		return err
	}

	err = PriceValidator(req.GetPriceUsd())
	if err != nil {
		return err
	}

	return nil
}

func ValidatorListAllCryptos(req *proto.SortCryptosReq) error {
	err := SortValidator(req.GetFieldSort())
	if err != nil {
		return err
	}

	err = OrderByValidator(req.GetOrderBy())
	if err != nil {
		return err
	}

	return nil
}

func IdValidator(id string) error {
	_, err := primitive.ObjectIDFromHex(id)
	if id == "" || len(id) <= 2 || err != nil {
		return errors.New("id is invalid: " + id + " err: " + err.Error())
	}
	return nil
}

func NameValidator(name string) error {
	if name == "" || len(name) <= 2 || !OnlyLetter(name) {
		return errors.New("name is invalid: " + name)
	}
	return nil
}

func AssetValidator(asset string) error {
	if asset == "" || len(asset) < 2 {
		return errors.New("asset_id is invalid: " + asset)
	}
	return nil
}

func PriceValidator(price float64) error {
	priceStr := fmt.Sprintf("%f", price)
	if price < 0 || OnlyLetter(priceStr) {
		return errors.New("price_usd is invalid: " + priceStr)
	}
	return nil
}

func SortValidator(field string) error {
	if field == "" || len(field) <= 3 {
		return errors.New("field is invalid: " + field)
	}
	return nil
}

func OrderByValidator(orderBy bool) error {
	if orderBy != true && orderBy != false {
		return errors.New("orderBy is invalid")
	}
	return nil
}
