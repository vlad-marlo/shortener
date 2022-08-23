package filebased

import (
	"encoding/json"
	"os"

	"github.com/vlad-marlo/shortener/internal/store/model"
)

type customer struct {
	file    *os.File
	decoder *json.Decoder
}

func newCustomer(filename string) (*customer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &customer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (c *customer) ReadURL(id string) (u *model.URL, err error) {
	var urls URLStore
	if err = c.decoder.Decode(&urls); err != nil {
		return
	}
	u = urls.URLs[id]
	return
}

func (c *customer) Close() error {
	return c.file.Close()
}
