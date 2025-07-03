package scheduler

import (
	"log"
	"sync"
	"time"
	"todolist/database"
	"todolist/model"

	"github.com/go-co-op/gocron/v2"
)

func UpdateStatusFromDueDate() {
	// scheduler due date
	s, err := gocron.NewScheduler()
	if err != nil {
		log.Fatal(err)
	}
	j, err := s.NewJob(
		gocron.DurationJob(
			1*time.Second,
		),
		gocron.NewTask(
			func() {
				now := time.Now()
				// log.Println(now.Format("2006-01-02 15:04:05"))

				// list all tasks
				var todos []model.Todo
				find := database.Postgres.Where("due_date <= ? AND status = 'pending'", now).Find(&todos)
				if find.Error != nil {
					log.Fatal(find.Error)
				}

				var wg sync.WaitGroup
				wg.Add(len(todos))
				for i := 1; i <= len(todos); i++ {
					go worker(&wg, todos[i-1])
				}
				wg.Wait() // tunggu semua goroutine selesai
			},
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Job ID: %s\n", j.ID())
	s.Start()
}

func worker(wg *sync.WaitGroup, todo model.Todo) {
	defer wg.Done()

	// update status menjadi completed
	todo.Status = "completed"
	update := database.Postgres.Save(&todo)
	if update.Error != nil {
		log.Fatal(update.Error)
	}
	log.Printf("Task %s updated successfully\n", todo.Title)
}
