package decoder

import (
	"encoding/base64"
	"fmt"
	"pixeltrace/internal/tools"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	data = `{
    "shop_name": "shoplebon.com",
    "raw_body": {
        "identities": {
            "$identity_cookie_id": "18bb2b4522e1d43-068cbf50b2ac5a4-26031151-3686400-18bb2b4522f109b"
        },
        "distinct_id": "18bb2b4522e1d43-068cbf50b2ac5a4-26031151-3686400-18bb2b4522f109b",
        "lib": {
            "$lib": "js",
            "$lib_method": "code",
            "$lib_version": "1.25.21"
        },
        "properties": {
            "$timezone_offset": -480,
            "$screen_height": 1440,
            "$screen_width": 2560,
            "$viewport_height": 581,
            "$viewport_width": 2560,
            "$lib": "js",
            "$lib_version": "1.25.21",
            "$latest_traffic_source_type": "引荐流量",
            "$latest_search_keyword": "未取到值",
            "$latest_referrer": "https://shoplebon.com/",
            "$referrer": "https://shoplebon.com/",
            "$url": "https://shoplebon.com/Women-s-Fashion/",
            "$url_path": "/Women-s-Fashion/",
            "$title": "Women's Fashion",
            "$is_first_day": true,
            "$is_first_time": false,
            "$referrer_host": "shoplebon.com"
        },
        "anonymous_id": "18bb2b4522e1d43-068cbf50b2ac5a4-26031151-3686400-18bb2b4522f109b",
        "type": "track",
        "event": "$pageview",
        "time": 1699511669829,
        "_track_id": 223549830,
        "_flush_time": 1699511669830
    },
    "m": {
        "ip": "::1",
        "time": "2023-11-09T15:46:22+08:00",
        "ua": "PostmanRuntime/7.33.0"
    }
}`
)

func TestEncodeImgData(t *testing.T) {
	dest := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(dest, tools.StringToSliceByte(data))
	fmt.Println(string(dest))
	assert.NotNil(t, dest)
}
