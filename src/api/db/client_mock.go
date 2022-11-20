package db

import (
	"database/sql"
	"fmt"

	"github.com/stretchr/testify/mock"
)

type ClientMock struct {
	mock.Mock
}

func NewClientMock() *ClientMock {
	return &ClientMock{}
}

func mockRow(args mock.Arguments, index int) *sql.Row {
	obj := args.Get(index)
	var s *sql.Row
	var ok bool
	if s, ok = obj.(*sql.Row); !ok {
		panic(fmt.Sprintf("assert: arguments: Row(%d) failed because object wasn't correct type: %v", index, args.Get(index)))
	}
	return s
}

func mockRows(args mock.Arguments, index int) *sql.Rows {
	obj := args.Get(index)
	var s *sql.Rows
	var ok bool
	if s, ok = obj.(*sql.Rows); !ok {
		panic(fmt.Sprintf("assert: arguments: Rows(%d) failed because object wasn't correct type: %v", index, args.Get(index)))
	}
	return s
}

func (m *ClientMock) Query(query string, params ...any) (*sql.Rows, error) {
	args := m.Called(query, params)
	return mockRows(args, 0), args.Error(1)
}

func (m *ClientMock) QueryRow(query string, params ...any) *sql.Row {
	args := m.Called(query, params)
	return mockRow(args, 0)
}
