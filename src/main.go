package main

import (
	"app/src/api/discord"
	"app/src/api/kitsu"
	"app/src/model"
	"app/src/utils/basicauth"
	"app/src/utils/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/beefsack/go-rate"

	"github.com/pieterclaerhout/go-waitgroup"
	"github.com/robfig/cron/v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func MakeKitsuResponse(conf config.Config) []kitsu.MessagePayload {

	tasks := kitsu.GetTasks()
	if conf.Log {
		fmt.Println("Got tasks: " + strconv.Itoa(len(tasks.Each)))
	}

	taskStatuses := kitsu.GetTaskStatuses()
	if conf.Log {
		fmt.Println("Got taskStatuses: " + strconv.Itoa(len(taskStatuses.Each)))
	}

	entities := kitsu.GetEntities()
	if conf.Log {
		fmt.Println("Got entities: " + strconv.Itoa(len(entities.Each)))
	}

	enitityTypes := kitsu.GetEntityTypes()
	if conf.Log {
		fmt.Println("Got enitityTypes: " + strconv.Itoa(len(enitityTypes.Each)))
	}

	projects := kitsu.GetProjects()
	if conf.Log {
		fmt.Println("Got projects: " + strconv.Itoa(len(projects.Each)))
	}

	taskTypes := kitsu.GetTaskTypes()
	if conf.Log {
		fmt.Println("Got taskTypes: " + strconv.Itoa(len(taskTypes.Each)))
	}

	persons := kitsu.GetPersons()
	if conf.Log {
		fmt.Println("Got persons: " + strconv.Itoa(len(persons.Each)))
	}

	var comments kitsu.Comments
	if conf.Kitsu.SkipComments == false {
		comments = kitsu.GetComments()
		if conf.Log {
			fmt.Println("Got comments: " + strconv.Itoa(len(comments.Each)))
		}
	}

	start := time.Now()

	response := make([]kitsu.MessagePayload, len(tasks.Each))

	wg := waitgroup.NewWaitGroup(conf.Threads)

	// tasks
	for i := 0; i < len(response); i++ {
		wg.BlockAdd()
		go func(i int) {
			defer wg.Done()

			// Ignore old messages
			layout := "2006-01-02T15:04:05"
			taskDate, err := time.Parse(layout, tasks.Each[i].UpdatedAt)
			if err != nil {
				fmt.Println(err)
			}
			daysCount := int(start.Sub(taskDate).Hours() / 24)

			if conf.IgnoreMessagesDaysOld != 0 && daysCount > conf.IgnoreMessagesDaysOld {
				return
			}

			// Store task
			response[i].Task.Task = tasks.Each[i]

			// Store task status
			for _, elem := range taskStatuses.Each {
				if elem.ID == tasks.Each[i].TaskStatusID {
					response[i].TaskStatus.TaskStatus = elem
					break
				}
			}

			// Store entity
			for _, elem := range entities.Each {
				if elem.ID == tasks.Each[i].EntityID {
					response[i].Entity.Entity = elem
					break
				}
			}

			// Store entity type
			for _, elem := range enitityTypes.Each {
				if elem.ID == response[i].Entity.Entity.EntityTypeID {
					response[i].EntityType.EntityType = elem
					break
				}
			}

			// Store parent
			for _, elem := range entities.Each {
				if elem.ID == response[i].Entity.Entity.ParentID {
					response[i].Parent.Entity = elem
				}
			}

			// Store project
			for _, elem := range projects.Each {
				if elem.ID == response[i].Entity.Entity.ProjectID {
					response[i].Project.Project = elem
					break
				}
			}

			// Store task type
			for _, elem := range taskTypes.Each {
				if elem.ID == tasks.Each[i].TaskTypeID {
					response[i].TaskType.TaskType = elem
					break
				}
			}

			// Store comments
			if conf.Kitsu.SkipComments == false {
				var taskComments kitsu.Comments
				for _, elem := range comments.Each {
					if elem.ObjectID == tasks.Each[i].ID {
						taskComments.Each = append(taskComments.Each, elem)
					}
				}

				if len(taskComments.Each) > 0 {
					// find the most recent comment in array
					sort.Slice(taskComments.Each, func(i, j int) bool {
						layout := "2006-01-02T15:04:05"
						a, err := time.Parse(layout, taskComments.Each[i].UpdatedAt)
						if err != nil {
							fmt.Println(err)
						}
						b, err := time.Parse(layout, taskComments.Each[j].UpdatedAt)
						if err != nil {
							fmt.Println(err)
						}
						return a.Unix() > b.Unix()
					})

					response[i].LatestComment.Comment.Comment = taskComments.Each[0]

				}

				// Store comment author
				for _, elem := range persons.Each {
					if len(taskComments.Each) > 0 {
						if elem.ID == taskComments.Each[0].PersonID {
							response[i].LatestComment.Author.Person = elem
							break
						}
					}
				}
			}

			// Store assignee
			if len(tasks.Each[i].Assignees) > 0 {
				for _, assigneeID := range tasks.Each[i].Assignees {
					for _, person := range persons.Each {
						if assigneeID == person.ID {
							response[i].Assignees = append(response[i].Assignees, person)
						}
					}
				}
			}

		}(i)
	}
	wg.Wait()

	if conf.Log {
		log.Printf("Done primary loop in %s", time.Since(start))
	}
	//return response

	// Remove empty elems
	var out []kitsu.MessagePayload
	for _, elem := range response {
		if len(elem.Task.Task.ID) > 0 {
			out = append(out, elem)
		}
	}

	if conf.Log {
		log.Printf("Done secondary loop in %s", time.Since(start))
	}

	return out
}

