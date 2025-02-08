package service

import (
	pb "github.com/ordarr/books/v1"
	"github.com/ordarr/data/core"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *BookTestSuite) TestGetBookByTitle() {
	t := suite.T()

	suite.Run("ReturnsPopulatedBook", func() {
		suite.mockRepo.(*MockRepo).On("GetByName", []string{"Name One"}).Return([]*core.Book{
			{
				BaseTable: core.BaseTable{ID: "12345"},
				Name:      "Name One",
			},
		}, nil)

		out, _ := suite.client.GetBooks(suite.ctx, &pb.GetBooksRequest{Names: []string{"Name One"}})

		assert.NotNil(t, out)
		assert.Equal(t, "12345", out.Content[0].Id)
	})

	suite.Run("ErrorWhenBookDoesntExist", func() {
		suite.mockRepo.(*MockRepo).On("GetByName", []string{"Name One"}).Return(nil, status.Error(codes.NotFound, "book not found"))
		_, err := suite.client.GetBooks(suite.ctx, &pb.GetBooksRequest{Names: []string{"Name One"}})

		assert.NotNil(t, err)
		assert.ErrorIs(t, err, status.Error(codes.NotFound, "book not found"))
	})
}
