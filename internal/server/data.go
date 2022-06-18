package server

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
