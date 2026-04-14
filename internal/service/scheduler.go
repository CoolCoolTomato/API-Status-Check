package service

import (
	"log"
	"sync"
	"time"
)

type Scheduler struct {
	checkService *CheckService
	mu           sync.Mutex
	running      bool
}

func NewScheduler(checkService *CheckService) *Scheduler {
	return &Scheduler{checkService: checkService}
}

func (s *Scheduler) Start() {
	log.Println("Starting scheduler...")
	go s.run()
}

func (s *Scheduler) run() {
	s.checkService.RunCheck()

	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		if s.running {
			s.mu.Unlock()
			log.Println("Previous check still running, skipping...")
			continue
		}
		s.running = true
		s.mu.Unlock()

		log.Println("Running scheduled check...")
		s.checkService.RunCheck()

		s.mu.Lock()
		s.running = false
		s.mu.Unlock()
	}
}
