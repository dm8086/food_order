package common

import "gopkg.in/mgo.v2/bson"

type ApiLog struct {
	ID            bson.ObjectId          `bson:"_id"`
	Time          string                 `json:"time" bson:"time"`
	RequestId     string                 `json:"requestId" bson:"requestId"`
	ResponseTime  string                 `json:"responseTime" bson:"responseTime"`
	TTL           int                    `json:"ttl" bson:"ttl"`
	Method        string                 `json:"method" bson:"method"`
	ContentType   string                 `json:"contentType" bson:"contentType"`
	Uri           string                 `json:"uri" bson:"uri"`
	ClientIP      string                 `json:"clientIP" bson:"clientIP"`
	RequestHeader map[string]string      `json:"requestHeader" bson:"requestHeader"`
	RequestParam  interface{}            `json:"requestParam" bson:"requestParam"`
	ResponseMap   map[string]interface{} `json:"responseMap" bson:"responseMap"`
}
