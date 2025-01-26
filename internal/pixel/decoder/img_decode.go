package decoder

import (
	"encoding/base64"
	"encoding/json"
	"pixeltrace/internal/tools"
)

func DecodeImgData(data string) (json.RawMessage, error) {
	dest := make(json.RawMessage, len(data))
	_, err := base64.StdEncoding.Decode(dest, tools.StringToSliceByte(data))
	if err != nil {
		return nil, err
	}
	return dest, nil
}
