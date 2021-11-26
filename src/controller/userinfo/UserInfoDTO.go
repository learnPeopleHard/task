package userinfo

import (
	"bytes"
	"encoding/gob"
)

type UserInfoDTO struct{
	Id int64
	UserId int64
	UserName string
	Nickname string
	ProfilePicture string
}

func InfoEncode(data UserInfoDTO) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func InfoDecode(b []byte) (UserInfoDTO, error) {
	buf := bytes.NewBuffer(b)
	decoder := gob.NewDecoder(buf)
	var data UserInfoDTO
	if err := decoder.Decode(&data); err != nil {
		return UserInfoDTO{}, err
	}
	return data, nil
}