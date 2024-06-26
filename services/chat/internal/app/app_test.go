package app

import (
	"chat/internal/app/mocks"
	"chat/internal/domain"
	"chat/internal/repository/errs"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApp_New(t *testing.T) {
	type testcase struct {
		conf *Config
		repo *mocks.LoadSaver
	}
	tests := []testcase{
		{
			conf: &Config{MessagesToLoad: 10},
			repo: mocks.NewLoadSaver(t),
		},
		{
			conf: &Config{MessagesToLoad: 1},
			repo: mocks.NewLoadSaver(t),
		},
		{
			conf: &Config{MessagesToLoad: 5},
			repo: mocks.NewLoadSaver(t),
		},
		{
			conf: &Config{MessagesToLoad: 6},
			repo: mocks.NewLoadSaver(t),
		},
		{
			conf: &Config{MessagesToLoad: 100},
			repo: mocks.NewLoadSaver(t),
		},
	}

	for _, test := range tests {
		app := New(test.repo, test.conf)
		assert.Equal(t, test.repo, app.repo)
		assert.Equal(t, test.conf.MessagesToLoad, app.messagesToLoad)
	}
}

func TestApp_SaveMessage_NoRepoError(t *testing.T) {
	type testcase struct {
		username      string
		message       string
		returnedError error
	}

	tests := [][]testcase{
		{
			testcase{username: "danil", message: "bred", returnedError: nil},
			testcase{username: "gleb", message: "hello", returnedError: nil},
			testcase{username: "maks", message: "bye-bye", returnedError: nil},
			testcase{username: "maks", message: "some message", returnedError: nil},
		},
		{
			testcase{username: "danil", message: "bred", returnedError: nil},
		},
		{
			testcase{username: "danil", message: "bred", returnedError: nil},
			testcase{username: "gleb", message: "hello", returnedError: nil},
			testcase{username: "maks", message: "bye-bye", returnedError: nil},
			testcase{username: "maks", message: "some message", returnedError: nil},
			testcase{username: "danil", message: "privet", returnedError: nil},
			testcase{username: "gleb", message: "bim bim bom bom", returnedError: nil},
			testcase{username: "gleb", message: "gew w2vt242vt3tv2v3", returnedError: nil},
			testcase{username: "danil", message: "vv2 32f32", returnedError: nil},
		},
	}

	for _, test := range tests {
		repo := mocks.NewLoadSaver(t)
		for _, tc := range test {
			repo.On(
				"SaveMessage",
				context.Background(),
				domain.Message{Username: tc.username, Text: tc.message},
			).
				Return(tc.returnedError)
		}

		app := New(repo, &Config{MessagesToLoad: 10})
		for _, tc := range test {
			err := app.SaveMessage(tc.message, tc.username)
			assert.Equal(t, tc.returnedError, err)
		}
	}
}

func TestApp_SaveMessage_WithRepoError(t *testing.T) {
	type testcase struct {
		username      string
		message       string
		returnedError error
	}

	tests := [][]testcase{
		{
			testcase{username: "danil", message: "bred", returnedError: errs.ErrInternal},
			testcase{username: "gleb", message: "hello", returnedError: errs.ErrInternal},
			testcase{username: "maks", message: "bye-bye", returnedError: errs.ErrInternal},
			testcase{username: "maks", message: "some message", returnedError: errs.ErrInternal},
		},
		{
			testcase{username: "danil", message: "bred", returnedError: errs.ErrInternal},
		},
		{
			testcase{username: "danil", message: "bred", returnedError: errs.ErrInternal},
			testcase{username: "gleb", message: "hello", returnedError: errs.ErrInternal},
			testcase{username: "maks", message: "bye-bye", returnedError: errs.ErrInternal},
			testcase{username: "maks", message: "some message", returnedError: errs.ErrInternal},
			testcase{username: "danil", message: "privet", returnedError: errs.ErrInternal},
			testcase{username: "gleb", message: "bim bim bom bom", returnedError: errs.ErrInternal},
			testcase{username: "gleb", message: "gew w2vt242vt3tv2v3", returnedError: errs.ErrInternal},
			testcase{username: "danil", message: "vv2 32f32", returnedError: errs.ErrInternal},
		},
	}

	for _, test := range tests {
		repo := mocks.NewLoadSaver(t)
		for _, tc := range test {
			repo.On(
				"SaveMessage",
				context.Background(),
				domain.Message{Username: tc.username, Text: tc.message},
			).
				Return(tc.returnedError)
		}

		app := New(repo, &Config{MessagesToLoad: 10})
		for _, tc := range test {
			err := app.SaveMessage(tc.message, tc.username)
			assert.Error(t, err)
		}
	}
}

func TestApp_LoadLastMessages(t *testing.T) {
	type testcase struct {
		count    int
		messages []domain.Message
		err      error
	}

	tests := []testcase{
		{
			count: 8,
			messages: []domain.Message{
				{Username: "danil", Text: "Hello, World"},
				{Username: "gleb", Text: "Hello"},
				{Username: "gleb", Text: "He  ewf  fewllo, World"},
				{Username: "maks", Text: "ggggg"},
				{Username: "danil", Text: "POKA"},
				{Username: "gleb", Text: "fwefwefwfewf f ew"},
				{Username: "gleb", Text: "Hell"},
				{Username: "gleb", Text: "Hellorld"},
			},
			err: nil,
		},
		{
			count: 1,
			messages: []domain.Message{
				{Username: "danil", Text: "Hello, World"},
			},
			err: nil,
		},
		{
			count:    10,
			messages: nil,
			err:      errs.ErrNotFound,
		},
	}

	for _, test := range tests {
		repo := mocks.NewLoadSaver(t)
		repo.On(
			"LoadMessages",
			context.Background(),
			test.count,
		).Return(test.messages, test.err)

		app := New(repo, &Config{MessagesToLoad: test.count})
		messages, err := app.LoadLastMessages()
		assert.Equal(t, test.messages, messages)
		if test.err != nil {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
