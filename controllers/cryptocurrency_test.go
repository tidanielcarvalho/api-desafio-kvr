package controllers

import (
	"api-desafio-kvr/models"
	"api-desafio-kvr/proto"
	"api-desafio-kvr/repositories"
	"api-desafio-kvr/repositories/mongodb"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
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

func returnMockProtoModelToDeleteCrypto() proto.DeleteCryptoReq {
	return proto.DeleteCryptoReq{
		Id: primitive.NewObjectID().Hex(),
	}
}

func returnMockProtoModelToFindCrypto() proto.FindCryptoReq {
	return proto.FindCryptoReq{
		Id: primitive.NewObjectID().Hex(),
	}
}

func returnMockProtoModelToVote() proto.VoteReq {
	return proto.VoteReq{
		Id: primitive.NewObjectID().Hex(),
	}
}

func returnMockProtoModelToMonitorVotes() proto.MonitorVotesReq {
	return proto.MonitorVotesReq{
		Id: primitive.NewObjectID().Hex(),
	}
}

func returnMockProtoModelToSortCryptos() proto.SortCryptosReq {
	return proto.SortCryptosReq{
		FieldSort: "name",
		OrderBy:   true,
	}
}

func returnMockDbListAll() []models.CryptoCurrency {
	cryptos := []models.CryptoCurrency{
		{
			Id:        primitive.NewObjectID(),
			Name:      "test name 2",
			AssetId:   "tn2",
			PriceUsd:  2.5,
			Votes:     int32(10),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Id:        primitive.NewObjectID(),
			Name:      "test name 1",
			AssetId:   "tn1",
			PriceUsd:  1.5,
			Votes:     int32(5),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	return cryptos
}

func returnMockModelCryptoCurrency() models.CryptoCurrency {
	return models.CryptoCurrency{
		Id:        primitive.NewObjectID(),
		Name:      "Created Crypto Test",
		AssetId:   "TCR",
		PriceUsd:  1.5,
		Votes:     0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func returnMockModelCryptoCurrencyEmpty() models.CryptoCurrency {
	return models.CryptoCurrency{}
}

type Mock_EndPointCryptos_MonitorVotesServer struct {
	grpc.ServerStream
	Results []*proto.CryptoCurrency
}

func (mock *Mock_EndPointCryptos_MonitorVotesServer) Send(crypto *proto.CryptoCurrency) error {
	mock.Results = append(mock.Results, crypto)
	return nil
}

// Testing crypto create with invalid name
func TestCreateCryptoWithNameInvalid(t *testing.T) {
	server := AppServer{}
	crypto := returnMockProtoModelCreateCrypto()
	crypto.Name = "Crypto 123"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	result, err := server.CreateCrypto(ctx, &crypto)

	assert.Equal(t, "rpc error: code = InvalidArgument desc = name is invalid: Crypto 123", err.Error())
	assert.Empty(t, result.Id)
	assert.Empty(t, result.CreatedAt)
	assert.Empty(t, result.UpdatedAt)

	defer cancel()
}

// Testing crypto create with invalid asset_id
func TestCreateCryptoWithAssetIdInvalid(t *testing.T) {
	server := AppServer{}
	crypto := returnMockProtoModelCreateCrypto()
	crypto.AssetId = "a"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	result, err := server.CreateCrypto(ctx, &crypto)

	assert.Equal(t, "rpc error: code = InvalidArgument desc = asset_id is invalid: a", err.Error())
	assert.Empty(t, result.Id)
	assert.Empty(t, result.CreatedAt)
	assert.Empty(t, result.UpdatedAt)

	defer cancel()
}

// Testing crypto create with invalid price_usd
func TestCreateCryptoWithPriceUsdInvalid(t *testing.T) {
	server := AppServer{}
	crypto := returnMockProtoModelCreateCrypto()
	crypto.PriceUsd = -5

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	result, err := server.CreateCrypto(ctx, &crypto)

	assert.Equal(t, "rpc error: code = InvalidArgument desc = price_usd is invalid: -5.000000", err.Error())
	assert.Empty(t, result.Id)
	assert.Empty(t, result.CreatedAt)
	assert.Empty(t, result.UpdatedAt)

	defer cancel()
}

// Testing create crypto with error
func TestCreateCryptoWithError(t *testing.T) {
	server := AppServer{}
	crypto := returnMockProtoModelCreateCrypto()

	mongodb.InsertCryptos = func(coll mongodb.IMCollection, crypto models.CryptoCurrency) (models.CryptoCurrency, error) {
		return models.CryptoCurrency{}, errors.New("test create error")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	result, err := server.CreateCrypto(ctx, &crypto)

	assert.NotNil(t, err)
	assert.Equal(t, "rpc error: code = Internal desc = test create error", err.Error())
	assert.Empty(t, result.Id)
	assert.Empty(t, result.Votes)
	assert.Empty(t, result.CreatedAt)
	assert.Empty(t, result.UpdatedAt)

	defer cancel()
}

// Testing create crypto successful
func TestCreateCryptoWithSuccess(t *testing.T) {
	server := AppServer{}
	crypto := returnMockProtoModelCreateCrypto()

	mongodb.InsertCryptos = func(coll mongodb.IMCollection, crypto models.CryptoCurrency) (models.CryptoCurrency, error) {
		return returnMockModelCryptoCurrency(), nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	result, err := server.CreateCrypto(ctx, &crypto)

	assert.Nil(t, err)
	assert.NotEmpty(t, result.Id)
	assert.Equal(t, crypto.Name, result.Name)
	assert.Equal(t, crypto.AssetId, result.AssetId)
	assert.Equal(t, crypto.PriceUsd, result.PriceUsd)
	assert.Empty(t, result.Votes)
	assert.NotEmpty(t, result.CreatedAt)
	assert.NotEmpty(t, result.UpdatedAt)

	defer cancel()
}

// Testing edit crypto with invalid id
func TestEditCryptoWithIdInvalid(t *testing.T) {
	server := AppServer{}
	crypto := returnMockProtoModelToEditCreateCrypto()
	crypto.Id = "123abc"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	result, err := server.EditCrypto(ctx, &crypto)

	assert.Equal(t, "rpc error: code = InvalidArgument desc = id is invalid: 123abc err: the provided hex string is not a valid ObjectID", err.Error())
	assert.Empty(t, result.Id)
	assert.Empty(t, result.CreatedAt)
	assert.Empty(t, result.UpdatedAt)

	defer cancel()
}

// Testing edit crypto with invalid name
func TestEditCryptoWithNameInvalid(t *testing.T) {
	server := AppServer{}
	crypto := returnMockProtoModelToEditCreateCrypto()
	crypto.Name = "Crypto 123"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	result, err := server.EditCrypto(ctx, &crypto)

	assert.Equal(t, "rpc error: code = InvalidArgument desc = name is invalid: Crypto 123", err.Error())
	assert.Empty(t, result.Id)
	assert.Empty(t, result.CreatedAt)
	assert.Empty(t, result.UpdatedAt)

	defer cancel()
}

// Testing edit crypto with invalid asset_id
func TestEditCryptoWithAssetIdInvalid(t *testing.T) {
	server := AppServer{}
	crypto := returnMockProtoModelToEditCreateCrypto()
	crypto.AssetId = "a"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	result, err := server.EditCrypto(ctx, &crypto)

	assert.Equal(t, "rpc error: code = InvalidArgument desc = asset_id is invalid: a", err.Error())
	assert.Empty(t, result.Id)
	assert.Empty(t, result.CreatedAt)
	assert.Empty(t, result.UpdatedAt)

	defer cancel()
}

// Testing edit crypto with invalid price_usd
func TestEditCryptoWithPriceUsdInvalid(t *testing.T) {
	server := AppServer{}
	crypto := returnMockProtoModelToEditCreateCrypto()
	crypto.PriceUsd = -5

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	result, err := server.EditCrypto(ctx, &crypto)

	assert.Equal(t, "rpc error: code = InvalidArgument desc = price_usd is invalid: -5.000000", err.Error())
	assert.Empty(t, result.Id)
	assert.Empty(t, result.CreatedAt)
	assert.Empty(t, result.UpdatedAt)

	defer cancel()
}

// Testing edit crypto with update error
func TestEditCryptoWithUpdateCryptoError(t *testing.T) {
	server := AppServer{}
	crypto := returnMockProtoModelToEditCreateCrypto()

	mongodb.UpdateCrypto = func(coll mongodb.IMCollection, crypto models.CryptoCurrency) (models.CryptoCurrency, int64, error) {
		return models.CryptoCurrency{Id: crypto.Id}, 0, errors.New("test update error")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	result, err := server.EditCrypto(ctx, &crypto)

	assert.NotNil(t, err)
	assert.Equal(t, "rpc error: code = Internal desc = test update error", err.Error())
	assert.NotNil(t, result.Id)
	assert.NotNil(t, result.Votes)
	assert.Empty(t, result.CreatedAt)
	assert.Empty(t, result.UpdatedAt)

	defer cancel()
}

// Testing edit crypto with getbyid error
func TestEditCryptoWithGetByIdError(t *testing.T) {
	server := AppServer{}
	crypto := returnMockProtoModelToEditCreateCrypto()

	mongodb.UpdateCrypto = func(coll mongodb.IMCollection, crypto models.CryptoCurrency) (models.CryptoCurrency, int64, error) {
		return models.CryptoCurrency{Id: crypto.Id}, 1, nil
	}

	mongodb.GetById = func(coll mongodb.IMCollection, id primitive.ObjectID) (crypto models.CryptoCurrency, err error) {
		return returnMockModelCryptoCurrencyEmpty(), errors.New("test getbyid error")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	result, err := server.EditCrypto(ctx, &crypto)

	assert.NotNil(t, err)
	assert.Equal(t, "rpc error: code = NotFound desc = test getbyid error", err.Error())
	assert.Empty(t, result.Id)
	assert.Empty(t, result.Votes)
	assert.Empty(t, result.CreatedAt)
	assert.Empty(t, result.UpdatedAt)

	defer cancel()
}

// Testing edit crypto successful
func TestEditCryptoWithSuccess(t *testing.T) {
	server := AppServer{}
	crypto := returnMockProtoModelToEditCreateCrypto()

	mongodb.UpdateCrypto = func(coll mongodb.IMCollection, crypto models.CryptoCurrency) (models.CryptoCurrency, int64, error) {
		return models.CryptoCurrency{Id: crypto.Id}, 1, nil
	}

	mongodb.GetById = func(coll mongodb.IMCollection, id primitive.ObjectID) (models.CryptoCurrency, error) {
		cryptoId, _ := primitive.ObjectIDFromHex(crypto.Id)
		return models.CryptoCurrency{
			Id:       cryptoId,
			Name:     crypto.Name,
			AssetId:  crypto.AssetId,
			PriceUsd: crypto.PriceUsd,
		}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	result, err := server.EditCrypto(ctx, &crypto)

	assert.Nil(t, err)
	assert.NotEmpty(t, result.Id)
	assert.Equal(t, crypto.Name, result.Name)
	assert.Equal(t, crypto.AssetId, result.AssetId)
	assert.Equal(t, crypto.PriceUsd, result.PriceUsd)
	assert.Empty(t, result.Votes)
	assert.NotEmpty(t, result.CreatedAt)
	assert.NotEmpty(t, result.UpdatedAt)

	defer cancel()
}

// Testing delete crypto with invalid id
func TestDeleteCryptoWithIdInvalid(t *testing.T) {
	server := AppServer{}
	crypto := returnMockProtoModelToDeleteCrypto()
	crypto.Id = "123abc"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	result, err := server.DeleteCrypo(ctx, &crypto)

	assert.Equal(t, "rpc error: code = InvalidArgument desc = id is invalid: 123abc err: the provided hex string is not a valid ObjectID", err.Error())
	assert.Empty(t, result.Message)

	defer cancel()
}

// Testing delete crypto with deletebyid did not find document
func TestDeleteCryptoWithDeleteByIdErrorErrNoDocuments(t *testing.T) {
	server := AppServer{}
	crypto := returnMockProtoModelToDeleteCrypto()

	mongodb.DeleteById = func(coll mongodb.IMCollection, id primitive.ObjectID) (primitive.ObjectID, error) {
		return id, mongo.ErrNoDocuments
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	result, err := server.DeleteCrypo(ctx, &crypto)

	assert.Equal(t, "rpc error: code = NotFound desc = mongo: no documents in result", err.Error())
	assert.Empty(t, result.Message)

	defer cancel()
}

// Testing delete crypto with deletebyid error
func TestDeleteCryptoWithDeleteByIdError(t *testing.T) {
	server := AppServer{}
	crypto := returnMockProtoModelToDeleteCrypto()

	mongodb.DeleteById = func(coll mongodb.IMCollection, id primitive.ObjectID) (primitive.ObjectID, error) {
		return id, errors.New("testing DeleteCrypo with error in DeleteById")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	result, err := server.DeleteCrypo(ctx, &crypto)

	assert.Equal(t, "rpc error: code = Internal desc = testing DeleteCrypo with error in DeleteById", err.Error())
	assert.Empty(t, result.Message)

	defer cancel()
}

// Testing delete crypto successful
func TestDeleteCryptoWithSuccess(t *testing.T) {
	server := AppServer{}
	crypto := returnMockProtoModelToDeleteCrypto()

	mongodb.DeleteById = func(coll mongodb.IMCollection, id primitive.ObjectID) (primitive.ObjectID, error) {
		return id, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	result, err := server.DeleteCrypo(ctx, &crypto)

	assert.Nil(t, err)
	assert.Equal(t, crypto.Id, result.Id)
	assert.Equal(t, "deleted successful", result.Message)

	defer cancel()
}

// Testing find crypto with invalid id
func TestFindCryptoWithIdInvalid(t *testing.T) {
	server := AppServer{}
	crypto := returnMockProtoModelToFindCrypto()
	crypto.Id = "123abc"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	_, err := server.FindCrypto(ctx, &crypto)

	assert.NotNil(t, err)
	assert.Equal(t, "rpc error: code = InvalidArgument desc = id is invalid: 123abc err: the provided hex string is not a valid ObjectID", err.Error())

	defer cancel()
}

// Testing find crypto with deletebyid did not find document
func TestFindCryptoWithDeleteByIdErrorErrNoDocuments(t *testing.T) {
	server := AppServer{}
	crypto := returnMockProtoModelToFindCrypto()

	mongodb.GetById = func(coll mongodb.IMCollection, id primitive.ObjectID) (models.CryptoCurrency, error) {
		return returnMockModelCryptoCurrency(), mongo.ErrNoDocuments
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	_, err := server.FindCrypto(ctx, &crypto)

	assert.NotNil(t, err)
	assert.Equal(t, "rpc error: code = NotFound desc = mongo: no documents in result", err.Error())

	defer cancel()
}

// Testing find crypto with deletebyid error
func TestFindCryptoWithDeleteByIdError(t *testing.T) {
	server := AppServer{}
	crypto := returnMockProtoModelToFindCrypto()

	mongodb.GetById = func(coll mongodb.IMCollection, id primitive.ObjectID) (models.CryptoCurrency, error) {
		return returnMockModelCryptoCurrency(), errors.New("testing FindCrypo with error in GetById")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	_, err := server.FindCrypto(ctx, &crypto)

	assert.NotNil(t, err)
	assert.Equal(t, "rpc error: code = Internal desc = testing FindCrypo with error in GetById", err.Error())

	defer cancel()
}

// Testing find crypto successful
func TestFindCryptoWithSuccess(t *testing.T) {
	server := AppServer{}
	crypto := returnMockProtoModelToFindCrypto()
	mockResponse := returnMockModelCryptoCurrency()
	crypto.Id = mockResponse.Id.Hex()

	mongodb.GetById = func(coll mongodb.IMCollection, id primitive.ObjectID) (models.CryptoCurrency, error) {
		return mockResponse, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	result, err := server.FindCrypto(ctx, &crypto)

	assert.Nil(t, err)
	assert.Equal(t, mockResponse.Id.Hex(), result.Id)
	assert.Equal(t, mockResponse.Name, result.Name)
	assert.Equal(t, mockResponse.AssetId, result.AssetId)
	assert.Equal(t, mockResponse.PriceUsd, result.PriceUsd)
	assert.Equal(t, mockResponse.Votes, result.Votes)
	assert.NotEmpty(t, result.CreatedAt)
	assert.NotEmpty(t, result.UpdatedAt)

	defer cancel()
}

// Testing list all cryptos with sort params invalid
func TestListAllCryptosWithSortParamsInvalid(t *testing.T) {
	server := AppServer{}
	sortParams := returnMockProtoModelToSortCryptos()
	sortParams.FieldSort = "a"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	_, err := server.ListAllCryptos(ctx, &sortParams)

	assert.NotNil(t, err)
	assert.Equal(t, "rpc error: code = InvalidArgument desc = field is invalid: a", err.Error())

	defer cancel()
}

// Testing list all cryptos with listall error
func TestListAllCryptosWithListAllError(t *testing.T) {
	server := AppServer{}
	sortParams := returnMockProtoModelToSortCryptos()

	mongodb.ListAll = func(coll mongodb.IMCollection, sort repositories.SortParams) (result []models.CryptoCurrency, err error) {
		return []models.CryptoCurrency{}, errors.New("testing ListAllCryptos with error in ListAll")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	_, err := server.ListAllCryptos(ctx, &sortParams)

	assert.NotNil(t, err)
	assert.Equal(t, "rpc error: code = Internal desc = testing ListAllCryptos with error in ListAll", err.Error())

	defer cancel()
}

// Testing list all cryptos with listall empty
func TestListAllCryptosWithListAllEmpty(t *testing.T) {
	server := AppServer{}
	sortParams := returnMockProtoModelToSortCryptos()
	mockCryptoEmpty := proto.ListCryptosResp{}
	mockCryptoEmpty.Crypto = []*proto.CryptoCurrency{}

	mongodb.ListAll = func(coll mongodb.IMCollection, sort repositories.SortParams) (result []models.CryptoCurrency, err error) {
		return []models.CryptoCurrency{}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	result, err := server.ListAllCryptos(ctx, &sortParams)

	assert.Nil(t, err)
	assert.Equal(t, mockCryptoEmpty.Crypto, result.Crypto)

	defer cancel()
}

// Testing list all cryptos successful
func TestListAllCryptosWithSuccess(t *testing.T) {
	server := AppServer{}
	sortParams := returnMockProtoModelToSortCryptos()

	mongodb.ListAll = func(coll mongodb.IMCollection, sort repositories.SortParams) (result []models.CryptoCurrency, err error) {
		return returnMockDbListAll(), nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	result, err := server.ListAllCryptos(ctx, &sortParams)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(result.Crypto))

	defer cancel()
}

// Testing upvote with invalid id
func TestUpvoteWithIdInvalid(t *testing.T) {
	server := AppServer{}
	crypto := returnMockProtoModelToVote()
	crypto.Id = "123abc"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	_, err := server.Upvote(ctx, &crypto)

	assert.NotNil(t, err)
	assert.Equal(t, "rpc error: code = InvalidArgument desc = id is invalid: 123abc err: the provided hex string is not a valid ObjectID", err.Error())

	defer cancel()
}

// Testing upvote with updatecrypto error
func TestUpvoteWithUpdateCryptoError(t *testing.T) {
	server := AppServer{}
	crypto := returnMockProtoModelToVote()

	mongodb.UpdateCrypto = func(coll mongodb.IMCollection, crypto models.CryptoCurrency) (models.CryptoCurrency, int64, error) {
		return models.CryptoCurrency{}, 0, errors.New("testing Upvote with error in UpdateCrypto")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	_, err := server.Upvote(ctx, &crypto)

	assert.NotNil(t, err)
	assert.Equal(t, "rpc error: code = Internal desc = testing Upvote with error in UpdateCrypto", err.Error())

	defer cancel()
}

// Testing upvote with updatecrypto matchedCount value is zero and getbyid error
func TestUpvoteWithUpdateCryptoMatchedCountZeroAndGetByIdError(t *testing.T) {
	server := AppServer{}
	crypto := returnMockProtoModelToVote()

	mongodb.UpdateCrypto = func(coll mongodb.IMCollection, crypto models.CryptoCurrency) (models.CryptoCurrency, int64, error) {
		return crypto, 0, nil
	}

	mongodb.GetById = func(coll mongodb.IMCollection, id primitive.ObjectID) (crypto models.CryptoCurrency, err error) {
		return models.CryptoCurrency{}, errors.New("testing Upvote with error in GetById")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	_, err := server.Upvote(ctx, &crypto)

	assert.NotNil(t, err)
	assert.Equal(t, "rpc error: code = NotFound desc = testing Upvote with error in GetById", err.Error())

	defer cancel()
}

// Testing upvote with updatecrypto matchedCount value is zero and getbyid empty
func TestUpvoteWithUpdateCryptoMatchedCountZeroButGetByIdCryptoNotExist(t *testing.T) {
	server := AppServer{}
	crypto := returnMockProtoModelToVote()

	mongodb.UpdateCrypto = func(coll mongodb.IMCollection, crypto models.CryptoCurrency) (models.CryptoCurrency, int64, error) {
		return crypto, 0, nil
	}

	mongodb.GetById = func(coll mongodb.IMCollection, id primitive.ObjectID) (crypto models.CryptoCurrency, err error) {
		return models.CryptoCurrency{}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	_, err := server.Upvote(ctx, &crypto)

	assert.NotNil(t, err)
	assert.Equal(t, "rpc error: code = Internal desc = crypto not exist", err.Error())

	defer cancel()
}

// Testing upvote successful
func TestUpvoteWithSuccess(t *testing.T) {
	server := AppServer{}
	crypto := returnMockProtoModelToVote()

	mongodb.UpdateCrypto = func(coll mongodb.IMCollection, crypto models.CryptoCurrency) (models.CryptoCurrency, int64, error) {
		return crypto, 1, nil
	}

	mongodb.GetById = func(coll mongodb.IMCollection, id primitive.ObjectID) (crypto models.CryptoCurrency, err error) {
		return models.CryptoCurrency{
			Votes: 1,
		}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	result, err := server.Upvote(ctx, &crypto)

	assert.Nil(t, err)
	assert.Equal(t, crypto.Id, result.Id)
	assert.Equal(t, "registered upvote successful", result.Message)

	defer cancel()
}

// Testing downvote with invalid id
func TestDownvoteWithIdInvalid(t *testing.T) {
	server := AppServer{}
	crypto := returnMockProtoModelToVote()
	crypto.Id = "123abc"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	_, err := server.Downvote(ctx, &crypto)

	assert.NotNil(t, err)
	assert.Equal(t, "rpc error: code = InvalidArgument desc = id is invalid: 123abc err: the provided hex string is not a valid ObjectID", err.Error())

	defer cancel()
}

// Testing downvote with updatecrypto error
func TestDownvoteWithUpdateCryptoError(t *testing.T) {
	server := AppServer{}
	crypto := returnMockProtoModelToVote()

	mongodb.UpdateCrypto = func(coll mongodb.IMCollection, crypto models.CryptoCurrency) (models.CryptoCurrency, int64, error) {
		return models.CryptoCurrency{}, 0, errors.New("testing Downvote with error in UpdateCrypto")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	_, err := server.Downvote(ctx, &crypto)

	assert.NotNil(t, err)
	assert.Equal(t, "rpc error: code = Internal desc = testing Downvote with error in UpdateCrypto", err.Error())

	defer cancel()
}

// Testing downvote with updatecrypto matchedCount value is zero and getbyid error
func TestDownvoteWithUpdateCryptoMatchedCountZeroAndGetByIdError(t *testing.T) {
	server := AppServer{}
	crypto := returnMockProtoModelToVote()

	mongodb.UpdateCrypto = func(coll mongodb.IMCollection, crypto models.CryptoCurrency) (models.CryptoCurrency, int64, error) {
		return crypto, 0, nil
	}

	mongodb.GetById = func(coll mongodb.IMCollection, id primitive.ObjectID) (crypto models.CryptoCurrency, err error) {
		return models.CryptoCurrency{}, errors.New("testing Downvote with error in GetById")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	_, err := server.Downvote(ctx, &crypto)

	assert.NotNil(t, err)
	assert.Equal(t, "rpc error: code = NotFound desc = testing Downvote with error in GetById", err.Error())

	defer cancel()
}

// Testing downvote with vote equal zero
func TestDownvoteWithVoteEqualZero(t *testing.T) {
	server := AppServer{}
	crypto := returnMockProtoModelToVote()

	mongodb.UpdateCrypto = func(coll mongodb.IMCollection, crypto models.CryptoCurrency) (models.CryptoCurrency, int64, error) {
		return crypto, 0, nil
	}

	mongodb.GetById = func(coll mongodb.IMCollection, id primitive.ObjectID) (crypto models.CryptoCurrency, err error) {
		return models.CryptoCurrency{}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	_, err := server.Downvote(ctx, &crypto)

	assert.NotNil(t, err)
	assert.Equal(t, "rpc error: code = Internal desc = unchanged crypto because vote = 0", err.Error())

	defer cancel()
}

// Testing upvote successful
func TestDownvoteWithSuccess(t *testing.T) {
	server := AppServer{}
	crypto := returnMockProtoModelToVote()

	mongodb.UpdateCrypto = func(coll mongodb.IMCollection, crypto models.CryptoCurrency) (models.CryptoCurrency, int64, error) {
		return crypto, 1, nil
	}

	mongodb.GetById = func(coll mongodb.IMCollection, id primitive.ObjectID) (crypto models.CryptoCurrency, err error) {
		return models.CryptoCurrency{
			Votes: 0,
		}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	result, err := server.Downvote(ctx, &crypto)

	assert.Nil(t, err)
	assert.Equal(t, crypto.Id, result.Id)
	assert.Equal(t, "registered downvote successful", result.Message)

	defer cancel()
}

// Testing monitor votes with invalid id
func TestMonitorVotesWithIdInvalid(t *testing.T) {
	server := AppServer{}
	crypto := returnMockProtoModelToMonitorVotes()
	crypto.Id = "123abc"
	mockStream := Mock_EndPointCryptos_MonitorVotesServer{}

	_, cancel := context.WithTimeout(context.Background(), time.Second*5)
	err := server.MonitorVotes(&crypto, &mockStream)

	assert.NotNil(t, err)
	assert.Equal(t, "rpc error: code = InvalidArgument desc = id is invalid: 123abc err: the provided hex string is not a valid ObjectID", err.Error())

	defer cancel()
}

// Help function to TestMonitorVotesWithGetByIdError and TestMonitorVotesWithSuccess
func mockUpdateToStream(server AppServer) proto.VoteReq {
	cryptoToUpVote := returnMockProtoModelToVote()

	mongodb.UpdateCrypto = func(coll mongodb.IMCollection, crypto models.CryptoCurrency) (models.CryptoCurrency, int64, error) {
		return crypto, 1, nil
	}

	mongodb.GetById = func(coll mongodb.IMCollection, id primitive.ObjectID) (crypto models.CryptoCurrency, err error) {
		return models.CryptoCurrency{}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	_, _ = server.Upvote(ctx, &cryptoToUpVote)

	defer cancel()

	return cryptoToUpVote
}

// Testing monitor votes with getbyid error
func TestMonitorVotesWithGetByIdError(t *testing.T) {
	server := AppServer{}

	cryptoMonitor := returnMockProtoModelToMonitorVotes()
	mockStream := Mock_EndPointCryptos_MonitorVotesServer{}

	StartChanToStream()
	SetObserver = func(id string) {
		observer <- cryptoMonitor.Id
	}

	cryptoReceivedUpVote := mockUpdateToStream(server)
	cryptoMonitor.Id = cryptoReceivedUpVote.Id

	mongodb.GetById = func(coll mongodb.IMCollection, id primitive.ObjectID) (models.CryptoCurrency, error) {
		return models.CryptoCurrency{}, errors.New("testing MonitorVotes with error in GetById")
	}

	_, cancel := context.WithTimeout(context.Background(), time.Second*5)
	err := server.MonitorVotes(&cryptoMonitor, &mockStream)

	defer cancel()

	assert.NotNil(t, err)
	assert.Equal(t, "rpc error: code = Internal desc = testing MonitorVotes with error in GetById", err.Error())
}

// Testing monitor votes successful
func TestMonitorVotesWithSuccess(t *testing.T) {
	server := AppServer{}

	cryptoResponseStream := returnMockModelCryptoCurrency()
	cryptoMonitorSendToMonitor := returnMockProtoModelToMonitorVotes()
	mockStream := Mock_EndPointCryptos_MonitorVotesServer{}

	StartChanToStream()
	SetObserver = func(id string) {
		observer <- cryptoMonitorSendToMonitor.Id
	}

	cryptoReceivedUpVote := mockUpdateToStream(server)
	cryptoMonitorSendToMonitor.Id = cryptoReceivedUpVote.Id
	objId, _ := primitive.ObjectIDFromHex(cryptoReceivedUpVote.Id)
	cryptoResponseStream.Id = objId
	cryptoResponseStream.Votes += 1

	mongodb.GetById = func(coll mongodb.IMCollection, id primitive.ObjectID) (models.CryptoCurrency, error) {
		return cryptoResponseStream, nil
	}

	// Set timeout to testing finish
	timeout := time.After(3 * time.Second)
	go func() {
		_, cancel := context.WithTimeout(context.Background(), time.Second*5)
		server.MonitorVotes(&cryptoMonitorSendToMonitor, &mockStream)
		defer cancel()
	}()

	select {
	case <-timeout:
		logger.Info("TESTING", "Timeout successful ")
	}

	assert.Equal(t, 1, len(mockStream.Results))
	assert.Equal(t, cryptoResponseStream.Id.Hex(), mockStream.Results[0].Id)
	assert.Equal(t, cryptoResponseStream.Name, mockStream.Results[0].Name)
}
