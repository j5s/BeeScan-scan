package servercheck

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

/*
创建人员：云深不知处
创建时间：2022/1/5
程序功能：mysql、pgsql、mssql检测单元
*/

func SQLCheck(hostType string, ip string, port string, user string, pwd string) bool {
	connectStr := ""

	switch hostType {
	case "mysql":
		connectStr = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?timeout=%ds",
			user, pwd, ip, port, "", 6,
		)
	case "postgres":
		connectStr = fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s sslmode=disable password=%s timeout=%ds",
			ip, port, user, "", pwd, 6,
		)
	case "mssql":
		connectStr = fmt.Sprintf(
			"server=%s;user id=%s;password=%s;port=%s;database=%s;timeout=%ds",
			ip, user, pwd, port, "", 6,
		)
	}

	db, err := gorm.Open(hostType, connectStr)
	if err != nil {
		return false
	}

	db.Close()
	return true
}