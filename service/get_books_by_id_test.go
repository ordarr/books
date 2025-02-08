package service

import (
	pb "github.com/ordarr/books/v1"
	"github.com/ordarr/data/core"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *BookTestSuite) TestGetBookById() {
	t := suite.T()

	suite.Run("ReturnsPopulatedBook", func() {
		suite.mockRepo.(*MockRepo).On("GetByID", []string{"12345"}).Return([]*core.Book{
			{
				BaseTable: core.BaseTable{ID: "12345"},
				Name:      "Name One",
			},
		}, nil)

		out, _ := suite.client.GetBooks(suite.ctx, &pb.GetBooksRequest{Ids: []string{"12345"}})

		assert.NotNil(t, out)
		assert.Equal(t, "Name One", out.Content[0].Name)
	})

	suite.Run("ErrorWhenBookDoesntExist", func() {
		suite.mockRepo.(*MockRepo).On("GetByID", []string{"12345"}).Return(nil, status.Error(codes.NotFound, "book not found"))
		_, err := suite.client.GetBooks(suite.ctx, &pb.GetBooksRequest{Ids: []string{"12345"}})

		assert.NotNil(t, err)
		assert.ErrorIs(t, err, status.Error(codes.NotFound, "book not found"))
	})
}
