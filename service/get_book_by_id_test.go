package service

import (
	books "github.com/ordarr/books/v1"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *BooksTestSuite) TestGetBookById() {
	t := suite.T()

	suite.Run("ReturnsPopulatedBook", func() {
		inserted := suite.populate()

		out, _ := suite.client.GetBookById(suite.ctx, &books.ValueRequest{Value: inserted[0].ID})

		assert.NotNil(t, out)
		assert.Equal(t, inserted[0].Title, out.Content.Title)
	})

	suite.Run("ErrorWhenBookDoesntExist", func() {
		t := suite.T()

		_, err := suite.client.GetBookById(suite.ctx, &books.ValueRequest{Value: "4783e133-d856-43f4-8d38-9e50c5996cad"})

		assert.NotNil(t, err)
		assert.ErrorIs(t, err, status.Error(codes.NotFound, "book not found"))
	})
}
