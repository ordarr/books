package service

import (
	"context"
	books "github.com/ordarr/books/v1"
	"github.com/ordarr/data/core"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"os"
	"testing"
)

type BooksTestSuite struct {
	suite.Suite
	ctx      context.Context
	client   books.BooksClient
	db       *gorm.DB
	populate func() []*core.Book
	closer   func()
}

func insertTestBooks(db *gorm.DB, ctx context.Context) []*core.Book {
	book1 := &core.Book{
		Ids: core.Ids{
			Calibre:  1,
			Koreader: 2,
		},
		Title: "Title One",
	}
	book2 := &core.Book{
		Ids: core.Ids{
			Calibre:  2,
			Koreader: 3,
		},
		Title: "Title Two",
	}
	session := db.Session(&gorm.Session{Context: ctx})
	session.Create(&book1)
	session.Create(&book2)

	return []*core.Book{
		book1, book2,
	}
}

func (suite *BooksTestSuite) SetupSubTest() {
	suite.ctx = context.Background()

	_db := core.Connect(&core.Config{
		Type:    "sqlite",
		Name:    "ordarr.db",
		LogMode: true,
	})

	suite.populate = func() []*core.Book {
		return insertTestBooks(_db, suite.ctx)
	}

	_client, _closer := Server(core.BookRepository{DB: _db})

	suite.db = _db
	suite.client = _client
	suite.closer = _closer
}

func (suite *BooksTestSuite) TearDownSubTest() {
	_ = os.Remove("ordarr.db")
	suite.closer()
}

func TestBooksTestSuite(t *testing.T) {
	suite.Run(t, new(BooksTestSuite))
}
