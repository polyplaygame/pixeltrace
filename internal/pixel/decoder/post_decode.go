package decoder

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

const (
	ZippedFlag = "1"
)

// decodeBase64AndUnzip 处理通用的解码和解压逻辑
func decodeBase64AndUnzip(data, ziped string) ([]byte, error) {
	if data == "" {
		return nil, fmt.Errorf("empty data")
	}

	// base64解码
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, errors.WithMessage(err, "base64 decode failed")
	}

	// 解压
	if ziped == ZippedFlag {
		decoded, err = unGzip(decoded)
		if err != nil {
			return nil, errors.WithMessage(err, "ungzip failed")
		}
	}

	return decoded, nil
}

// DecodePostDataList 解析 post 数据列表
func DecodePostDataList(data, ziped string) ([]json.RawMessage, error) {
	decoded, err := decodeBase64AndUnzip(data, ziped)
	if err != nil {
		return nil, errors.WithMessage(err, "decode data failed")
	}

	var dataList []json.RawMessage
	if err := jsoniter.Unmarshal(decoded, &dataList); err != nil {
		return nil, errors.WithMessage(err, "unmarshal json failed")
	}
	return dataList, nil
}

// DecodePostData 解析 post 数据
func DecodePostData(data, ziped string) (json.RawMessage, error) {
	decoded, err := decodeBase64AndUnzip(data, ziped)
	if err != nil {
		return nil, errors.WithMessage(err, "decode data failed")
	}
	return json.RawMessage(decoded), nil
}

func unGzip(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, errors.WithMessage(err, "new gzip reader failed")
	}
	defer reader.Close()
	uncompressed, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.WithMessage(err, "read gzip data failed")
	}

	return uncompressed, nil
}
