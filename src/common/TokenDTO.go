package common

import (
	"bytes"
	"encoding/gob"
)

type TokenDTO struct{
	UserId int64
	UserName string
	Token string
	LoginTime int64
}
func Encode(data TokenDTO) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Decode(b []byte) (TokenDTO, error) {
	buf := bytes.NewBuffer(b)
	decoder := gob.NewDecoder(buf)
	var data TokenDTO
	if err := decoder.Decode(&data); err != nil {
		return TokenDTO{}, err
	}
	return data, nil
}