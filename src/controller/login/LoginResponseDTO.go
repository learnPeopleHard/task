package login

import (
	"bytes"
	"encoding/gob"
)

type LoginResponseDTO struct{
	Code int
	Message string
	UserId int64
	UserName string
	LoginToken string
}

func encode(data LoginResponseDTO) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func decode(b []byte) (LoginResponseDTO, error) {
	buf := bytes.NewBuffer(b)
	decoder := gob.NewDecoder(buf)
	var data LoginResponseDTO
	if err := decoder.Decode(&data); err != nil {
		return LoginResponseDTO{}, err
	}
	return data, nil
}
