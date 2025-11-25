package db

import (
	"fmt"
	"gin-init/config" // 导入配置包

	"gorm.io/driver/mysql" // 根据数据库类型选择驱动
	"gorm.io/gorm"
)

var DB *gorm.DB // 全局 DB 实例，供其他模块调用

type Page struct {
	PageNum  int `json:"pageNum"`  // 当前页码
	PageSize int `json:"pageSize"` // 每页条数
}

// Init 初始化 GORM 连接
func Init() error {
	var err error
	// 根据配置中的数据库驱动初始化
	switch config.Cfg.Database.Driver {
	case "mysql":
		DB, err = gorm.Open(mysql.Open(config.Cfg.Database.Dsn), &gorm.Config{})
		// 其他数据库（如 PostgreSQL、SQLite）可在此扩展
	default:
		return fmt.Errorf("不支持的数据库驱动: %s", config.Cfg.Database.Driver)
	}
	if err != nil {
		return fmt.Errorf("数据库连接失败: %v", err)
	}

	// 配置连接池（可选，推荐）
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接池失败: %v", err)
	}
	sqlDB.SetMaxOpenConns(config.Cfg.Database.MaxOpenConns) // 从配置读取
	sqlDB.SetMaxIdleConns(config.Cfg.Database.MaxIdleConns)

	return nil
}

// SelectPage 执行原生SQL并返回分页数据
func SelectPage(db *gorm.DB, sql string, page Page, args ...interface{}) ([]map[string]interface{}, error) {
	var total int64
	// 拼接sql， 查询总行数
	tmp := fmt.Sprintf("SELECT COUNT(1) as total FROM (%s) AS t", sql)
	err := db.Raw(tmp, args...).Row().Scan(&total)
	if err != nil {
		return nil, err
	}

	// 2. 计算偏移量（防止 pageNum 小于 1）
	if page.PageNum <= 0 {
		page.PageNum = 1
	}
	offset := (page.PageNum - 1) * page.PageSize

	sql = fmt.Sprintf("SELECT * FROM (%s) AS t LIMIT %d OFFSET %d", sql, page.PageSize, offset)

	rows, err := db.Raw(sql, args...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for rows.Next() {
		// 创建一个map存储当前行数据
		row := make(map[string]interface{})
		// 为每个字段创建一个指针，用于接收数据
		values := make([]interface{}, len(columns))
		for i := range columns {
			var v interface{}
			values[i] = &v // 指针指向空接口
		}

		// 扫描数据到values
		if err := rows.Scan(values...); err != nil {
			return nil, err
		}

		// 将扫描到的值映射到map
		for i, col := range columns {
			// 取值时需要解引用（因为values[i]是指针）
			val := *(values[i].(*interface{}))
			// 处理字节数组为字符串（如varchar类型）
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}

		result = append(result, row)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// SelectList 执行原生SQL并返回 []map[string]interface{}
func SelectList(db *gorm.DB, sql string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := db.Raw(sql, args...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for rows.Next() {
		// 创建一个map存储当前行数据
		row := make(map[string]interface{})
		// 为每个字段创建一个指针，用于接收数据
		values := make([]interface{}, len(columns))
		for i := range columns {
			var v interface{}
			values[i] = &v // 指针指向空接口
		}

		// 扫描数据到values
		if err := rows.Scan(values...); err != nil {
			return nil, err
		}

		// 将扫描到的值映射到map
		for i, col := range columns {
			// 取值时需要解引用（因为values[i]是指针）
			val := *(values[i].(*interface{}))
			// 处理字节数组为字符串（如varchar类型）
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}

		result = append(result, row)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// SelectOne 执行原生SQL并返回 map[string]interface{}
func SelectOne(db *gorm.DB, sql string, args ...interface{}) (map[string]interface{}, error) {
	rows, err := db.Raw(sql, args...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// 只读取第一行
	if !rows.Next() {
		return nil, nil // 无记录
	}

	// 绑定字段值
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range columns {
		valuePtrs[i] = &values[i]
	}
	if err := rows.Scan(valuePtrs...); err != nil {
		return nil, err
	}

	// 转换为 map
	row := make(map[string]interface{})
	for i, col := range columns {
		val := values[i]
		if b, ok := val.([]byte); ok {
			row[col] = string(b)
		} else {
			row[col] = val
		}
	}

	return row, nil
}
