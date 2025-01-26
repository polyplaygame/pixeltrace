package main

import (
	"pixeltrace/biz/dal/model"

	"gorm.io/gen"
)

//go:generate go run main.go

func main() {
	// 	type Config struct {
	//     OutPath      string // 查询类代码的输出路径
	//     OutFile      string // 代码输出文件名，默认: gen.go
	//     ModelPkgPath string // 生成的 model 包名
	//     WithUnitTest bool   // 是否为生成的查询类代码生成单元测试
	//     FieldNullable     bool // generate pointer when field is nullable
	//     FieldCoverable    bool // generate pointer when field has default value, to fix problem zero value cannot be assign: https://gorm.io/docs/create.html#Default-Values
	//     FieldSignable     bool // detect integer field's unsigned type, adjust generated data type
	//     FieldWithIndexTag bool // generate with gorm index tag
	//     FieldWithTypeTag  bool // generate with gorm column type tag
	//     Mode GenerateMode // generator modes
	// }
	cfg := gen.Config{
		OutPath: "../query",
		// WithUnitTest: true,
		// FieldNullable:  true,
		// FieldCoverable: true,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
		Mode:              gen.WithDefaultQuery | gen.WithQueryInterface,
	}
	g := gen.NewGenerator(cfg)
	g.ApplyBasic(
		&model.AppCode{},
	)
	// g.ApplyInterface(func(models.CommonOp) {},
	// 	models.Contact{},
	// 	models.ProductCategory{},
	// )
	g.Execute()
}
