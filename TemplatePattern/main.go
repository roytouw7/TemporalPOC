package main

func main() {
	var workflow BusinessWorkflow
	workflow = &TemporalWorkflow{}
	workflow.Execute(workflow)
}
