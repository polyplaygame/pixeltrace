package geoip

import (
	"context"
	"net"
	"os"
	"path/filepath"
	"sync"

	mm "github.com/oschwald/geoip2-golang"
	"github.com/pkg/errors"
)

var _ Reader = (*MaxMindReader)(nil)

const (
	// CityDBName city数据库名称
	CityDBName = "GeoLite2-City.mmdb"
	// AsnDBName asn数据库名称
	AsnDBName = "GeoLite2-ASN.mmdb"
)

var _cityDB *mm.Reader
var _asnDB *mm.Reader
var mmReader = &MaxMindReader{}
var mmOnce sync.Once

type MaxMindReader struct {
	ctx      context.Context
	cancel   context.CancelFunc
	CityPath string
	AsnPath  string
}

// NewMaxMindReader 创建MaxMindReader
// dirPath 数据库文件目录
// 说明：文件名称固定，GeoLite2-City.mmdb、GeoLite2-ASN.mmdb
func NewMaxMindReader(ctx context.Context, dirPath string) (*MaxMindReader, error) {
	mmOnce.Do(func() {
		mmReader.ctx, mmReader.cancel = context.WithCancel(ctx)
		cityPath := filepath.Join(dirPath, CityDBName)
		asnPath := filepath.Join(dirPath, AsnDBName)
		mmReader.CityPath = cityPath
		mmReader.AsnPath = asnPath
		if err := loadCityDB(cityPath); err != nil {
			panic(errors.Wrap(err, "初始化city数据库错误"))
		}
		if err := loadASNDB(asnPath); err != nil {
			panic(errors.Wrap(err, "初始化asn数据库错误"))
		}
		go watchForChanges(ctx, cityPath, loadCityDB)
		go watchForChanges(ctx, asnPath, loadASNDB)
	})
	return mmReader, nil
}

func loadCityDB(dbPath string) error {
	// 判断文件是否存在
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return errors.WithMessage(err, "db file not exist")
	}
	cdb, err := mm.Open(dbPath)
	if err != nil {
		return errors.WithMessage(err, "open db file failed")
	}
	_cityDB = cdb
	return nil
}

func loadASNDB(dbPath string) error {
	// 判断文件是否存在
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return errors.WithMessage(err, "db file not exist")
	}
	adb, err := mm.Open(dbPath)
	if err != nil {
		return errors.WithMessage(err, "open db file failed")
	}
	_asnDB = adb
	return nil
}

// Lookup 查询IP 信息
// Todo: locale 的处理
func (m *MaxMindReader) Lookup(ipStr, locale string) (*Result, error) {
	ip := net.ParseIP(ipStr)
	cityRecord, err := _cityDB.City(ip)
	if err != nil {
		return nil, err
	}
	if _asnDB != nil {
		asnRecord, err := _asnDB.ASN(ip)
		if err != nil {
			return nil, err
		}
		return newResult(ipStr, cityRecord, asnRecord, locale)
	}
	return newResult(ipStr, cityRecord, nil, locale)
}

func newResult(ipStr string, city *mm.City, asn *mm.ASN, locale string) (*Result, error) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	var result *Result
	if city != nil {
		result = &Result{
			IP:                     ipStr,
			PostalCode:             city.Postal.Code,
			LocationTimezone:       city.Location.TimeZone,
			LocationLatitude:       city.Location.Latitude,
			LocationLongitude:      city.Location.Longitude,
			LocationMetroCode:      city.Location.MetroCode,
			LocationAccuracyRadius: city.Location.AccuracyRadius,
			IsInEU:                 city.Country.IsInEuropeanUnion,
			IsAnonymousProxy:       city.Traits.IsAnonymousProxy,
			IsSatelliteProvider:    city.Traits.IsSatelliteProvider,
		}
	}
	if city.City.Names != nil {
		result.CityNameID = city.City.GeoNameID
		name, ok := city.City.Names[locale]
		if ok {
			result.CityName = name
		} else {
			result.CityName = city.City.Names["en"]
		}
	}
	if city.Continent.Names != nil {
		result.ContinentID = city.Continent.GeoNameID
		result.ContinentCode = city.Continent.Code
		if name, ok := city.Continent.Names[locale]; ok {
			result.Continent = name
		} else {
			result.Continent = city.Continent.Names["en"]
		}
	}
	if city.Country.Names != nil {
		result.CountryID = city.Country.GeoNameID
		result.CountryCode = city.Country.IsoCode
		if name, ok := city.Country.Names[locale]; ok {
			result.CountryName = name
		} else {
			result.CountryName = city.Country.Names["en"]
		}
	}
	if city.RegisteredCountry.Names != nil {
		result.RegisteredCountryID = city.RegisteredCountry.GeoNameID
		result.RegisteredCountryCode = city.RegisteredCountry.IsoCode
		if name, ok := city.RegisteredCountry.Names[locale]; ok {
			result.RegisteredCountryName = name
		} else {
			result.RegisteredCountryName = city.RegisteredCountry.Names["en"]
		}
	}
	if city.RepresentedCountry.Names != nil {
		result.RepresentedCountryID = city.RepresentedCountry.GeoNameID
		result.RepresentedCountryCode = city.RepresentedCountry.IsoCode
		if name, ok := city.RepresentedCountry.Names[locale]; ok {
			result.RepresentedCountryName = name
		} else {
			result.RepresentedCountryName = city.RepresentedCountry.Names["en"]
		}
	}
	if city.Subdivisions != nil && len(city.Subdivisions) > 0 {
		result.SubdivisionID = city.Subdivisions[0].GeoNameID
		result.SubdivisionCode = city.Subdivisions[0].IsoCode
		if name, ok := city.Subdivisions[0].Names[locale]; ok {
			result.SubdivisionName = name
		} else {
			result.SubdivisionName = city.Subdivisions[0].Names["en"]
		}
	}
	if asn != nil {
		result.AutonomousSystemNumber = asn.AutonomousSystemNumber
		result.AutonomousSystemOrg = asn.AutonomousSystemOrganization
	}
	return result, nil
}
