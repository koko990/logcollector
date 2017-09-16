package types

import (
	"fmt"
	"time"
)

type ServiceDashboard struct {
	Id                           int    `json:"id" orm:"pk;auto"`
	ServiceName                  string `json:"service_name" orm:"column(service_name)"`
	PodNumber                    int    `json:"pod_number" orm:"column(pod_number)"`
	ContainerNumber              int    `json:"container_number" orm:"column(container_number)"`
	Timestamp                    int    `json:"timestamp" orm:"column(time_list_id)"`
	AverageMinutePodNumber       int    `json:"average_minute_pod_number" orm:"column(average_minute_pod_number)"`
	AverageMinuteContainerNumber int    `json:"average_minute_container_number" orm:"column(average_minute_container_number)"`
	AverageHourPodNumber         int    `json:"average_hour_pod_number" orm:"column(average_hour_pod_number)"`
	AverageHourContainerNumber   int    `json:"average_hour_container_number" orm:"column(average_hour_container_number)"`
	AverageDayPodNumber          int    `json:"average_day_pod_number" orm:"column(average_day_pod_number)"`
	AverageDayContainerNumber    int    `json:"average_day_container_number" orm:"column(average_day_container_number)"`
}
type DashboardNode struct {
	Id           int64    `json:"id" orm:"pk;auto"`
	NodeName     string `json:"pod_name" orm:"column(node_name)"`
	InternalIp   string `json:"ip" orm:"column(ip)"`
	CpuUsage     float32 `json:"cpu_usage" orm:"column(cpu_usage)"`
	MemUsage     float32 `json:"mem_usage" orm:"column(mem_usage)"`
	Timestamp    int `json:"pod_name" orm:"column(time_list_id)"`
	StorageTotal int64 `json:"pod_name" orm:"column(storage_total)"`
	StorageUse   int64 `json:"pod_name" orm:"column(storage_use)"`
}

func (s ServiceDashboard) TableName() string {
	timeTemp := time.Now().Format("2006_01_02")
	table := fmt.Sprintf(`dashboard_service_%s`, timeTemp)
	return table
}
func (s DashboardNode) TableName() string {
	timeTemp := time.Now().Format("2006_01_02")
	table := fmt.Sprintf(`dashboard_node_%s`, timeTemp)
	return table
}
