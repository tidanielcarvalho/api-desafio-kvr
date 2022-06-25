package controllers

import (
	"api-desafio-kvr/helpers"
	"api-desafio-kvr/models"
	"api-desafio-kvr/proto"
	"api-desafio-kvr/repositories"
	db "api-desafio-kvr/repositories/mongodb"
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

	byteCrypto, _ := json.Marshal(insertedCrypto)
	json.Unmarshal(byteCrypto, &cryptoResponse)

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

	objId, _ := primitive.ObjectIDFromHex(req.GetId())
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

	byteCrypto, _ := json.Marshal(crypto)
	json.Unmarshal(byteCrypto, &cryptoResponse)

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

	objId, _ := primitive.ObjectIDFromHex(req.GetId())
	_, err = db.DeleteById(a.Database, objId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Error(req.GetId(), "Delete crypto error: "+err.Error())
			return &messageResponse, status.Errorf(5, err.Error())
		}
		logger.Error("", "Crypto not deleted "+req.String()+" error: "+err.Error())
		return &messageResponse, status.Errorf(13, err.Error())
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

	objId, _ := primitive.ObjectIDFromHex(req.GetId())
	findResp, err := db.GetById(a.Database, objId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Error(req.GetId(), "Find crypto error: "+err.Error())
			return &cryptoResponse, status.Errorf(5, err.Error())
		}
		logger.Error("", "Crypto not found because error "+req.String()+" error: "+err.Error())
		return &cryptoResponse, status.Errorf(13, err.Error())
	}

	byteFind, _ := json.Marshal(findResp)
	json.Unmarshal(byteFind, &cryptoResponse)

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
		cryptoList = append(cryptoList, &proto.CryptoCurrency{
			Id:        value.Id.Hex(),
			Name:      value.Name,
			AssetId:   value.AssetId,
			PriceUsd:  value.PriceUsd,
			Votes:     int32(value.Votes),
			CreatedAt: value.CreatedAt.Format("2006-01-02T15:04:05.999Z"),
			UpdatedAt: value.UpdatedAt.Format("2006-01-02T15:04:05.999Z"),
		})
	}

	amount := len(cryptoList)
	cryptoListResponse.Crypto = cryptoList

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

	objId, _ := primitive.ObjectIDFromHex(req.GetId())
	crypto := models.CryptoCurrency{
		Id:         objId,
		UpdateType: models.UpVote,
	}

	_, matchedCount, err := db.UpdateCrypto(a.Database, crypto)
	if err != nil {
		logger.Error(req.GetId(), "Crypto upvote error: "+err.Error())
		return &responseMessage, status.Errorf(13, err.Error())
	}

	// if document not updated, verify if exist
	if matchedCount == 0 {
		_, errGet := db.GetById(a.Database, crypto.Id)
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

	objId, _ := primitive.ObjectIDFromHex(req.GetId())
	crypto := models.CryptoCurrency{
		Id:         objId,
		UpdateType: models.DownVote,
	}

	_, matchedCount, err := db.UpdateCrypto(a.Database, crypto)
	if err != nil {
		logger.Error(req.GetId(), "Crypto downvote error: "+err.Error())
		return &responseMessage, status.Errorf(13, err.Error())
	}

	if matchedCount == 0 {
		_, errGet := db.GetById(a.Database, crypto.Id)
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

			objId, _ := primitive.ObjectIDFromHex(cryptoUpdatedId)
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

			out, _ := json.Marshal(streamCrypto)
			logger.Info(req.GetId(), "Streaming in Crypto "+string(out))
		}
	}
}
