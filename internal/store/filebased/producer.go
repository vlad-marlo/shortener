package filebased

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/vlad-marlo/shortener/internal/store"
	"github.com/vlad-marlo/shortener/internal/store/model"
)

type producer struct {
	file    *os.File
	encoder *json.Encoder
	decoder *json.Decoder
}

func newProducer(filename string) (*producer, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &producer{
		file:    file,
		encoder: json.NewEncoder(file),
		decoder: json.NewDecoder(file),
	}, nil
}

func (p *producer) CreateURL(u *model.URL) error {
	return p.encoder.Encode(&u)
}

func (p *producer) GetURLByID(id string) (u *model.URL, err error) {
	for {
		err := p.decoder.Decode(&u)
		log.Println(u)
		if u != nil && u.ID == id {
			return u, nil
		}
		if err != nil {
			if err == io.EOF {
				return nil, store.ErrNotFound
			}
			return nil, err
		}
	}
}

func (p *producer) Close() error {
	return p.file.Close()
}
