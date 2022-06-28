package controllers

import (
	"api-desafio-kvr/helpers"
	"api-desafio-kvr/models"
	"api-desafio-kvr/proto"
	"api-desafio-kvr/repositories"
	db "api-desafio-kvr/repositories/mongodb"
	rds "api-desafio-kvr/repositories/redis"
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/grpc/status"
)

// Codes of return google.golang.org/grpc/codes

var logger = &helpers.Log{}
var observer chan string

type AppServer struct {
	proto.UnimplementedEndPointCryptosServer
	Database *mongo.Collection
}

func StartChanToStream() {
	logger.Info("", "Starting channel for stream")
	observer = make(chan string)
}

var SetObserver = func(id string) {
	observer <- id
}

func (a *AppServer) CreateCrypto(ctx context.Context, req *proto.CreateCryptoReq) (*proto.CryptoCurrency, error) {
	logger.Debug("", "Creating crypto received params "+req.String())
	cryptoResponse := proto.CryptoCurrency{}

	err := helpers.ValidatorInCreateCrypto(req)
	if err != nil {
		logger.Error("", "Params create crypto is invalid "+req.String())
		return &cryptoResponse, status.Errorf(3, err.Error())
	}

	cryptoDb := models.CryptoCurrency{
		Name:     cases.Title(language.AmericanEnglish).String(req.GetName()),
		AssetId:  cases.Upper(language.AmericanEnglish).String(req.GetAssetId()),
		PriceUsd: req.GetPriceUsd(),
	}

	insertedCrypto, err := db.InsertCryptos(a.Database, cryptoDb)
	if err != nil {
		logger.Error("", "Crypto not created "+req.String()+" error: "+err.Error())
		return &cryptoResponse, status.Errorf(13, err.Error())
	}

	// Set cache in Redis
	err = rds.Set(insertedCrypto.Id.Hex(), insertedCrypto, rds.YesDeleteAll)
	if err != nil {
		logger.Error(insertedCrypto.Id.Hex(), "Error to set cache in redis: "+err.Error())
	}

	byteCrypto, err := json.Marshal(insertedCrypto)
	if err != nil {
		logger.Error(insertedCrypto.Id.Hex(), "Error in response: "+err.Error())
		return &cryptoResponse, status.Errorf(13, err.Error())
	}
	err = json.Unmarshal(byteCrypto, &cryptoResponse)
	if err != nil {
		logger.Error(insertedCrypto.Id.Hex(), "Error in response: "+err.Error())
		return &cryptoResponse, status.Errorf(13, err.Error())
	}

	cryptoResponse.CreatedAt = insertedCrypto.CreatedAt.Format("2006-01-02T15:04:05.999Z")
	cryptoResponse.UpdatedAt = insertedCrypto.UpdatedAt.Format("2006-01-02T15:04:05.999Z")

	logger.Info(cryptoResponse.Id, "Crypto created successful")
	return &cryptoResponse, nil
}

func (a *AppServer) EditCrypto(ctx context.Context, req *proto.EditCryptoReq) (*proto.CryptoCurrency, error) {
	logger.Debug("", "Editing crypto received params "+req.String())
	cryptoResponse := proto.CryptoCurrency{}

	err := helpers.ValidatorInEditCrypto(req)
	if err != nil {
		logger.Error("", "Params edit crypto is invalid "+req.String())
		return &cryptoResponse, status.Errorf(3, err.Error())
	}

	objId, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		logger.Error(req.GetId(), "Params edit crypto is invalid "+req.String())
		return &cryptoResponse, status.Errorf(3, err.Error())
	}

	cryptoUpdate := models.CryptoCurrency{
		Id:         objId,
		Name:       cases.Title(language.AmericanEnglish).String(req.GetName()),
		AssetId:    cases.Upper(language.AmericanEnglish).String(req.GetAssetId()),
		PriceUsd:   req.GetPriceUsd(),
		UpdateType: models.UpdateOnly,
	}

	updatedCrypto, _, err := db.UpdateCrypto(a.Database, cryptoUpdate)
	if err != nil {
		logger.Error("", "Crypto not edited "+req.String()+" error: "+err.Error())
		return &cryptoResponse, status.Errorf(13, err.Error())
	}

	crypto, err := db.GetById(a.Database, updatedCrypto.Id)
	if err != nil {
		logger.Error("", "Crypto not find after update "+req.String()+" error: "+err.Error())
		return &cryptoResponse, status.Errorf(5, err.Error())
	}

	// Set cache in Redis
	err = rds.Set(crypto.Id.Hex(), crypto, rds.YesDeleteAll)
	if err != nil {
		logger.Error(req.GetId(), "Error to set cache in redis: "+err.Error())
	}

	byteCrypto, err := json.Marshal(crypto)
	if err != nil {
		logger.Error(crypto.Id.Hex(), "Error in response EditCrypto: "+err.Error())
		return &cryptoResponse, status.Errorf(13, err.Error())
	}
	err = json.Unmarshal(byteCrypto, &cryptoResponse)
	if err != nil {
		logger.Error(crypto.Id.Hex(), "Error in response EditCrypto: "+err.Error())
		return &cryptoResponse, status.Errorf(13, err.Error())
	}

	logger.Info(cryptoResponse.Id, "Crypto updated successful")

	go SetObserver(req.GetId())
	return &cryptoResponse, nil
}

