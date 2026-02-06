package sea

import (
	"order_food/global"
	"time"
)

type EventRecord struct {
	global.SeaModel
	StoreId            string     `gorm:"column:store_id" json:"storeId"`
	StoreName          string     `gorm:"column:store_name" json:"storeName"`
	TableId            uint       `gorm:"column:table_id" json:"tableId"`
	DeskName           string     `gorm:"column:desk_name" json:"deskName"`
	AreaId             uint       `gorm:"column:area_id" json:"areaId"`
	AreaName           string     `gorm:"column:area_name" json:"areaName"`
	Date               string     `gorm:"column:date" json:"date"`
	RequestId          string     `gorm:"column:request_id" json:"requestId"`
	RequestTime        *time.Time `gorm:"column:request_time" json:"requestTime"`
	ResponseTime       *time.Time `gorm:"column:response_time" json:"responseTime"`
	CompletionTime     *time.Time `gorm:"column:completion_time" json:"completionTime"`
	ResponseDuration   int        `gorm:"column:response_duration" json:"responseDuration"`
	ProcessingDuration int        `gorm:"column:processing_duration" json:"processingDuration"`
	CompletionDuration int        `gorm:"column:completion_duration" json:"completionDuration"`
	EmployeeID         int        `gorm:"column:employee_id" json:"employeeId"`
	EventType          int        `gorm:"column:event_type" json:"eventType"`
	EventStatus        int        `gorm:"column:event_status" json:"eventStatus"`
}

// TableName 指定表名
func (EventRecord) TableName() string {
	return "event_records"
}

type StoreDailyStatistics struct {
	global.SeaModel
	StoreId               string  `gorm:"column:store_id" json:"storeId"`
	Date                  string  `gorm:"column:date" json:"date"`
	EventType             int     `gorm:"column:event_type" json:"eventType"`
	AvgResponseDuration   float64 `gorm:"column:avg_response_duration" json:"avgResponseDuration"`
	AvgProcessingDuration float64 `gorm:"column:avg_processing_duration" json:"avgProcessingDuration"`
	AvgCompletionDuration float64 `gorm:"column:avg_completion_duration" json:"avgCompletionDuration"`
	TotalEvents           int     `gorm:"column:total_events" json:"totalEvents"`
}

func (StoreDailyStatistics) TableName() string {
	return "store_daily_statistics"
}
