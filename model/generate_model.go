package model

import (
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

// GenerateModel 根据表结构生成model文件
func GenerateModel(table string) {
	g := gen.NewGenerator(gen.Config{
		OutPath: "./model",            // struct输出目录
		Mode:    gen.WithDefaultQuery, // 生成默认查询方法
	})

	// 连接数据库（支持MySQL、PostgreSQL等）
	db, _ := gorm.Open(mysql.Open("root:password@tcp(127.0.0.1:3306)/gin-init?charset=utf8mb4"))
	g.UseDB(db)

	// 生成所有表的struct
	g.GenerateAllTable()
	// 或指定表：g.GenerateModel("users")、g.GenerateModels("users", "orders")

	g.Execute()
}
