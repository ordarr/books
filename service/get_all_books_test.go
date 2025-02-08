package service

import (
	pb "github.com/ordarr/books/v1"
	"github.com/ordarr/data/core"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/suite"
)

func (suite *BookTestSuite) TestGetAllBooks() {
	t := suite.T()

	suite.Run("ReturnsPopulatedList", func() {
		suite.mockRepo.(*MockRepo).On("GetAll").Return([]*core.Book{
			{
				BaseTable: core.BaseTable{ID: "12345"},
				Name:      "Book One",
			},
		}, nil)
		out, _ := suite.client.GetBooks(suite.ctx, &pb.GetBooksRequest{})

		assert.NotNil(t, out)
		assert.Len(t, out.Content, 1)
	})
}
