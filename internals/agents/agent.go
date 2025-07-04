package agents

import (
	"context"
	"sync"
)

// When we get a response from OpenAPI we want to
type FileTask struct {
	Path    string
	Content string
}

type Agent struct {
	openAI           *OpenAPI
	outputDir        string
	basePackage      string
	taskQueue        chan FileTask
	wg               sync.WaitGroup
	workerCount      int
	cntxt            context.Context
	cancel           context.CancelFunc
	fileWrittenMutex sync.Mutex
	filesWritten     map[string]bool
}

func NewAgent(cntxt context.Context, openAPI *OpenAPI, outputDir string, basePackage string, workerCount int) *Agent {
	// context.Context, context.CancelFunc
	cntxt, cancel := context.WithCancel(cntxt)
	return &Agent{
		openAI:       openAPI,
		outputDir:    outputDir,
		basePackage:  basePackage,
		taskQueue:    make(chan FileTask, 100),
		workerCount:  workerCount,
		cntxt:        cntxt,
		cancel:       cancel,
		filesWritten: make(map[string]bool),
	}
}
