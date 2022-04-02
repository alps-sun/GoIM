package sql

import "database/sql"

var s sql.NullString

func test(s *sql.NullString, db sql.DB) (string, error) {
	id := 1
	// 查询到结果后scan 扫描行
	err := db.QueryRow("select name from foo where id=?", id).Scan(&s)

	if s.Valid {
		return s.String, nil
	} else {
		return "", err
	}
}
