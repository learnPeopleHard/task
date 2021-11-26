package check

import (
	"bytes"
	"encoding/gob"
)

type CheckAuthDTO struct{
	Token string
	FlagGetAndSet bool
}

func Encode(data CheckAuthDTO) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Decode(b []byte) (CheckAuthDTO, error) {
	buf := bytes.NewBuffer(b)
	decoder := gob.NewDecoder(buf)
	var data CheckAuthDTO
	if err := decoder.Decode(&data); err != nil {
		return CheckAuthDTO{}, err
	}
	return data, nil
}
