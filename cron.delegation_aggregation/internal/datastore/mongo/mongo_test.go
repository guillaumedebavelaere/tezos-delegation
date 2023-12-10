package mongo_test

import (
	"context"
	"github.com/guillaumedebavelaere/tezos-delegation/pkg/test/docker"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	gomongo "go.mongodb.org/mongo-driver/mongo"

	mongosvc "github.com/guillaumedebavelaere/tezos-delegation/cron.delegation_aggregation/internal/datastore/mongo"
	"github.com/guillaumedebavelaere/tezos-delegation/pkg/mongo"
)

type MongoTestSuite struct {
	suite.Suite
	dockerTest  *docker.Docker
	mongoClient mongo.Client
	mongoSvc    *mongosvc.Datastore
	database    *gomongo.Database
	collection  *gomongo.Collection
}

func (suite *MongoTestSuite) SetupSuite() {
	var err error

	suite.dockerTest, err = docker.New(docker.NewMongoContainer())
	suite.Require().NoError(err)
	suite.Require().NoError(suite.dockerTest.Start())

	suite.mongoClient = mongo.New(&mongo.Config{
		URI:               docker.EndpointMongo(suite.dockerTest.GetResources()[0]),
		ConnectTimeout:    30 * time.Second,
		HeartbeatInterval: 30 * time.Second,
		Timeout:           30 * time.Second,
		Username:          "admin",
		Password:          "admin",
	})
	suite.Require().NoError(suite.mongoClient.Init())

	suite.mongoSvc = mongosvc.New(suite.mongoClient)
	suite.Require().NoError(suite.mongoSvc.Init())
}

func (suite *MongoTestSuite) SetupTest() {
	suite.database = suite.mongoClient.C().Database("tezos-delegations")
	suite.collection = suite.database.Collection("delegations")
}

func (suite *MongoTestSuite) TearDownTest() {
	suite.Require().NoError(suite.database.Drop(context.Background()))
}

func (suite *MongoTestSuite) TearDownSuite() {
	suite.Require().NoError(suite.dockerTest.Stop())
}

func TestMongoTestSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(MongoTestSuite))
}
