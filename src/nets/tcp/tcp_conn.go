/***********************************************************************
* @ tcp连接
* @ brief
	1、Notice：
		http的消息处理，是另开goroutine调用的，所以函数中可阻塞
		tcp的消息处理，是在readRoutine中及时调用的，所以函数中不能有阻塞调用
		否则“该条连接”的读会被挂起，c++中的话，整个系统的处理线程都会阻塞掉

	2、server端目前是一条连接两个goroutine(readLoop/writeLoop)
		假设5k玩家，就有1w个goroutine，太多了

	3、msghandler可考虑设计成：不执行逻辑，仅将消息加入buf队列，由一个goroutine来处理
		不过那5k个readRoutine貌似省不了哇，感觉单独一个goroutine处理消息也不会有性能提升
		且增加了风险，若某条消息有阻塞调用，后面的就得等了

	4、Rpc:
		g_rpc_response须加读写锁，与c++(多线程收-主线程顺序处理)不同，go是每个用户一条goroutine

	5、现在的架构是：每条连接各线程收数据，直接在io线程调注册的业务函数，对强交互的业务不友好
		要是做MMO之类的，可考虑像c++一样，io线程只负责收发，数据交付给全局队列，主线程逐帧处理，避免竞态

	6、io.ErrShortWrite
		这个错误比较奇怪，按道理只是本次io没有发完数据，底层网络仍是可以工作的，应该上层处理好多余数据，继续发才是……
		但我查了好些资料，还有github上的开源库，都是直接报错断开的
		有一种说法是这样的,如果发生short write了,说明网络很糟糕了,继续发送可能会更糟糕,断开才是对的
		这个看你策略是保证对,还是降低服务器压力
		C++库的方式一般是把数据积累到服务器缓存里，网络不好，积累爆炸后服务也会跪掉
		这种方式保证对,用内存做代价,其实较好的。go的channel就反过来,如果有问题要解决问题,而不是积累问题

* @ reconnect
	1、Accept返回的conn，先收一个包，内含connId
	2、connId为0表示新连接，挑选一个空闲的TCPConn，newTCPConn()
	3、不为0即重连，取对应TCPConn，若它关闭的，随即resetConn()

* @ 更稳定的连接
	*、参考项目
		【http://blog.codingnow.com/2014/02/connection_reuse.html】
		【https://github.com/funny/snet】
	*、新连接建立
		上行包：
			1、ConnID==0
			2、DH密钥
		下行包：
			1、加密的ConnID
			2、DH密钥
	*、连接修复(断线重连)
		上行包：
			1、旧有的ConnID
			2、Client已发送字节数
			3、Client已接收字节数
			4、密钥计算出的MD5
		服务器：
			1、验证合法性，失败立即断开
			2、根据ConnID定位旧连接，并下发“已发、已收字节数”作为重连回应
			3、再由Client上报的“已接收字节数”，计算出需重传数据，立即下发
			4、Client收到重连响应后，比较收发字节数差值来读取Server下发的重传数据

* @ author zhoumf
* @ date 2016-8-3
***********************************************************************/
package tcp

import (
	"common"
	"common/safe"
	"gamelog"
	"generate_out/rpc/enum"
	"net"
	"nets/rpc"
	"sync/atomic"
	"time"
)

const (
	kHeadLen          = 2 //头2字节存msgSize
	Msg_Size_Max      = 16 * 1024
	Writer_Cap        = 16 * 1024
	Delay_Delete_Conn = 60 * time.Second
)

func init() {
	rpc.G_HandleFunc[enum.Rpc_regist] = _Rpc_regist
	rpc.G_HandleFunc[enum.Rpc_unregist] = _Rpc_unregist
	rpc.G_HandleFunc[enum.Rpc_svr_accept] = _Rpc_svr_accept
}

type TCPConn struct {
	net.Conn
	reader       netbuf //包装conn减少conn.Read的io次数【client/test/net.go】
	writer       safe.ChanByte
	_isClose     uint32 //可设计成标记位
	delayDel     *time.Timer
	onDisConnect func()
	user         atomic.Value
}

func newTCPConn(conn net.Conn) *TCPConn {
	self := new(TCPConn)
	self.reader.Init(conn)
	self.writer.Init(2048)
	self.resetConn(conn)
	return self
}
func (self *TCPConn) resetConn(conn net.Conn) {
	self.Conn = conn
	conn.(*net.TCPConn).SetLinger(0) //关闭即丢弃数据
	self.reader.Reset(conn)
	self.writer.IsStop = false
	atomic.StoreUint32(&self._isClose, 0)
	if self.delayDel != nil {
		self.delayDel.Stop()
	}
}
func (self *TCPConn) Close() error {
	if !self.IsClose() {
		atomic.StoreUint32(&self._isClose, 1)
		return self.Conn.Close()
	}
	return nil
}
func (self *TCPConn) IsClose() bool { return atomic.LoadUint32(&self._isClose) == 1 }

func (self *TCPConn) GetUser() interface{}  { return self.user.Load() }
func (self *TCPConn) SetUser(v interface{}) { self.user.Store(v) }

func (self *TCPConn) CallRpc(msgId uint16, sendFun, recvFun func(*common.NetPack)) {
	req := common.NewNetPackCap(32)
	rpc.MakeReq(req, msgId, sendFun, recvFun)
	self.WriteMsg(req)
	req.Free()
}
func (self *TCPConn) WriteMsg(msg *common.NetPack) {
	if self.writer.AddMsg(msg.Data()) > Writer_Cap {
		gamelog.Error("Writer_Cap")
		self.Close()
	}
}
func (self *TCPConn) WriteBuf(buf []byte) { self.writer.Add(buf) }
func (self *TCPConn) writeLoop() {
loop:
	for {
		b := self.writer.WaitGet() //block
		for pos, total := 0, len(b); ; {
			if n, e := self.Conn.Write(b[pos:]); e != nil {
				gamelog.Debug(e.Error())
				break loop
			} else if pos += n; pos == total {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
	self.Close()
}
