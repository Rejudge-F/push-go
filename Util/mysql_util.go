package Util

import (
	"context"
	"errors"
	"fmt"
	log "github.com/cihub/seelog"
	"push-go/common"
	"push-go/dbservice"
)

var (
	mysqlDB = &dbservice.MysqlDB{
		nil,
		false,
	}
)

func InsertDataToDB() {
	log.Info("insert data to mysql....")
	TruncateTable("xiaomi_mall", "user")
	mysqlConn, err := mysqlDB.Get()
	if err != nil {
		log.Error("get mysql conn error", err)
		return
	}

	defer mysqlConn.Close()
	for i := 0; i < DATANUM; i++ {
		sqlStr := fmt.Sprintf("INSERT INTO xiaomi_mall.user VALUES(%d, \"%d\")", i, 99)
		mysqlConn.ExecContext(context.Background(), sqlStr)
	}
	log.Info("insert ok!")
}

func GetUserDatas(db, table string, offset *int, nums int) ([]*common.UserData, error) {
	if db == "" || table == "" || *offset < 0 || nums <= 0 {
		return nil, errors.New(fmt.Sprintf("GetUserDatas: db: %s, table: %s, offset: %d, nums: %d", db, table, offset, nums))
	}

	mysqlConn, err := mysqlDB.Get()
	defer mysqlConn.Close()

	if err != nil {
		log.Error("get mysql conn error", err)
		return nil, err
	}

	sqlStr := fmt.Sprintf("SELECT id, name FROM %s.%s LIMIT %d, %d", db, table, *offset, nums)
	rows, err := mysqlDB.Pool.Query(sqlStr)
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		return nil, err
	}

	result := make([]*common.UserData, 0, nums)

	for rows.Next() {
		tUserData := &common.UserData{}
		err = rows.Scan(&tUserData.UID, &tUserData.Name)
		if err != nil {
			return nil, err
		}
		result = append(result, tUserData)
	}

	*offset = *offset + len(result)
	//log.Debug("len = ", len(result))
	return result, nil
}

func TruncateTable(db, tab string) {
	mysqlConn, _ := mysqlDB.Get()
	mysqlConn.ExecContext(context.Background(), fmt.Sprintf("truncate TABLE %s.%s", db, tab))
}
