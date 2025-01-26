package geoip

type Reader interface {
	Lookup(ipStr, locale string) (*Result, error)
}

type Result struct {
	IP                     string  `json:"ip" mapstructure:"ip"`
	PostalCode             string  `json:"postal_code,omitempty" mapstructure:"postal_code"`                           // 邮政编码
	CityNameID             uint    `json:"city_name_id,omitempty" mapstructure:"city_name_id"`                         // 城市名称ID
	CityName               string  `json:"city_name,omitempty" mapstructure:"city_name"`                               // 城市名称
	ContinentID            uint    `json:"continent_id,omitempty" mapstructure:"continent_id"`                         // 大洲ID
	ContinentCode          string  `json:"continent_code,omitempty" mapstructure:"continent_code"`                     // 大洲代码
	Continent              string  `json:"continent,omitempty" mapstructure:"continent"`                               // 大洲名称
	CountryID              uint    `json:"country_id,omitempty" mapstructure:"country_id"`                             // 国家ID
	CountryCode            string  `json:"country_code,omitempty" mapstructure:"country_code"`                         // 国家代码
	CountryName            string  `json:"country_name,omitempty" mapstructure:"country_name"`                         // 国家名称
	RegisteredCountryID    uint    `json:"registered_country_id,omitempty" mapstructure:"registered_country_id"`       // 注册国家ID
	RegisteredCountryName  string  `json:"registered_country_name,omitempty" mapstructure:"registered_country_name"`   // 注册国家名称
	RegisteredCountryCode  string  `json:"registered_country_code,omitempty" mapstructure:"registered_country_code"`   // 注册国家代码
	RepresentedCountryID   uint    `json:"represented_country_id,omitempty" mapstructure:"represented_country_id"`     // 代表国家ID
	RepresentedCountryName string  `json:"represented_country_name,omitempty" mapstructure:"represented_country_name"` // 代表国家名称
	RepresentedCountryCode string  `json:"represented_country_code,omitempty" mapstructure:"represented_country_code"` // 代表国家代码
	SubdivisionID          uint    `json:"subdivision_id,omitempty" mapstructure:"subdivision_id"`                     // 行政区ID
	SubdivisionName        string  `json:"subdivision_name,omitempty" mapstructure:"subdivision_name"`                 // 行政区名称
	SubdivisionCode        string  `json:"subdivision_code,omitempty" mapstructure:"subdivision_code"`                 // 行政区代码
	LocationTimezone       string  `json:"location_timezone,omitempty" mapstructure:"location_timezone"`               // 时区
	LocationLatitude       float64 `json:"location_latitude,omitempty" mapstructure:"location_latitude"`               // 纬度
	LocationLongitude      float64 `json:"location_longitude,omitempty" mapstructure:"location_longitude"`             // 经度
	LocationMetroCode      uint    `json:"location_metro_code,omitempty" mapstructure:"location_metro_code"`           // 地铁代码
	LocationAccuracyRadius uint16  `json:"location_accuracy_radius,omitempty" mapstructure:"location_accuracy_radius"` // 精度半径
	IsInEU                 bool    `json:"is_in_eu,omitempty" mapstructure:"is_in_eu"`                                 // 是否在欧盟
	IsAnonymousProxy       bool    `json:"is_anonymous_proxy,omitempty" mapstructure:"is_anonymous_proxy"`             // 是否匿名代理
	IsSatelliteProvider    bool    `json:"is_satellite_provider,omitempty" mapstructure:"is_satellite_provider"`       // 是否卫星提供商
	AutonomousSystemNumber uint    `json:"autonomous_system_number,omitempty" mapstructure:"autonomous_system_number"` // 自治系统号
	AutonomousSystemOrg    string  `json:"autonomous_system_org,omitempty" mapstructure:"autonomous_system_org"`       // 自治系统组织
}
