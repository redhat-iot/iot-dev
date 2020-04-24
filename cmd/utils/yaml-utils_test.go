package utils

import (
	"io/ioutil"
	"testing"
)

func TestCountPods(t *testing.T) {
	file, _ := ioutil.ReadFile("yaml-utils-test.txt")
	podStatus := NewpodStatus()
	podStatus.CountPods(file)

	if podStatus.Running != 9 {
		t.Errorf("Failed")
	} else if podStatus.Succeeded != 3 {
		t.Errorf("Failed")
	}

}
