package service

import (
	"github.com/ordarr/data/core"
	"github.com/stretchr/testify/mock"
)

type MockRepo struct {
	mock.Mock
}

func (repo *MockRepo) GetAll() ([]*core.Book, error) {
	args := repo.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*core.Book), nil
}

func (repo *MockRepo) GetByName(titles []string) ([]*core.Book, error) {
	args := repo.Called(titles)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*core.Book), nil
}

func (repo *MockRepo) GetByID(ids []string) ([]*core.Book, error) {
	args := repo.Called(ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*core.Book), nil
}

func (repo *MockRepo) Create(entity *core.Book) (*core.Book, error) {
	args := repo.Called(entity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*core.Book), nil
}
