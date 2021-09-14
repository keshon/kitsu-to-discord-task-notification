// Package kitsu provides methods for Kitsu task management software
package kitsu

import (
	"app/src/utils/config"
	"app/src/utils/request"
	"app/src/utils/sanitize"
	"app/src/utils/truncate"
	"fmt"
	"net/http"
	"os"
	"sort"
	"time"
)

type Task struct {
	Assignees       []string    `json:"assignees,omitempty"`
	ID              string      `json:"id,omitempty"`
	CreatedAt       string      `json:"created_at,omitempty"`
	UpdatedAt       string      `json:"updated_at,omitempty"`
	Name            string      `json:"name,omitempty"`
	LastCommentDate string      `json:"last_comment_date,omitempty"`
	Data            interface{} `json:"data,omitempty"`
	ProjectID       string      `json:"project_id,omitempty"`
	TaskTypeID      string      `json:"task_type_id,omitempty"`
	TaskStatusID    string      `json:"task_status_id,omitempty"`
	EntityID        string      `json:"entity_id,omitempty"`
	AssignerID      string      `json:"assigner_id,omitempty"`
	Type            string      `json:"type,omitempty"`
}
type Tasks struct {
	Each []Task
}

type Person struct {
	ID                        string `json:"id,omitempty"`
	CreatedAt                 string `json:"created_at,omitempty"`
	UpdatedAt                 string `json:"updated_at,omitempty"`
	FirstName                 string `json:"first_name,omitempty"`
	LastName                  string `json:"last_name,omitempty"`
	Email                     string `json:"email,omitempty"`
	Phone                     string `json:"phone,omitempty"`
	Active                    bool   `json:"active,omitempty"`
	LastPresence              string `json:"last_presence,omitempty"`
	DesktopLogin              string `json:"desktop_login,omitempty"`
	ShotgunID                 string `json:"shotgun_id,omitempty"`
	Timezone                  string `json:"timezone,omitempty"`
	Locale                    string `json:"locale,omitempty"`
	Data                      string `json:"data,omitempty"`
	Role                      string `json:"role,omitempty"`
	HasAvatar                 bool   `json:"has_avatar,omitempty"`
	NotificationsEnabled      bool   `json:"notifications_enabled,omitempty"`
	NotificationsSlackEnabled bool   `json:"notifications_slack_enabled,omitempty"`
	NotificationsSlackUserid  string `json:"notifications_slack_userid,omitempty"`
	Type                      string `json:"type,omitempty"`
	FullName                  string `json:"full_name,omitempty"`
}

type Entity struct {
	EntitiesOut     []interface{} `json:"entities_out,omitempty"`
	InstanceCasting []interface{} `json:"instance_casting,omitempty"`
	CreatedAt       string        `json:"created_at,omitempty"`
	UpdatedAt       string        `json:"updated_at,omitempty"`
	ID              string        `json:"id,omitempty"`
	Name            string        `json:"name,omitempty"`
	Code            interface{}   `json:"code,omitempty"`
	Description     interface{}   `json:"description,omitempty"`
	ShotgunID       interface{}   `json:"shotgun_id,omitempty"`
	Canceled        bool          `json:"canceled,omitempty"`
	NbFrames        interface{}   `json:"nb_frames,omitempty"`
	ProjectID       string        `json:"project_id,omitempty"`
	EntityTypeID    string        `json:"entity_type_id,omitempty"`
	ParentID        string        `json:"parent_id,omitempty"`
	SourceID        interface{}   `json:"source_id,omitempty"`
	PreviewFileID   interface{}   `json:"preview_file_id,omitempty"`
	Data            interface{}   `json:"data,omitempty"`
	EntitiesIn      []interface{} `json:"entities_in,omitempty"`
	Type            string        `json:"type,omitempty"`
}

type Entities struct {
	Each []Entity
}

