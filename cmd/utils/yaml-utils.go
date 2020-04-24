package utils

import (
	"log"

	ypg "gopkg.in/yaml.v2"
)

//PodStatus ...
type PodStatus struct {
	podCount  int
	Pending   int
	Running   int
	Succeeded int
	Failed    int
	Unknown   int
}

func NewpodStatus() *PodStatus {
	return &PodStatus{
		podCount:  0,
		Pending:   0,
		Running:   0,
		Succeeded: 0,
		Failed:    0,
		Unknown:   0,
	}
}

//CountPods ...
func (podStatus *PodStatus) CountPods(yaml []byte) {
	//Make sure all counts are 0
	podStatus.podCount = 0
	podStatus.Pending = 0
	podStatus.Running = 0
	podStatus.Succeeded = 0
	podStatus.Failed = 0
	podStatus.Unknown = 0

	m := make(map[string]interface{})
	//mItems := make(map[string]interface{})

	err := ypg.Unmarshal([]byte(yaml), &m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	//fmt.Printf("--- m:\n%v\n\n", m["items"].(string))
	for key, itemMap := range m {
		if key == "items" {
			for key := range itemMap.([]interface{}) {
				kind := itemMap.([]interface{})[key].(map[interface{}]interface{})["kind"].(string)
				statusMap := itemMap.([]interface{})[key].(map[interface{}]interface{})["status"]
				if kind == "Pod" {
					for statusKey, phase := range statusMap.(map[interface{}]interface{}) {
						if statusKey == "phase" {
							if phase.(string) == "Running" {
								podStatus.Running++
								podStatus.podCount++
							} else if phase.(string) == "Succeeded" {
								podStatus.Succeeded++
								podStatus.podCount++
							} else if phase.(string) == "Pending" {
								podStatus.Pending++
								podStatus.podCount++
							} else if phase.(string) == "Failed" {
								podStatus.Failed++
								podStatus.podCount++
							} else if phase.(string) == "Unknown" {
								podStatus.Unknown++
								podStatus.podCount++
							} else {
								continue
							}

						}
					}
				}

			}
		}

	}
}