func (a *AppServer) DeleteCrypo(ctx context.Context, req *proto.DeleteCryptoReq) (*proto.DefaultResp, error) {
	logger.Debug("", "Deleting crypto received params "+req.String())
	messageResponse := proto.DefaultResp{}

	err := helpers.IdValidator(req.GetId())
	if err != nil {
		logger.Error("", "Params to delete crypto is invalid "+req.String())
		return &messageResponse, status.Errorf(3, err.Error())
	}

	objId, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		logger.Error(req.GetId(), "Params edit crypto is invalid "+req.String())
		return &messageResponse, status.Errorf(3, err.Error())
	}

	_, err = db.DeleteById(a.Database, objId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Error(req.GetId(), "Delete crypto error: "+err.Error())
			return &messageResponse, status.Errorf(5, err.Error())
		}
		logger.Error("", "Crypto not deleted "+req.String()+" error: "+err.Error())
		return &messageResponse, status.Errorf(13, err.Error())
	}

	// Delete cache in Redis
	err = rds.Del(req.GetId())
	if err != nil {
		logger.Error(req.GetId(), "Error to delete cache in redis: "+err.Error())
	}

	messageResponse.Id = req.GetId()
	messageResponse.Message = "deleted successful"

	logger.Info(req.GetId(), "Crypto deleted successful")

	go SetObserver(req.GetId())
	return &messageResponse, err
}

func (a *AppServer) FindCrypto(ctx context.Context, req *proto.FindCryptoReq) (*proto.CryptoCurrency, error) {
	logger.Debug("", "Finding crypto received params "+req.String())
	cryptoResponse := proto.CryptoCurrency{}

	err := helpers.IdValidator(req.GetId())
	if err != nil {
		logger.Error("", "Params to find crypto is invalid "+req.String())
		return &cryptoResponse, status.Errorf(3, err.Error())
	}

	// Get cache in Redis
	cache := rds.Get(req.GetId())
	if cache != nil {
		err = json.Unmarshal([]byte(cache), &cryptoResponse)
		if err == nil {
			logger.Info(req.GetId(), "Found cache successful")
			return &cryptoResponse, nil
		}
		// else continue
		logger.Warn(req.GetId(), err.Error())
	}

	objId, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return &cryptoResponse, nil
	}

	findResp, err := db.GetById(a.Database, objId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Error(req.GetId(), "Find crypto error: "+err.Error())
			return &cryptoResponse, status.Errorf(5, err.Error())
		}
		logger.Error("", "Crypto not found because error "+req.String()+" error: "+err.Error())
		return &cryptoResponse, status.Errorf(13, err.Error())
	}

	// Set cache in Redis
	err = rds.Set(findResp.Id.Hex(), findResp, rds.YesDeleteAll)
	if err != nil {
		logger.Error(findResp.Id.Hex(), "Error to set cache in redis: "+err.Error())
	}

	byteFind, err := json.Marshal(findResp)
	if err != nil {
		logger.Error(findResp.Id.Hex(), "Error in response FindCrypto: "+err.Error())
		return &cryptoResponse, status.Errorf(13, err.Error())
	}
	err = json.Unmarshal(byteFind, &cryptoResponse)
	if err != nil {
		logger.Error(findResp.Id.Hex(), "Error in response FindCrypto: "+err.Error())
		return &cryptoResponse, status.Errorf(13, err.Error())
	}

	logger.Info(req.GetId(), "Crypto found successful")
	return &cryptoResponse, nil
}

