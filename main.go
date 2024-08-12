package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
)

const (
	defaultMongoUrl         = "mongodb://arex:iLoveArex@arex-helm-name-beta-arex-mongodb.arex.svc.cluster.local:27017/arex_storage_db"
	defaultScheduleEndpoint = "http://arex-helm-name-beta-arex-schedule.arex.svc.cluster.local:8080/api/createPlan"
)

func main() {
	appId := os.Getenv("APP_ID")
	targetHost := os.Getenv("TARGET_HOST")

	if appId == "" || targetHost == "" {
		log.Fatal("APP_ID and TARGET_HOST environment variables must be set")
	}

	mongoUrl := os.Getenv("MONGO_URL")
	if mongoUrl == "" {
		mongoUrl = defaultMongoUrl
	}

	scheduleEndpoint := os.Getenv("SCHEDULE_ENDPOINT")
	if scheduleEndpoint == "" {
		scheduleEndpoint = defaultScheduleEndpoint
	}

	clientOptions := options.Client().ApplyURI(mongoUrl)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	db := client.Database("arex_storage_db")
	fmt.Println("Connected successfully to server")

	// query auto pin cases
	pinnedCasesCursor, err := db.Collection("PinnedServletMocker").Find(context.TODO(), bson.M{"appId": appId})
	if err != nil {
		log.Fatal(err)
	}
	defer pinnedCasesCursor.Close(context.TODO())

	var pinnedCases []struct {
		Id            string `bson:"_id"`
		OperationName string `bson:"operationName"`
	}
	if err = pinnedCasesCursor.All(context.TODO(), &pinnedCases); err != nil {
		log.Fatal(err)
	}

	caseGroupByOp := make(map[string][]string)
	for _, item := range pinnedCases {
		caseGroupByOp[item.OperationName] = append(caseGroupByOp[item.OperationName], item.Id)
	}

	// query service operations
	operationNames := make([]string, 0, len(caseGroupByOp))
	for name := range caseGroupByOp {
		operationNames = append(operationNames, name)
	}

	operationsCursor, err := db.Collection("ServiceOperation").Find(context.TODO(), bson.M{
		"appId":         appId,
		"operationName": bson.M{"$in": operationNames},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer operationsCursor.Close(context.TODO())

	var operations []struct {
		Id            string `bson:"_id"`
		OperationName string `bson:"operationName"`
	}
	if err = operationsCursor.All(context.TODO(), &operations); err != nil {
		log.Fatal(err)
	}

	createPlanReq := struct {
		AppId                 string `json:"appId"`
		ReplayPlanType        int    `json:"replayPlanType"`
		TargetEnv             string `json:"targetEnv"`
		Operator              string `json:"operator"`
		OperationCaseInfoList []struct {
			OperationId  string   `json:"operationId"`
			ReplayIdList []string `json:"replayIdList"`
		} `json:"operationCaseInfoList"`
	}{
		AppId:          appId,
		ReplayPlanType: 2,
		TargetEnv:      targetHost,
		Operator:       "AREX",
		OperationCaseInfoList: make([]struct {
			OperationId  string   `json:"operationId"`
			ReplayIdList []string `json:"replayIdList"`
		}, len(operations)),
	}

	for i, op := range operations {
		createPlanReq.OperationCaseInfoList[i].OperationId = op.Id
		createPlanReq.OperationCaseInfoList[i].ReplayIdList = caseGroupByOp[op.OperationName]
	}

	reqBody, err := json.MarshalIndent(createPlanReq, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(reqBody))

	resp, err := http.Post(scheduleEndpoint, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.StatusCode)
}
