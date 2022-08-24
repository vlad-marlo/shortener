package filebased

import (
	"encoding/json"
	"log"
	"os"

	"github.com/vlad-marlo/shortener/internal/store"
	"github.com/vlad-marlo/shortener/internal/store/model"
)

type producer struct {
	file    *os.File
	decoder *json.Decoder
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
		decoder: json.NewDecoder(file),
	}, nil
}

func (p *producer) getURLs() (*URLs, error) {
	data := &URLs{}
	if err := p.decoder.Decode(&data); err != nil {
		return nil, err
	}
	log.Print("successfully get urls")
	return data, nil
}

func (p *producer) CreateURL(u *model.URL) error {
	data, err := p.getURLs()
	if err != nil {
		return err
	}
	if _, ok := data.URLs[u.ID]; ok {
		return store.ErrAlreadyExists
	}
	data.URLs[u.ID] = u
	return p.encoder.Encode(&data)
}

func (p *producer) GetURLByID(id string) (*model.URL, error) {
	data, err := p.getURLs()
	if err != nil {
		return nil, err
	}
	u, ok := data.URLs[id]
	if !ok {
		return nil, store.ErrNotFound
	}
	return u, nil
}

func (p *producer) Close() error {
	return p.file.Close()
}
