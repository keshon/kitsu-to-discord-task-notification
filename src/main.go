package main

import (
	"app/src/api/discord"
	"app/src/api/kitsu"
	"app/src/model"
	"app/src/utils/basicauth"
	"app/src/utils/config"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/robfig/cron/v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func mainTask(conf config.Config, db *gorm.DB) {
	// Get all Tasks
	tasks := kitsu.GetTasks()

	if len(tasks.Each) <= 0 {
		return
	}
	fmt.Println(strconv.Itoa(len(tasks.Each)) + " records found")

	switch threads := conf.Threads; {
	case threads < 0:
		fmt.Println("Thread mode is: Wait groups")

		var wg sync.WaitGroup
		wg.Add(len(tasks.Each))

		for _, elem := range tasks.Each {
			go func(elem kitsu.Task) {
				defer wg.Done()

				kitsu.MakeTaskResponse(conf, elem)

			}(elem)
		}
		wg.Wait()

	case threads > 0:
		fmt.Println("Thread mode is: Semafore")

		// Semafore async
		var sem = make(chan int, threads)

		for _, elem := range tasks.Each {
			sem <- 1
			go func() {

				kitsu.MakeTaskResponse(conf, elem)

				<-sem
			}()
		}

	default:
		fmt.Println("Thread mode is: Synced")

		for key, task := range tasks.Each {
			// Check DB first
			dbResult := model.FindTask(db, task.ID)

			if conf.Log == true {
				fmt.Println("=====================")
				fmt.Println(strconv.Itoa(key) + " record \n")

				fmt.Println("**DB data**")
				fmt.Println("ID: " + strconv.FormatUint(uint64(dbResult.ID), 10))
				fmt.Println("TaskID: " + dbResult.TaskID)
				fmt.Println("TaskUpdatedAt: " + dbResult.TaskUpdatedAt)
				fmt.Println("CommentID: " + dbResult.CommentID)
				fmt.Println("CommentUpdatedAt: " + dbResult.CommentUpdatedAt + "\n")
			}

			// Ignore DONE unchanged tasks
			if len(dbResult.TaskID) > 0 {
				for _, elem := range conf.Kitsu.IsDoneStatusNames {
					if dbResult.TaskStatus == elem && dbResult.TaskUpdatedAt == task.UpdatedAt {
						continue
					}
				}
			}

			// Get detailed reposponse for particular task
			taskResponse := kitsu.TaskResponse{}

			taskResponse = kitsu.MakeTaskResponse(conf, task)

			// Store current status from DB
			taskResponse.OldStatusName = dbResult.TaskStatus

			if len(dbResult.TaskID) > 0 {
				// check if status is different or last updated task date don't match
				if dbResult.TaskStatus != taskResponse.StatusName || dbResult.TaskUpdatedAt != taskResponse.TaskUpdatedAt {
					// update
					model.UpdateTask(db, task.ID, taskResponse.TaskUpdatedAt, taskResponse.StatusName, taskResponse.CommentID, taskResponse.CommentUpdatedAt)
				} else {
					continue
				}
			} else {
				// create
				model.CreateTask(db, task.ID, taskResponse.TaskUpdatedAt, taskResponse.StatusName, taskResponse.CommentID, taskResponse.CommentUpdatedAt)

			}

			if conf.Log == true {
				fmt.Println("**Kitsu data**")
				fmt.Println("ProjectID: " + taskResponse.ProjectID)
				fmt.Println("ProjectName: " + taskResponse.ProjectName)
				fmt.Println("TaskUpdatedAt: " + taskResponse.TaskUpdatedAt)
				fmt.Println("TaskName: " + taskResponse.TaskName)
				fmt.Println("SubTaskName: " + taskResponse.SubTaskName)
				fmt.Println("StatusName: " + taskResponse.StatusName)
				fmt.Println("OldStatusName: " + taskResponse.OldStatusName)
				fmt.Println("CommentID: " + taskResponse.CommentID)
				fmt.Println("CommentAuthor: " + taskResponse.CommentAuthor)
				fmt.Println("CommentMessage: " + taskResponse.CommentMessage)
				fmt.Println("CommentUpdatedAt: " + taskResponse.CommentUpdatedAt + "\n")

				fmt.Println("=====================")
			}

			if conf.SilentUpdate != true {

				webhookURLsByStatus := conf.Discord.WebhookURLsByStatus
				messageSent := false

				if len(webhookURLsByStatus) > 0 {
					for _, elem := range webhookURLsByStatus {

						kitsuStatus := strings.Split(elem, ":")[0]
						discordURL := strings.Split(elem, ":")[1]

						if strings.ToLower(taskResponse.StatusName) == kitsuStatus {
							discord.SendMessage(conf, taskResponse, discordURL)

							if conf.Log == true {
								fmt.Println("Send message to webhookURLsByStatus for" + kitsuStatus)
							}

							messageSent = true
						}

					}
				}

				if conf.SuppressUndefinedRoles != true && messageSent == false {
					discord.SendMessage(conf, taskResponse, conf.Discord.WebhookURL)

					if conf.Log == true {
						fmt.Println("Send message to webhookURL")
					}
				}
			}
		}

		fmt.Println("Done")
	}
}

func main() {
	// Load config
	conf := config.Read()

	// Debug
	if conf.Debug == true {
		os.Setenv("Debug", "true")
	}

	// Connect to DB
	db, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&model.Task{})

	// Create Cron
	c := cron.New(cron.WithChain(
		cron.DelayIfStillRunning(cron.DefaultLogger),
	))

	// Kitsu auth
	token := basicauth.AuthForJWTToken(conf.Kitsu.Hostname+"api/auth/login", conf.Kitsu.Email, conf.Kitsu.Password)
	os.Setenv("KitsuJWTToken", token)

	fmt.Println("Got new token")
	c.AddFunc("@every 1h", func() {
		token := basicauth.AuthForJWTToken(conf.Kitsu.Hostname+"api/auth/login", conf.Kitsu.Email, conf.Kitsu.Password)
		os.Setenv("KitsuJWTToken", token)

		fmt.Println("Got new token")
	})

	mainTask(conf, db)
	c.AddFunc("@every "+strconv.Itoa(conf.PollInterval)+"s", func() {
		mainTask(conf, db)
	})

	c.Run()
}
