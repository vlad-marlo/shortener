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

func (p *producer) truncateFile() error {
	fileInfo, err := p.file.Stat()
	if err != nil {
		return err
	}
	return p.file.Truncate(fileInfo.Size())
}

func (p *producer) getURLs() (data *URLs, err error) {
	if err = p.decoder.Decode(&data); err != nil && err != io.EOF {
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
	if err = p.truncateFile(); err != nil {
		return err
	}
	if _, ok := data.URLs[u.ID]; ok {
		// без этого при попытке сохранить урл с существующим ID все данные удалятся
		if err = p.encoder.Encode(&data); err != nil {
			return err
		}
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
