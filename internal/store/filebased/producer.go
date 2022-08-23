package filebased

import (
	"encoding/json"
	"os"

	"github.com/vlad-marlo/shortener/internal/store/model"
)

type producer struct {
	file    *os.File
	encoder *json.Encoder
}

func newProducer(filename string) (*producer, error) {
	file, err := os.OpenFile(filename, .O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return nil, err
	}
	return &producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *producer) WriteURL(u *model.URL) error {
	return p.encoder.Encode(&u)
}
