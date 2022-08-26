package filebased

import (
	"encoding/json"
	"os"

	"github.com/vlad-marlo/shortener/internal/store"
	"github.com/vlad-marlo/shortener/internal/store/model"
)

type producer struct {
	file    *os.File
	encoder *json.Encoder
}

func newProducer(filename string) (*producer, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *producer) CreateURL(u *model.URL) error {
	return p.encoder.Encode(&u)
}

func (p *producer) GetURLByID(id string) (u *model.URL, err error) {
	// TODO: write getting url logic
	return nil, store.ErrNotFound
}

func (p *producer) Close() error {
	return p.file.Close()
}
