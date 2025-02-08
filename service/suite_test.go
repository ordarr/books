package service

import (
	"context"
	pb "github.com/ordarr/books/v1"
	"github.com/ordarr/data/core"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type BookTestSuite struct {
	suite.Suite
	ctx      context.Context
	client   pb.BooksClient
	mockRepo core.IBookRepository
	closer   func()
}

func (suite *BookTestSuite) SetupSubTest() {
	suite.ctx = context.Background()
	suite.mockRepo = &MockRepo{}
	suite.client, suite.closer = CreateClient(suite.mockRepo)
}

func (suite *BookTestSuite) TearDownSubTest() {
	_ = os.Remove("ordarr.db")
	suite.closer()
}

func TestBookTestSuite(t *testing.T) {
	suite.Run(t, new(BookTestSuite))
}
