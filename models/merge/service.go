package merge

import (
	"github.com/koko990/logcollector/service/storage"
	"strconv"

	"container/list"
	"github.com/koko990/logcollector/models/util/calculate"

	"github.com/koko990/logcollector/models/types"
	"github.com/koko990/logcollector/models/util/pool"

	modelK8s "k8s.io/api/core/v1"
)

func RunService(timestampArg string) {
	var pod modelK8s.PodList
	var serviceK8s modelK8s.ServiceList
	var ok bool
	if checkServiceCache(timestampArg) == true {
		is := pool.IndexService{Timestamp: timestampArg, ServiceKey: int(storage.Service)}
		vs := pool.MergeServiceCache[is]
		serviceK8s, ok = vs.(modelK8s.ServiceList)
		delete(pool.MergeServiceCache, is)
		ip := pool.IndexService{Timestamp: timestampArg, ServiceKey: int(storage.Pod)}
		vp := pool.MergeServiceCache[ip]
		pod, ok = vp.(modelK8s.PodList)
		if ok == false {
			panic(ok)
		}
		delete(pool.MergeServiceCache, ip)
		handleService(serviceK8s, pod, timestampArg)
	}
}

//
func checkServiceCache(timestampArg string) bool {
	if pool.MergeServiceCache[pool.IndexService{Timestamp: timestampArg, ServiceKey: int(storage.Service)}] != nil &&
		pool.MergeServiceCache[pool.IndexService{Timestamp: timestampArg, ServiceKey: int(storage.Pod)}] != nil {
		return true
	}
	return false
}

//
func handleService(serviceK8s modelK8s.ServiceList, pods modelK8s.PodList, timestamp string) {
	timestampI, _ := strconv.Atoi(timestamp)
	mergeServiceToDAOModel(serviceK8s, pods, timestampI)
}

//
func mergeServiceToDAOModel(serviceK8s modelK8s.ServiceList, pods modelK8s.PodList, timestamp int) {
	duplicationCheck := make(map[string]int)
	for _, services := range serviceK8s.Items {
		if duplicationCheck[services.Name] == 1 {
			continue
		}
		for k, v := range services.Spec.Selector {
			podNum, containerNum := searchPod(k, v, pods)
			if podNum == 0 {
				continue
			}
			servPreparation := types.ServiceDashboard{
				ServiceName:     services.Name,
				PodNumber:       podNum,
				ContainerNumber: containerNum,
				Timestamp:       timestamp,
			}
			if pool.ServiceDAOBuffer.Back() == nil {
				servRes := initAverage(servPreparation)
				pool.ServiceDAOBuffer.PushBack(servRes)
			} else {
				bottomElem := pool.ServiceDAOBuffer.Back()
				prevServ := getServiceInCache(servPreparation, bottomElem)
				servRes := calcServiceMovingAverage(servPreparation, prevServ)
				pool.ServiceDAOBuffer.PushBack(servRes)
				duplicationCheck[services.Name] = 1
			}
			break
		}

	}
}

//
func calcServiceMovingAverage(current types.ServiceDashboard, Prev types.ServiceDashboard) types.ServiceDashboard {
	containerCalc := calculate.CalcParam{
		CurrentVal:       current.ContainerNumber,
		PrevVal:          Prev.ContainerNumber,
		CurrentTimestamp: current.Timestamp,
		PrevTimestamp:    Prev.Timestamp,
	}
	current.AverageDayContainerNumber = containerCalc.GetDayAver()
	current.AverageHourContainerNumber = containerCalc.GetHourAver()
	current.AverageMinuteContainerNumber = containerCalc.GetMinuteAver()
	podCalc := calculate.CalcParam{
		CurrentVal:       current.PodNumber,
		PrevVal:          Prev.PodNumber,
		CurrentTimestamp: current.Timestamp,
		PrevTimestamp:    Prev.Timestamp,
	}
	current.AverageMinutePodNumber = podCalc.GetMinuteAver()
	current.AverageHourPodNumber = podCalc.GetHourAver()
	current.AverageDayPodNumber = podCalc.GetDayAver()
	return current
}

func getServiceInCache(currentService types.ServiceDashboard, prev *list.Element) types.ServiceDashboard {
	val := prev.Value.(types.ServiceDashboard)
	if val.ServiceName != currentService.ServiceName {
		prev = prev.Prev()
		if prev == nil {
			return initAverage(currentService)
		}
		val = getServiceInCache(currentService, prev)
	}
	return val
}

//
func initAverage(s types.ServiceDashboard) types.ServiceDashboard {
	s.AverageDayContainerNumber = s.ContainerNumber
	s.AverageHourContainerNumber = s.ContainerNumber
	s.AverageMinuteContainerNumber = s.ContainerNumber
	s.AverageMinutePodNumber = s.PodNumber
	s.AverageHourPodNumber = s.PodNumber
	s.AverageDayPodNumber = s.PodNumber
	return s
}

//
func searchPod(k string, v string, pods modelK8s.PodList) (int, int) {
	var c, p int
	for _, pod := range pods.Items {
		for labelKey, labelVal := range pod.Labels {
			if k == labelKey && v == labelVal {
				c = c + len(pod.Spec.Containers)
				p = p + 1
				break
			}
		}
	}
	return p, c
}