type EntityType struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type TaskStatus struct {
	ID              string      `json:"id,omitempty"`
	CreatedAt       string      `json:"created_at,omitempty"`
	UpdatedAt       string      `json:"updated_at,omitempty"`
	Name            string      `json:"name,omitempty"`
	ShortName       string      `json:"short_name,omitempty"`
	Color           string      `json:"color,omitempty"`
	IsDone          bool        `json:"is_done,omitempty"`
	IsArtistAllowed bool        `json:"is_artist_allowed,omitempty"`
	IsClientAllowed bool        `json:"is_client_allowed,omitempty"`
	IsRetake        bool        `json:"is_retake,omitempty"`
	ShotgunID       interface{} `json:"shotgun_id,omitempty"`
	IsReviewable    bool        `json:"is_reviewable,omitempty"`
	Type            string      `json:"type,omitempty"`
}

type Comment struct {
	ID        string      `json:"id,omitempty"`
	CreatedAt string      `json:"created_at,omitempty"`
	UpdatedAt string      `json:"updated_at,omitempty"`
	ShotgunID interface{} `json:"shotgun_id,omitempty"`
	ObjectID  string      `json:"object_id,omitempty"`
	PersonID  string      `json:"person_id,omitempty"`
	Text      string      `json:"text,omitempty"`
}

type Comments struct {
	Each []Comment
}

type TaskType struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	ShortName string `json:"short_name,omitempty"`
}

type TaskTypes struct {
	Each []TaskType
}

