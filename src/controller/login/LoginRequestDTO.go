package login

import (
	"bytes"
	"encoding/gob"
)

type LoginRequestDTO struct{
	Name string
	Password string
}

func reqencode(data LoginRequestDTO) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func reqdecode(b []byte) (LoginRequestDTO, error) {
	buf := bytes.NewBuffer(b)
	decoder := gob.NewDecoder(buf)
	var data LoginRequestDTO
	if err := decoder.Decode(&data); err != nil {
		return LoginRequestDTO{}, err
	}
	return data, nil
}
