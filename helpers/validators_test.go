package helpers

import (
	"api-desafio-kvr/proto"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func returnMockProtoModelCreateCrypto() proto.CreateCryptoReq {
	return proto.CreateCryptoReq{
		Name:     "Created Crypto Test",
		AssetId:  "TCR",
		PriceUsd: 1.5,
	}
}

func returnMockProtoModelToEditCreateCrypto() proto.EditCryptoReq {
	return proto.EditCryptoReq{
		Id:       primitive.NewObjectID().Hex(),
		Name:     "Edited Crypto Test",
		AssetId:  "tce",
		PriceUsd: 0,
	}
}

func returnMockProtoModelToSortCryptos() proto.SortCryptosReq {
	return proto.SortCryptosReq{
		FieldSort: "name",
		OrderBy:   true,
	}
}

func TestValidatorInCreateCryptoWithNameInvalid(t *testing.T) {
	crypto := returnMockProtoModelCreateCrypto()
	crypto.Name = "Crypto 123"

	err := ValidatorInCreateCrypto(&crypto)
	assert.NotNil(t, err)
	assert.Equal(t, "name is invalid: Crypto 123", err.Error())
}

func TestValidatorInCreateCryptoWithAssetIdInvalid(t *testing.T) {
	crypto := returnMockProtoModelCreateCrypto()
	crypto.AssetId = "a"

	err := ValidatorInCreateCrypto(&crypto)
	assert.NotNil(t, err)
	assert.Equal(t, "asset_id is invalid: a", err.Error())
}

func TestValidatorInCreateCryptoWithPriceUsdInvalid(t *testing.T) {
	crypto := returnMockProtoModelCreateCrypto()
	crypto.PriceUsd = -5

	err := ValidatorInCreateCrypto(&crypto)
	assert.NotNil(t, err)
	assert.Equal(t, "price_usd is invalid: -5.000000", err.Error())
}

func TestValidatorInCreateCryptoWithSuccess(t *testing.T) {
	crypto := returnMockProtoModelCreateCrypto()

	err := ValidatorInCreateCrypto(&crypto)
	assert.Nil(t, err)
}

func TestValidatorInEditCryptoWithIdInvalid(t *testing.T) {
	crypto := returnMockProtoModelToEditCreateCrypto()
	crypto.Id = "123abc"

	err := ValidatorInEditCrypto(&crypto)
	assert.NotNil(t, err)
	assert.Equal(t, "id is invalid: 123abc err: the provided hex string is not a valid ObjectID", err.Error())
}

func TestValidatorInEditCryptoWithNameInvalid(t *testing.T) {
	crypto := returnMockProtoModelToEditCreateCrypto()
	crypto.Name = "Crypto 123"

	err := ValidatorInEditCrypto(&crypto)
	assert.NotNil(t, err)
	assert.Equal(t, "name is invalid: Crypto 123", err.Error())
}

func TestValidatorInEditCryptoWithAssetIdInvalid(t *testing.T) {
	crypto := returnMockProtoModelToEditCreateCrypto()
	crypto.AssetId = "a"

	err := ValidatorInEditCrypto(&crypto)
	assert.NotNil(t, err)
	assert.Equal(t, "asset_id is invalid: a", err.Error())
}

func TestValidatorInEditCryptoWithPriceUsdInvalid(t *testing.T) {
	crypto := returnMockProtoModelToEditCreateCrypto()
	crypto.PriceUsd = -5

	err := ValidatorInEditCrypto(&crypto)
	assert.NotNil(t, err)
	assert.Equal(t, "price_usd is invalid: -5.000000", err.Error())
}

func TestValidatorInEditCryptoWithSuccess(t *testing.T) {
	crypto := returnMockProtoModelToEditCreateCrypto()

	err := ValidatorInEditCrypto(&crypto)
	assert.Nil(t, err)
}

func TestValidatorListAllCryptosWithFieldSortEmptyEqualInvalid(t *testing.T) {
	sortParams := returnMockProtoModelToSortCryptos()
	sortParams.FieldSort = ""

	err := ValidatorListAllCryptos(&sortParams)
	assert.NotNil(t, err)
	assert.Equal(t, "field is invalid: ", err.Error())
}

func TestValidatorListAllCryptosWithFieldSortLessThanThreeCaracterEqualInvalid(t *testing.T) {
	sortParams := returnMockProtoModelToSortCryptos()
	sortParams.FieldSort = "aaa"

	err := ValidatorListAllCryptos(&sortParams)
	assert.NotNil(t, err)
	assert.Equal(t, "field is invalid: aaa", err.Error())
}

func TestValidatorListAllCryptosWithFieldSortBiggerThanThreeCaracterEqualSuccess(t *testing.T) {
	sortParams := returnMockProtoModelToSortCryptos()
	sortParams.FieldSort = "aaaa"

	err := ValidatorListAllCryptos(&sortParams)
	assert.Nil(t, err)
}

func TestValidatorListAllCryptosWithSuccess(t *testing.T) {
	sortParams := returnMockProtoModelToSortCryptos()

	err := ValidatorListAllCryptos(&sortParams)
	assert.Nil(t, err)
}

func TestIdValidatorWithInvalid(t *testing.T) {
	id := "123abc"

	err := IdValidator(id)
	assert.NotNil(t, err)
	assert.Equal(t, "id is invalid: 123abc err: the provided hex string is not a valid ObjectID", err.Error())
}

func TestIdValidatorWithEmptyAndInvalid(t *testing.T) {
	id := ""

	err := IdValidator(id)
	assert.NotNil(t, err)
	assert.Equal(t, "id is invalid:  err: the provided hex string is not a valid ObjectID", err.Error())
}

func TestIdValidatorWithLessThanTwoCaracterInvalid(t *testing.T) {
	id := "12"

	err := IdValidator(id)
	assert.NotNil(t, err)
	assert.Equal(t, "id is invalid: 12 err: the provided hex string is not a valid ObjectID", err.Error())
}

func TestIdValidatorWithSuccess(t *testing.T) {
	id := primitive.NewObjectID().Hex()

	err := IdValidator(id)
	assert.Nil(t, err)
}

func TestNameValidatorWithInvalid(t *testing.T) {
	name := "123abc"

	err := NameValidator(name)
	assert.NotNil(t, err)
	assert.Equal(t, "name is invalid: 123abc", err.Error())
}

func TestNameValidatorWithEmptyAndInvalid(t *testing.T) {
	name := ""

	err := NameValidator(name)
	assert.NotNil(t, err)
	assert.Equal(t, "name is invalid: ", err.Error())
}

func TestNameValidatorWithSuccess(t *testing.T) {
	name := "teste crypto"

	err := NameValidator(name)
	assert.Nil(t, err)
}

func TestAssetValidatorWithEmptyAndInvalid(t *testing.T) {
	asset := ""

	err := AssetValidator(asset)
	assert.NotNil(t, err)
	assert.Equal(t, "asset_id is invalid: ", err.Error())
}

func TestAssetValidatorWithInvalid(t *testing.T) {
	asset := "a"

	err := AssetValidator(asset)
	assert.NotNil(t, err)
	assert.Equal(t, "asset_id is invalid: a", err.Error())
}

func TestAssetValidatorWithSuccess(t *testing.T) {
	name := "teste crypto"

	err := AssetValidator(name)
	assert.Nil(t, err)
}

func TestPriceValidatorWithInvalid(t *testing.T) {
	price := float64(-5)

	err := PriceValidator(price)
	assert.NotNil(t, err)
	assert.Equal(t, "price_usd is invalid: -5.000000", err.Error())
}

func TestPriceValidatorWithSuccess(t *testing.T) {
	price := float64(5)

	err := PriceValidator(price)
	assert.Nil(t, err)
}

func TestSortValidatorWithInvalid(t *testing.T) {
	field := "aa"

	err := SortValidator(field)
	assert.NotNil(t, err)
	assert.Equal(t, "field is invalid: aa", err.Error())
}

func TestSortValidatorWithEmptyAndInvalid(t *testing.T) {
	field := ""

	err := SortValidator(field)
	assert.NotNil(t, err)
	assert.Equal(t, "field is invalid: ", err.Error())
}

func TestSortValidatorWithSuccess(t *testing.T) {
	field := "aaaa"

	err := SortValidator(field)
	assert.Nil(t, err)
}

func TestOrderByValidatorWithSuccess(t *testing.T) {
	field := true

	err := OrderByValidator(field)
	assert.Nil(t, err)
}
