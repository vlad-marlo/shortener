package filebased

import (
	"encoding/json"
	"os"

	"github.com/vlad-marlo/shortener/internal/store"
	"github.com/vlad-marlo/shortener/internal/store/model"
)

type producer struct {
	file    *os.File
	decoder *json.Decoder
	encoder *json.Encoder
}

func newProducer(filename string) (p *producer, err error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	p.file = file
	return
}

func (p *producer) getURLs() (data urls, err error) {
	err = p.decoder.Decode(&data)
	return
}

func (p *producer) WriteURL(u *model.URL) (err error) {
	data, err := p.getURLs()
	if err != nil {
		return
	}
	if _, ok := data.URLs[u.ID]; ok {
		err = store.ErrAlreadyExists
		return
	}
	data.URLs[u.ID] = u
	return
}

func (p *producer) ReadURL(id string) (u *model.URL, err error) {
	data, err := p.getURLs()
	if err != nil {
		return nil, err
	}
	u, ok := data.URLs[id]
	if !ok {
		return nil, store.ErrNotFound
	}
	return
}

func (p *producer) Close() error {
	return p.file.Close()
}
