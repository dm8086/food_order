package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"order_food/global"
	"order_food/model/common"
	"order_food/utils"
	"strings"
	"time"

	"go.uber.org/zap"
	"gopkg.in/mgo.v2/bson"

	"github.com/gin-gonic/gin"
)

// LogLayout 日志layout
type LogLayout struct {
	Time      time.Time
	Metadata  map[string]interface{} // 存储自定义原数据
	Path      string                 // 访问路径
	Query     string                 // 携带query
	Body      string                 // 携带body数据
	IP        string                 // ip地址
	UserAgent string                 // 代理
	Error     string                 // 错误
	Cost      time.Duration          // 花费时间
	Source    string                 // 来源
}

type Logger struct {
	// Filter 用户自定义过滤
	Filter func(c *gin.Context) bool
	// FilterKeyword 关键字过滤(key)
	FilterKeyword func(layout *LogLayout) bool
	// AuthProcess 鉴权处理
	AuthProcess func(c *gin.Context, layout *LogLayout)
	// 日志处理
	Print func(LogLayout)
	// Source 服务唯一标识
	Source string
}

func (l Logger) SetLoggerMiddleware() gin.HandlerFunc {
	go handleAccessChannel()

	return func(c *gin.Context) {
		// 增加当是配置里不需要记录日志的接口时直接跳过
		apiMps := make(map[string]bool)
		apis := strings.Split(strings.Trim(global.GVA_CONFIG.NoApilog.Apis, ","), ",")
		for _, api := range apis {
			apiMps[api] = true
		}

		if apiMps[c.Request.URL.Path] {
			c.Next()
			return
		}

		bodyLogWriter := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bodyLogWriter

		// 开始时间
		start := time.Now()

		data, err := c.GetRawData()
		if err != nil {
			global.GVA_LOG.Error(fmt.Sprintf("GetRawData error:%s", err.Error()))
		}
		body := string(data)

		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data)) // 关键点
		// 处理请求
		c.Next()

		responseBody := bodyLogWriter.body.String()

		var req map[string]interface{}
		var result map[string]interface{}

		// 日志格式
		if strings.Contains(c.Request.RequestURI, "/docs") || c.Request.RequestURI == "/" {
			return
		}

		if responseBody != "" && responseBody[0:1] == "{" {
			err := json.Unmarshal([]byte(responseBody), &result)
			if err != nil {
				result = map[string]interface{}{"status": -1, "msg": "解析异常:" + err.Error()}
			}
		}

		// 结束时间
		endTime := time.Now()

		// 日志格式
		var params interface{}
		if strings.Contains(c.ContentType(), "application/json") && body != "" {
			utils.FromJSON(body, &req)
			params = req
		} else if strings.Contains(c.ContentType(), "x-www-form-urlencoded") || strings.Contains(c.ContentType(), "multipart/form-data") {
			params = utils.GinParamMap(c)
		} else if c.Request.Method != "DELETE" {
			return
		}
		// 增加GET请求也需要写入日志 202501021203 vii
		// else if c.Request.Method == "GET" {
		// 	return
		// }

		postLog := new(common.ApiLog)
		postLog.ID = bson.NewObjectId()
		postLog.Time = start.Format("2006-01-02 15:04:05")
		postLog.Uri = c.Request.RequestURI
		postLog.Method = c.Request.Method
		postLog.ContentType = c.ContentType()
		postLog.RequestHeader = utils.GinHeaders(c)
		ip := c.GetHeader("X-Forward-For")
		if ip == "" {
			ip = c.GetHeader("X-Real-IP")
			if ip == "" {
				ip = c.ClientIP()
			}
		}
		postLog.ClientIP = ip
		postLog.RequestParam = params
		postLog.ResponseTime = endTime.Format("2006-01-02 15:04:05")
		postLog.ResponseMap = result
		postLog.TTL = int(endTime.UnixNano()/1e6 - start.UnixNano()/1e6)

		accessLog := "|" + c.Request.Method + "|" + postLog.Uri + "|" + c.ClientIP() + "|" + endTime.Format("2006-01-02 15:04:05.012") + "|" + fmt.Sprintf("%vms", endTime.UnixNano()/1e6-start.UnixNano()/1e6)
		global.GVA_LOG.Debug(accessLog)
		// global.GVA_LOG.Debug(fmt.Sprintf("请求参数:%v", params))
		// global.GVA_LOG.Debug(fmt.Sprintf("请求头:%v", postLog.RequestHeader))
		// global.GVA_LOG.Debug(fmt.Sprintf("接口返回:%v", result))
		accessChannel <- utils.ToJSON(postLog)
	}
}

func DefaultLogger() gin.HandlerFunc {
	return Logger{
		Print: func(layout LogLayout) {
			// 标准输出,k8s做收集
			// v, _ := json.Marshal(layout)
			//fmt.Println(string(v))
			// global.GVA_LOG.Debug(string(v))
		},
		Source: "GVA",
	}.SetLoggerMiddleware()
}

var accessChannel = make(chan string, 100)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func handleAccessChannel() {
	if global.GVA_CONFIG.ApiLog.Enabled {
		for accessLog := range accessChannel {
			var postLog common.ApiLog
			json.Unmarshal([]byte(accessLog), &postLog)

			err := global.GVA_APILOG.C(global.GVA_CONFIG.ApiLog.Collection).Insert(postLog)
			if err != nil {
				global.GVA_LOG.Error("MongoDB写入错误", zap.Error(err))
			}
		}
		return
	}
}