type Project struct {
	ID              string `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	ProjectStatusID string `json:"project_status_id,omitempty"`
}

type Projects struct {
	Each []Project
}

type ProjectStatus struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type ProjectStatuses struct {
	Each []ProjectStatus
}

func GetComment(objectID string) Comments {
	path := config.Read().Kitsu.Hostname + "api/data/comments?object_id=" + objectID
	response := Comments{}
	request.Do(os.Getenv("KitsuJWTToken"), http.MethodGet, path, nil, &response.Each)

	return response
}

func GetTasks() Tasks {
	path := config.Read().Kitsu.Hostname + "api/data/tasks/"
	response := Tasks{}
	println(os.Getenv("KitsuJWTToken"))
	request.Do(os.Getenv("KitsuJWTToken"), http.MethodGet, path, nil, &response.Each)

	return response
}

func GetTask(taskID string) Task {
	path := config.Read().Kitsu.Hostname + "api/data/tasks/" + taskID
	response := Task{}
	request.Do(os.Getenv("KitsuJWTToken"), http.MethodGet, path, nil, &response)

	return response
}

func GetPerson(personID string) Person {
	path := config.Read().Kitsu.Hostname + "api/data/persons/" + personID
	response := Person{}
	request.Do(os.Getenv("KitsuJWTToken"), http.MethodGet, path, nil, &response)

	return response
}

func GetEntities(EntityID string) Entities {
	path := config.Read().Kitsu.Hostname + "api/data/entities/"
	response := Entities{}
	request.Do(os.Getenv("KitsuJWTToken"), http.MethodGet, path, nil, &response.Each)

	return response
}

func GetEntity(EntityID string) Entity {
	path := config.Read().Kitsu.Hostname + "api/data/entities/" + EntityID
	response := Entity{}
	request.Do(os.Getenv("KitsuJWTToken"), http.MethodGet, path, nil, &response)

	return response
}

func GetEntityType(entityTypeID string) EntityType {
	path := config.Read().Kitsu.Hostname + "api/data/entity-types/" + entityTypeID
	response := EntityType{}
	request.Do(os.Getenv("KitsuJWTToken"), http.MethodGet, path, nil, &response)

	return response
}

func GetTaskStatus(taskStatusID string) TaskStatus {
	path := config.Read().Kitsu.Hostname + "api/data/task-status/" + taskStatusID
	response := TaskStatus{}
	request.Do(os.Getenv("KitsuJWTToken"), http.MethodGet, path, nil, &response)

	return response
}

func GetTaskType(taskID string) TaskType {
	path := config.Read().Kitsu.Hostname + "api/data/task-types/" + taskID
	response := TaskType{}
	request.Do(os.Getenv("KitsuJWTToken"), http.MethodGet, path, nil, &response)

	return response
}

func GetProject(projectID string) Project {
	path := config.Read().Kitsu.Hostname + "api/data/projects/" + projectID
	response := Project{}
	request.Do(os.Getenv("KitsuJWTToken"), http.MethodGet, path, nil, &response)

	return response
}

func GetProjectStatus(projectStatusID string) ProjectStatus {
	path := config.Read().Kitsu.Hostname + "api/data/project-status/" + projectStatusID
	response := ProjectStatus{}
	request.Do(os.Getenv("KitsuJWTToken"), http.MethodGet, path, nil, &response)

	return response
}

type TaskResponse struct {
	ProjectName      string
	ProjectID        string
	TaskName         string
	TaskUpdatedAt    string
	TaskType         string
	SubTaskName      string
	StatusName       string
	OldStatusName    string
	AssigneesList    string
	CommentID        string
	CommentMessage   string
	CommentAuthor    string
	CommentUpdatedAt string
}

func MakeTaskResponse(conf config.Config, task Task) TaskResponse {
	response := TaskResponse{}

	// Get human readable status
	response.StatusName = GetTaskStatus(task.TaskStatusID).ShortName

	// Get entity name (Top Task)
	entity := GetEntity(task.EntityID)
	response.TaskName = entity.Name
	response.TaskUpdatedAt = entity.UpdatedAt
	entityType := GetEntityType(entity.EntityTypeID)
	response.TaskType = entityType.Name

	// Parse project name (production)
	if conf.Kitsu.SkipProject == false {
		project := GetProject(entity.ProjectID)

		if project.Name != "" {
			response.ProjectName = sanitize.Sanitize(project.Name)
			response.ProjectID = project.ID
		}
	}

	// Get subtask name (task type)
	taskType := GetTaskType(task.TaskTypeID)
	if taskType.Name != "" {
		response.SubTaskName = sanitize.Sanitize(taskType.Name)
	}

	// Get assingee for the Task and his phone data (we store Telegram nicknames there)
	detailedTask := GetTask(task.ID)
	if len(detailedTask.Assignees) > 0 && conf.Kitsu.SkipMentions != true {
		for key, elem := range detailedTask.Assignees {
			assingnee := GetPerson(elem)
			if assingnee.FullName != "" {
				response.AssigneesList += assingnee.FirstName + " " + assingnee.LastName
				if len(detailedTask.Assignees) > 1 && len(detailedTask.Assignees)-1 == key {
					response.AssigneesList += ", "
				}
			}
		}
	}

	// Get comment
	if conf.Kitsu.SkipComments != true {
		comments := GetComment(detailedTask.ID)
		if len(comments.Each) > 0 {
			// find the most recent comment in array
			sort.Slice(comments.Each, func(i, j int) bool {
				layout := "2006-01-02T15:04:05"
				a, err := time.Parse(layout, comments.Each[i].UpdatedAt)
				if err != nil {
					fmt.Println(err)
				}
				b, err := time.Parse(layout, comments.Each[j].UpdatedAt)
				if err != nil {
					fmt.Println(err)
				}
				return a.Unix() > b.Unix()
			})

			response.CommentID = comments.Each[0].ID
			response.CommentUpdatedAt = comments.Each[0].UpdatedAt

			if comments.Each[0].Text != "" {
				commentMessage := truncate.TruncateString(comments.Each[0].Text, 128)
				if commentMessage != comments.Each[0].Text {
					commentMessage += "..."
				}
				response.CommentMessage = commentMessage
				response.CommentAuthor = GetPerson(comments.Each[0].PersonID).FullName
			}
		}
	}

	return response
}
