package inmemory

import (
	"context"
	"encoding/json"
	"server/internal/domain"
	"sync"
)

type Repository struct {
	data [][]byte
	mu   *sync.Mutex
}

func NewRepository() *Repository {
	return &Repository{
		mu:   &sync.Mutex{},
		data: make([][]byte, 0),
	}
}

func (r *Repository) SaveMessage(_ context.Context, message domain.Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	r.mu.Lock()
	r.data = append(r.data, data)
	r.mu.Unlock()
	return nil
}

func (r *Repository) LoadMessages(_ context.Context, count int) ([]domain.Message, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	count = min(count, len(r.data))
	res := make([]domain.Message, count)

	for i := len(r.data) - count; i < len(r.data); i++ {
		var msg domain.Message

		err := json.Unmarshal(r.data[i], &msg)
		if err != nil {
			return nil, err
		}

		res[i-len(r.data)+count] = msg
	}

	return res, nil
}