func (a *AppServer) ListAllCryptos(ctx context.Context, req *proto.SortCryptosReq) (*proto.ListCryptosResp, error) {
	logger.Debug("", "Listing crypto received params "+req.String())
	cryptoListResponse := proto.ListCryptosResp{}

	err := helpers.ValidatorListAllCryptos(req)
	if err != nil {
		logger.Error("", "Params to list crypto is invalid "+req.String())
		return &cryptoListResponse, status.Errorf(3, err.Error())
	}

	// Get cache in Redis
	key := rds.PrefixDeleteAll + "-" + req.GetFieldSort() + "-" + strconv.FormatBool(req.GetOrderBy())
	cryptos := rds.Get(key)
	if cryptos != nil {
		err = json.Unmarshal([]byte(cryptos), &cryptoListResponse)
		if err == nil {
			logger.Info(key, "Found cache successful")
			return &cryptoListResponse, nil
		}
		// else continue
		logger.Warn(key, err.Error())
	}

	// if GetOrderBy == true then orderBy is ASC, else orderBy is DESC
	sort := repositories.SortParams{
		Field: req.GetFieldSort(),
		Asc:   req.GetOrderBy(),
	}

	response, err := db.ListAll(a.Database, sort)
	if err != nil {
		logger.Error("", "Cryptos not listed because error "+req.String()+" error: "+err.Error())
		return &cryptoListResponse, status.Errorf(13, err.Error())
	}

	cryptoList := []*proto.CryptoCurrency{}
	for _, value := range response {
		proto := value.ToProtoCrypto()
		cryptoList = append(cryptoList, &proto)
	}

	amount := len(cryptoList)
	cryptoListResponse.Crypto = cryptoList

	byteCrypto, err := json.Marshal(cryptoListResponse)
	if err != nil {
		logger.Error("", "Error to set cache in redis: "+err.Error())
		return &cryptoListResponse, nil
	}
	// Set cache in Redis
	key = rds.PrefixDeleteAll + "-" + sort.Field + "-" + strconv.FormatBool(sort.Asc)
	err = rds.SetByByte(key, string(byteCrypto), rds.NoDeleteAll)
	if err != nil {
		logger.Error("", "Error to set cache in redis: "+err.Error())
	}

	logger.Info("", "Listed "+strconv.Itoa(amount)+" crypto successful")
	return &cryptoListResponse, nil
}

func (a *AppServer) Upvote(ctx context.Context, req *proto.VoteReq) (*proto.DefaultResp, error) {
	logger.Debug(req.GetId(), "Upvoting crypto received params "+req.String())
	responseMessage := proto.DefaultResp{}

	err := helpers.IdValidator(req.GetId())
	if err != nil {
		logger.Error(req.GetId(), "Params to upvote crypto is invalid "+req.String())
		return &responseMessage, status.Errorf(3, err.Error())
	}

	objId, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return &responseMessage, status.Errorf(3, err.Error())
	}

	crypto := models.CryptoCurrency{
		Id:         objId,
		UpdateType: models.UpVote,
	}

	_, matchedCount, err := db.UpdateCrypto(a.Database, crypto)
	if err != nil {
		logger.Error(req.GetId(), "Crypto upvote error: "+err.Error())
		return &responseMessage, status.Errorf(13, err.Error())
	}

	crypto, errGet := db.GetById(a.Database, crypto.Id)

	// if document not updated, verify if exist
	if matchedCount == 0 {
		// Verifying reason to mathedCount == 0 maybe crypto not exists
		if errGet != nil {
			logger.Error(req.GetId(), "Upvote error: "+errGet.Error())
			return &responseMessage, status.Errorf(5, errGet.Error())
		}

		err = errors.New("crypto not exist")
		logger.Error(req.GetId(), err.Error())
		return &responseMessage, status.Errorf(13, err.Error())
	}

	responseMessage.Id = req.GetId()
	responseMessage.Message = "registered upvote successful"

	logger.Info(req.GetId(), "Crypto upvote successful")

	// Set cache in Redis
	err = rds.Set(crypto.Id.Hex(), crypto, rds.YesDeleteAll)
	if err != nil {
		logger.Error(req.GetId(), "Error to set cache in redis: "+err.Error())
	}

	go SetObserver(req.GetId())
	return &responseMessage, nil
}

