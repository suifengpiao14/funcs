package funcs_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/spf13/cast"
	"github.com/suifengpiao14/funcs"
)

type TrafficAllocation struct {
	Name   string `json:"name"`
	Weight string `json:"weight"`
}
type TrafficAllocations []TrafficAllocation

func (t TrafficAllocation) GetWeight() int {
	return cast.ToInt(t.Weight)
}

type Result []string

func (r Result) Times() (times map[string]int) {
	times = make(map[string]int)
	for _, v := range r {
		if _, ok := times[v]; !ok {
			times[v] = 0
		}
		times[v]++
	}
	return times
}

func (r Result) Percent() (p map[string]string) {
	times := r.Times()
	count := len(r)
	p = make(map[string]string)
	for k, v := range times {
		rate := float64(v) / float64(count) * 100
		p[k] = fmt.Sprintf("%.2f", rate)
	}
	return p
}

func TestNewWeightedRandomSelections(t *testing.T) {
	trafficAllocations := TrafficAllocations{
		{Name: "创蓝", Weight: "50"},
		{Name: "玄武", Weight: "50"},
	}
	weightedRandomSelections := funcs.NewWeightedRandomSelections(trafficAllocations...)
	allResult := make(Result, 0)
	g := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		g.Add(1)
		go func() {
			defer g.Done()
			for i := 0; i < 100000; i++ {
				server := weightedRandomSelections.Select()
				allResult = append(allResult, server.Name)
			}
		}()
	}
	g.Wait()
	p := allResult.Percent()
	fmt.Println(p)
}
