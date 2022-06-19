package server

import (
	"encoding/base64"
	"fmt"
	"hash/crc64"
)

type FontStyles int32

const (
	FontStyles_NORMAL FontStyles = 0
	FontStyles_ITALIC FontStyles = 1
)

type ConvertRequestData struct {
	InputText string     `json:"input_text"`
	FontSize  int32      `json:"font_size"`
	FontFile  string     `json:"font_file"`
	FontStyle FontStyles `json:"font_style"`
}

type ConvertRequest struct {
	ConvertRequestData

	ConvID string `json:"conv_id"`
}

var crctab = crc64.MakeTable(crc64.ISO)

func (cr ConvertRequestData) UniqueID() string {
	str := fmt.Sprintf("%s%d%s%d", cr.InputText, cr.FontSize, cr.FontFile, cr.FontStyle)
	crcval := crc64.Checksum([]byte(str), crctab)
	strval := fmt.Sprint(crcval)

	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte(strval))
}