func DumpToFile(data []kitsu.MessagePayload, filename string) {

	file, _ := json.MarshalIndent(data, "", "    ")
	_ = ioutil.WriteFile("dump/"+filename+".json", file, 0644)

}

func DiscordQueue(data []kitsu.MessagePayload, conf config.Config, db *gorm.DB) []kitsu.MessagePayload {
	if len(data) == 0 {
		if conf.Log {
			fmt.Printf("Nothing to do\n")
		}
		return []kitsu.MessagePayload{}
	}

	rl := rate.New(conf.Discord.RequestsPerMinute, time.Minute) // 50 times per minute

	var payload []kitsu.MessagePayload
	var out []kitsu.MessagePayload
	for i := 0; i < len(data); i++ {

		dbResult := model.FindTask(db, data[i].Task.ID)

		// DB verify
		data[i].PreviousStatusName = dbResult.TaskStatus

		if len(dbResult.TaskID) > 0 {
			// check if status is different or last updated date don't match
			if dbResult.TaskStatus != data[i].TaskStatus.TaskStatus.ShortName || dbResult.TaskUpdatedAt != data[i].Task.Task.UpdatedAt {
				// update
				model.UpdateTask(db, data[i].Task.Task.ID, data[i].Task.Task.UpdatedAt, data[i].TaskStatus.TaskStatus.ShortName, data[i].LatestComment.Comment.ID, data[i].LatestComment.Comment.UpdatedAt)
			} else {
				continue
			}
		} else {
			// create
			model.CreateTask(db, data[i].Task.Task.ID, data[i].Task.Task.UpdatedAt, data[i].TaskStatus.TaskStatus.ShortName, data[i].LatestComment.Comment.ID, data[i].LatestComment.Comment.UpdatedAt)
		}

		if conf.SilentUpdateDB == true {
			if conf.Log {
				log.Printf("Ignoring message\n")
			}
			continue
		}

		payload = append(payload, data[i])
		out = append(out, payload...)

		if len(payload) >= conf.Discord.EmbedsPerRequests || i == len(data)-1 {
			if conf.Log {
				log.Printf("Sending bunch of messages: " + strconv.Itoa(len(payload)))
			}
			discord.SendMessageBunch(conf, payload, conf.Discord.WebhookURL)

			payload = nil

			rl.Wait()
		}
	}

	return out
}

func main() {
	start := time.Now()

	// Load config
	conf := config.Read()

	// Debug
	if conf.Debug {
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

	if conf.Log {
		log.Printf("Connected to database in %s", time.Since(start))

		if _, err := os.Stat("./dump"); os.IsNotExist(err) {
			err := os.Mkdir("./dump", os.ModeDir)
			if err != nil {
				panic("failed to create dir")
			}
		}
	}

	// Create Cron
	c := cron.New(cron.WithChain(
		cron.DelayIfStillRunning(cron.DefaultLogger),
	))

	// Kitsu auth
	token := basicauth.AuthForJWTToken(conf.Kitsu.Hostname+"api/auth/login", conf.Kitsu.Email, conf.Kitsu.Password)
	os.Setenv("KitsuJWTToken", token)
	if conf.Log {
		log.Printf("Connected to Kitsu in %s", time.Since(start))
	}

	c.AddFunc("@every 1h", func() {
		token := basicauth.AuthForJWTToken(conf.Kitsu.Hostname+"api/auth/login", conf.Kitsu.Email, conf.Kitsu.Password)
		os.Setenv("KitsuJWTToken", token)
		if conf.Log {
			fmt.Println("Got new Kitsu token")
		}
	})

	// Request data from Kitsu
	kitsuResponse := MakeKitsuResponse(conf)
	if conf.Log {
		DumpToFile(kitsuResponse, "kitsu_response")
		log.Printf("Done MakeKitsuResponse in %s", time.Since(start))
	}

	// Prepare messages to Discord
	discordPayload := DiscordQueue(kitsuResponse, conf, db)
	if conf.Log {
		DumpToFile(discordPayload, "discord_payload")
		log.Printf("Done DiscordQueue in %s", time.Since(start))
	}

	c.AddFunc("@every "+strconv.Itoa(conf.Kitsu.RequestInterval)+"m", func() {
		// Request data from Kitsu
		kitsuResponse := MakeKitsuResponse(conf)
		if conf.Log {
			DumpToFile(kitsuResponse, "kitsu_response")
			log.Printf("Done MakeKitsuResponse in %s", time.Since(start))
		}

		// Prepare messages to Discord
		DiscordQueue(kitsuResponse, conf, db)
		if conf.Log {
			log.Printf("Done DiscordQueue in %s", time.Since(start))
		}
	})

	c.Run()
}
