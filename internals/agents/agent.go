package agents

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
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

func (a *Agent) Start() {
	log.Printf("Starting %d workers.....\n", a.workerCount)

	for i := 0; i < a.workerCount; i++ {
		a.wg.Add(1)
		go a.worker(i)
	}
}

func (a *Agent) worker(id int) {
	defer a.wg.Done()
	log.Printf("Worker %d started\n", id)
	for {
		select {
		case task, ok := <-a.taskQueue:
			if ok == false {
				log.Printf("Worker %d stopping\n", id)
				return
			}
			a.fileWrittenMutex.Lock()
			if a.filesWritten[task.Path] {
				log.Printf("File %s already written, skipping\n", task.Path)
				a.fileWrittenMutex.Unlock()
				continue
			}
			a.filesWritten[task.Path] = true
			a.fileWrittenMutex.Unlock()
			err := a.writeFile(task)
			if err != nil {
				log.Printf("error writing file ( worker %d) %s: %v\n", id, task.Path, err)
			} else {
				log.Printf("Worker %d wrote file %s\n", id, task.Path)
			}
		case <-a.cntxt.Done():
			log.Printf("Worker %d, received cancel signal", id)
			return
		}
	}
}

func (a *Agent) writeFile(task FileTask) error {
	fullPath := filepath.Join(a.outputDir, task.Path)
	dir := filepath.Dir(fullPath)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("Error creating directories: %s : %w\n", dir, err)
	}
	err = os.WriteFile(fullPath, []byte(task.Content), 0644)

	if err != nil {
		return fmt.Errorf("Error writing file %s: %w", fullPath, err)
	}
	log.Printf("Sucessfully wrote file: %s", fullPath)
	return nil
}

func (a *Agent) SendFileTask(path string, content string) {
	task := FileTask{
		Path:    path,
		Content: content,
	}

	go func() {
		a.taskQueue <- task
	}()
}

func (a *Agent) Stop() {
	log.Printf("Stopping Agent")
	close(a.taskQueue)
	a.cancel()
	a.wg.Wait()
}
