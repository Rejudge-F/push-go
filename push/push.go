package push

import (
	"fmt"
	log "github.com/cihub/seelog"
	"push-go/Util"
	"push-go/common"
	"sync"
	"time"
)

type IPush interface {
	Init()
	UserDataToChan(uData []*common.UserData)
}

type Push struct {
	DataChan1 chan *common.UserData
	DataChan2 chan *common.UserData
	offset    int
	nums      int
	mutex     sync.Mutex
	hasPush   int32
}

func (push *Push) Init() {
	push.DataChan1 = make(chan *common.UserData, Util.CHANSIZE)
	push.DataChan2 = make(chan *common.UserData, Util.CHANSIZE)
	push.offset = 0
	push.nums = Util.PERNUM
	push.hasPush = 0
}

func (push *Push) UserDataToChan(uData chan *common.UserData) {
	for {
		select {
		case d := <-uData:
			if len(push.DataChan1) < len(push.DataChan2) {
				push.DataChan1 <- d
			} else {
				push.DataChan2 <- d
			}
		}

	}
}

func PushMsg(dataChan chan *common.UserData, hasPush *int32) {
	for {
		select {
		case uData := <-dataChan:
			Util.PushMsg(uData, hasPush)
		}
		if *hasPush == Util.DATANUM {
			break
		}
	}
}

func (push *Push) pullMsg() {
	childChan := make(chan *common.UserData, Util.CHANSIZE)
	for {
		push.mutex.Lock()
		uData, err := Util.GetUserDatas("xiaomi_mall", "user", &push.offset, push.nums)
		push.mutex.Unlock()
		if err != nil {
			panic(err)
		}
		//if len(uData) == 0 {
		//	break
		//}

		go Util.UserDataToChan(uData, childChan)
		go push.UserDataToChan(childChan)

		if push.offset == Util.DATANUM {
			log.Info("exit pullMsg")
			break
		}
	}
}

func Run() {
	push := &Push{}
	push.Init()
	t := time.Now()
	for i := 0; i < Util.THREADNUM; i++ {
		go push.pullMsg()
		go PushMsg(push.DataChan1, &push.hasPush)
		go PushMsg(push.DataChan2, &push.hasPush)
	}

	ticker := time.NewTicker(time.Second)

	for {
		if push.hasPush == Util.DATANUM {
			eplased := time.Since(t)
			log.Debug(fmt.Sprintf("%d Msg has pushed in elapsed: %v\n", push.hasPush, eplased))
			break
		}
		select {
		case <-ticker.C:
			elapsed := time.Since(t)
			log.Info("hasPushed: ", push.hasPush, "time: ", elapsed)
		}

	}
}
