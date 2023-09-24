package main

import "fmt"

type BusinessWorkflow interface {
	Execute(workflow BusinessWorkflow)
	InvariantLogic([]int) int
	TemporalSpecificLogic()
	FetchResults()
}

type BaseWorkflow struct{}

func (bw *BaseWorkflow) Execute(workflow BusinessWorkflow) {
	nums := []int{1, 2, 3, 4, 5}
	sum := bw.InvariantLogic(nums)
	fmt.Printf("Sum calculated: %d\n", sum)
	workflow.TemporalSpecificLogic()
	workflow.FetchResults()
}

func (bw *BaseWorkflow) InvariantLogic(nums []int) int {
	sum := 0
	for _, num := range nums {
		sum += num
	}
	return sum
}

// TemporalSpecificLogic noop placeholder
func (bw *BaseWorkflow) TemporalSpecificLogic() {
	return
}

// FetchResults noop placeholder
func (bw *BaseWorkflow) FetchResults() {
	return
}
