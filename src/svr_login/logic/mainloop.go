package logic

import (
	"tcp"
	"time"
)

func MainLoop() {
	//timeNow, timeOld, time_elapse := time.Now().UnixNano()/int64(time.Millisecond), int64(0), 0
	for {
		//timeOld = timeNow
		//timeNow = time.Now().UnixNano() / int64(time.Millisecond)
		//time_elapse = int(timeNow - timeOld)

		tcp.G_RpcQueue.Update()

		time.Sleep(50 * time.Millisecond)
	}
}
