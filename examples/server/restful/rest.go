package restful

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/tinywell/fabclient/pkg/common"
	"github.com/tinywell/fabclient/pkg/handler"
)

// RESTful接口地址
const (
	WaitResponseTimeOut = time.Second * 60

	URLData = "/default/data"
)

// RestServer RESTful接口
type RestServer struct {
	addr     string             // 端口
	router   *httprouter.Router // 路由
	tls      bool
	certFile string
	keyFile  string
	token    bool
	msgs     chan handler.Message
}

// NewServer 生成新RestServer对象
// params：
//   - addr string  服务监听地址
// return:
//   - RestServer  RESTful接口对象
func NewServer(addr string) *RestServer {
	server := &RestServer{
		addr:   addr,
		router: httprouter.New(),
		msgs:   make(chan handler.Message),
	}
	server.run()
	return server
}

// ReceiveMessage return messager channel
func (server *RestServer) ReceiveMessage() <-chan handler.Message {
	return server.msgs
}

func (server *RestServer) run() error {

	server.router.POST(URLData, server.saveData)
	server.router.GET(URLData, server.readData)

	go func() {

		err := http.ListenAndServe(server.addr, server.router)
		if err != nil {
			panic(fmt.Errorf("http.ListenAndServe error:%s", err.Error()))
		}
		fmt.Printf("server listen at %s\n", server.addr)
	}()

	return nil
}

// SaveData save tx data
func (server *RestServer) saveData(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	// 从http请求body中读取数据
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Read Body error:%s", err.Error())
		return
	}

	qidata := handler.Message{
		TranCode: "EXP100",
		TranData: data,
	}

	// 此处不能用 goroutine 处理，当前方法结束后，w会关闭，导致返回信息无法写入
	server.waitResponse(w, qidata)
}

// ReadData read tx data
func (server *RestServer) readData(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// 获取 get 请求参数 key
	name := r.FormValue("name")
	// name := params.ByName("key")
	// 调用交易数据查询接口UnionChainReadData
	qidata := handler.Message{
		TranCode: "EXP200",
		TranData: []byte(name),
	}
	server.waitResponse(w, qidata)
}

// SendReturn 将处理结果返回请求端
// params：
// - w http.ResponseWriter        响应输出
// - data interface{}				 返回数据
func SendReturn(w http.ResponseWriter, data interface{}) {
	msg, err := json.Marshal(data)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	fmt.Fprintf(w, string(msg))
}

func (server *RestServer) waitResponse(w http.ResponseWriter, qidata handler.Message) {
	rst := make(chan handler.Result)
	qidata.Result = rst

	select {
	case server.msgs <- qidata:
	case <-time.After(WaitResponseTimeOut):
		SendReturn(w, handler.Result{
			RspCode: common.RspServerError,
			RspData: []byte("send data time out"),
		})
		return
	}

	select {
	case result := <-rst:
		if result.RspCode != common.RspSuccess {
		}
		SendReturn(w, result)
		return
	case <-time.After(WaitResponseTimeOut):
		SendReturn(w, handler.Result{
			RspCode: common.RspTimeout,
			RspData: []byte("wait response time out"),
		})
		return
	}
}
