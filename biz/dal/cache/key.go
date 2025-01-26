package cache

import "fmt"

const (
	PixelAppCodeKey = "pixel_app_code"
)

func GetAppCodeKey(code string) string {
	return fmt.Sprintf("%s:%s", PixelAppCodeKey, code)
}