func (a *AppServer) Downvote(ctx context.Context, req *proto.VoteReq) (*proto.DefaultResp, error) {
	logger.Debug("", "Downvoting crypto received params "+req.String())
	responseMessage := proto.DefaultResp{}
	responseMessage.Id = req.GetId()

	err := helpers.IdValidator(req.GetId())
	if err != nil {
		logger.Error(req.GetId(), "Params to downvote crypto is invalid")
		return &responseMessage, status.Errorf(3, err.Error())
	}

	objId, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return &responseMessage, status.Errorf(3, err.Error())
	}

	crypto := models.CryptoCurrency{
		Id:         objId,
		UpdateType: models.DownVote,
	}

	_, matchedCount, err := db.UpdateCrypto(a.Database, crypto)
	if err != nil {
		logger.Error(req.GetId(), "Crypto downvote error: "+err.Error())
		return &responseMessage, status.Errorf(13, err.Error())
	}

	crypto, errGet := db.GetById(a.Database, crypto.Id)

	if matchedCount == 0 {
		// Verifying reason to mathedCount == 0 maybe crypto not exists
		if errGet != nil {
			logger.Error(req.GetId(), "Downvote error: "+errGet.Error())
			return &responseMessage, status.Errorf(5, errGet.Error())
		}
		// vote is 0
		err = errors.New("unchanged crypto because vote = 0")
		logger.Error(req.GetId(), err.Error())
		return &responseMessage, status.Errorf(13, err.Error())
	}

	responseMessage.Message = "registered downvote successful"
	logger.Info(req.GetId(), "Crypto downvote successful")

	// Set cache in Redis
	err = rds.Set(crypto.Id.Hex(), crypto, rds.YesDeleteAll)
	if err != nil {
		logger.Error(req.GetId(), "Error to set cache in redis: "+err.Error())
	}

	go SetObserver(req.GetId())
	return &responseMessage, nil
}

func (a *AppServer) MonitorVotes(req *proto.MonitorVotesReq, stream proto.EndPointCryptos_MonitorVotesServer) error {
	err := helpers.IdValidator(req.GetId())
	if err != nil {
		logger.Error("", "Params to stream crypto is invalid "+req.String())
		return status.Errorf(3, err.Error())
	}

	logger.Info(req.GetId(), "Streaming crypto...")

	for {
		// received id to observer chan
		cryptoUpdatedId := <-observer

		if cryptoUpdatedId == req.GetId() {

			objId, err := primitive.ObjectIDFromHex(cryptoUpdatedId)
			if err != nil {
				return err
			}

			cryptoFound, err := db.GetById(a.Database, objId)
			if err != nil {
				logger.Error(req.GetId(), "Error to stram crypto: "+err.Error())
				return err
			}

			streamCrypto := &proto.CryptoCurrency{
				Id:        cryptoFound.Id.Hex(),
				Name:      cryptoFound.Name,
				AssetId:   cryptoFound.AssetId,
				PriceUsd:  cryptoFound.PriceUsd,
				Votes:     cryptoFound.Votes,
				CreatedAt: cryptoFound.CreatedAt.Format("2006-01-02T15:04:05.999Z"),
				UpdatedAt: cryptoFound.UpdatedAt.Format("2006-01-02T15:04:05.999Z"),
			}
			err = stream.Send(streamCrypto)

			if err != nil {
				errStatus := status.Convert(err)
				if cases.Lower(language.AmericanEnglish).String(errStatus.Message()) == "transport is closing" {
					logger.Warn(req.GetId(), "Stream "+errStatus.Message())
					return nil
				}
				return err
			}

			out, err := json.Marshal(streamCrypto)
			if err != nil {
				logger.Error(req.GetId(), err.Error())
			}
			logger.Info(req.GetId(), "Streaming in Crypto "+string(out))
		}
	}
}
