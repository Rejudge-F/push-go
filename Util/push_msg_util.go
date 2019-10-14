package Util

import (
	"push-go/common"
	. "sync/atomic"
	"time"
)

func PushMsg(data *common.UserData, nowPush *int32) {
	time.Sleep(time.Duration(PUSHTIME) * time.Millisecond)
	AddInt32(nowPush, 1)
	//log.Debug(data.UID, " ", *nowPush)
}

func UserDataToChan(uData []*common.UserData, dataChan chan *common.UserData) {
	for _, d := range uData {
		dataChan <- d
	}
}
