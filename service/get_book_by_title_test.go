package service

import (
	"github.com/google/uuid"
	books "github.com/ordarr/books/v1"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (suite *BooksTestSuite) TestGetBookByTitleReturnsBookWhenFound() {
	t := suite.T()

	suite.Run("ReturnsPopulatedBook", func() {
		suite.populate()

		out, _ := suite.client.GetBookByTitle(suite.ctx, &books.ValueRequest{Value: "Title One"})

		assert.NotNil(t, out)
		assert.NoError(t, uuid.Validate(out.Content.Id))
	})

	suite.Run("ErrorWhenBookDoesntExist", func() {
		_, err := suite.client.GetBookByTitle(suite.ctx, &books.ValueRequest{Value: "some-random-id"})

		assert.NotNil(t, err)
		assert.ErrorIs(t, err, status.Error(codes.NotFound, "book not found"))
	})
}
