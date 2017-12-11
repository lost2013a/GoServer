/***********************************************************************
* @ 充值订单数据库
* @ brief
    1、gamesvr先通知SDK进程，建立新充值订单

    2、第三方充值信息到达后，验证是否为有效订单，通过后入库

* @ author zhoumf
* @ date 2016-8-18
***********************************************************************/
package logic

import (
	"dbmgo"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type TOrderInfo struct {
	Order_id       string `bson:"_id"`
	Third_order_id string
	Pf_id          string //平台名（如oppo）
	Pk_id          string //包id（限定2位）
	Pay_id         int    //支付商id（不同平台下的同一种支付渠道pay_id必须一样）
	Op_id          int    //运营商 注：1代表移动。2代表联通 3代表电信
	App_id         string //应用id
	Server_id      int    //服务器id（只有一个服务器的话，默认传1）
	Account        string //账号（没有必须传空字符）
	Third_account  string //第三方账号，相当于登录账号
	Role_id        string //角色id
	Item_id        int    //物品id
	Item_name      string //物品名（中文需要urlencode）
	Item_count     int    //物品数量
	Item_price     int    //物品价格（单位是分）
	Total_price    int    //物品总价格（单位是分）
	Currency       string //货币（在大陆，直接填RMB即可）
	Code           string //计费点
	Imsi           string
	Imei           string
	Ip             string
	Net            string //网络类型 CMNET， WIFI等
	Status         int    //1成功 0失败（第三方通告）
	Can_send       int    //1能发货 （client发货后置0）
	Time           int64
}

var (
	g_order_map = make(map[string]*TOrderInfo, 1024)
)

func CreateOrderInDB(ptr *TOrderInfo) {
	//生成订单号
	ptr.Order_id = fmt.Sprintf("%03d%s%06d", ptr.Pay_id, time.Now().Format("060102"), dbmgo.GetNextIncId("OrderId"))
	ptr.Time = time.Now().Unix()
	g_order_map[ptr.Order_id] = ptr
	dbmgo.InsertSync("Order", ptr)
}
func FindOrder(orderId string) *TOrderInfo {
	if ptr, ok := g_order_map[orderId]; ok {
		return ptr
	}
	return nil
}
func InitDB() {
	//删除数据库里超过30天的订单
	dbmgo.RemoveSync("Order", bson.M{"time": bson.M{"$lt": time.Now().Unix() - 30*24*3600}})
	//载入所有未完成订单
	var list []TOrderInfo
	dbmgo.FindAll("Order", bson.M{"can_send": 0}, &list)
	for i := 0; i < len(list); i++ {
		g_order_map[list[i].Order_id] = &list[i]
	}
}
