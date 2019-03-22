package fasthttp

import (
	"common"
	"common/file"
	"conf"
	"encoding/json"
	"fmt"
	"gamelog"
	http "github.com/valyala/fasthttp"
	"netConfig/meta"
	"path"
	"svr_client/test/qps"
	"sync"
)

//idx1 := strings.Index(addr, "//") + 2
//idx2 := strings.LastIndex(addr, ":")
//ip := addr[idx1:idx2]
//port := common.Atoi(addr[idx2+1:])
func Addr(ip string, port uint16) string { return fmt.Sprintf("http://%s:%d", ip, port) }

var (
	_svr http.Server
	_map = map[string]http.RequestHandler{}
)

func NewHttpServer(port uint16, module string, svrId int) error {
	if conf.TestFlag_CalcQPS {
		go qps.WatchLoop()
	}
	g_svraddr_path = fmt.Sprintf("%s/%s/%d/reg_addr.csv", file.GetExeDir(), module, svrId)
	loadCacheNetMeta()
	_svr.Handler = func(ctx *http.RequestCtx) {
		if v, ok := _map[common.B2S(ctx.Path())]; ok {
			v(ctx)
		} else {
			ctx.NotFound()
		}
	}
	HandleFunc("/reg_to_svr", _reg_to_svr)
	return _svr.ListenAndServe(fmt.Sprintf(":%d", port))
}
func CloseServer() { _svr.Shutdown() }

func HandleFunc(path string, f http.RequestHandler) { _map[path] = f }

// ------------------------------------------------------------
//! 模块注册
func _reg_to_svr(ctx *http.RequestCtx) {
	pMeta := new(meta.Meta)
	if err := common.B2T(ctx.Request.Body(), pMeta); err != nil {
		gamelog.Error("common.B2T: " + err.Error())
		return
	}
	meta.AddMeta(pMeta)
	appendNetMeta(pMeta)
	fmt.Println("RegistToSvr: ", pMeta)
}

// ------------------------------------------------------------
//! 本地存储“远程服务”注册地址，以备Local Http NetSvr崩溃重启
//! 不似tcp，对端不知道这边傻逼了( ▔___▔)y
//! 采用追加方式，同个“远程服务”的地址，会被最新追加的覆盖掉
var (
	g_svraddr_path string
	_mutex         sync.Mutex
)

func loadCacheNetMeta() {
	records, err := file.ReadCsv(g_svraddr_path)
	if err != nil {
		return
	}
	metas := make([]meta.Meta, len(records))
	for i := 0; i < len(records); i++ {
		json.Unmarshal(common.S2B(records[i][0]), &metas[i])
		meta.AddMeta(&metas[i])
		fmt.Println("load meta cache: ", metas[i])
	}
}
func appendNetMeta(pMeta *meta.Meta) {
	if meta.GetMeta("zookeeper", 0) != nil {
		return //有zookeeper实现重启恢复，不必本地缓存
	}
	_mutex.Lock()
	records, err := file.ReadCsv(g_svraddr_path)
	_mutex.Unlock()
	if err == nil {
		pMeta2 := new(meta.Meta)
		for i := 0; i < len(records); i++ {
			json.Unmarshal(common.S2B(records[i][0]), pMeta2)
			if pMeta.IsSame(pMeta2) {
				return
			}
		}
	}

	b, _ := json.Marshal(pMeta)
	dir, name := path.Split(g_svraddr_path)
	_mutex.Lock()
	err = file.AppendCsv(dir, name, []string{common.B2S(b)})
	_mutex.Unlock()
	if err != nil {
		gamelog.Error("AppendSvrAddrCsv: " + err.Error())
	}
}