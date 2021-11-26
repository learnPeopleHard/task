package common

import (
	"bytes"
	"encoding/gob"
)

type BaseResponseDTO struct{
	Code int
	Message string
}

func encode(data BaseResponseDTO) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func decode(b []byte) (BaseResponseDTO, error) {
	buf := bytes.NewBuffer(b)
	decoder := gob.NewDecoder(buf)
	var data BaseResponseDTO
	if err := decoder.Decode(&data); err != nil {
		return BaseResponseDTO{}, err
	}
	return data, nil
}
