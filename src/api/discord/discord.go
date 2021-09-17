package discord

import (
	"app/src/api/kitsu"
	"app/src/utils/config"
	"app/src/utils/request"
	"app/src/utils/truncate"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

// Discord embed API
// https://discord.com/developers/docs/resources/channel#embed-object-embed-structure

type EmbedAuthor struct {
	Name string `json:"name"`
}

type EmbedFooter struct {
	Text string `json:"text"`
}

type Embed struct {
	Color       int         `json:"color,omitempty"`
	Title       string      `json:"title"` // 256 characters
	Url         string      `json:"url,omitempty"`
	Description string      `json:"description"` // 4096 characters
	Author      EmbedAuthor `json:"author,omitempty"`
	Footer      EmbedFooter `json:"footer,omitempty"`
}
type Payload struct {
	Username  string  `json:"username,omitempty"`
	AvatarUrl string  `json:"avatar_url,omitempty"`
	Content   string  `json:"content"`
	Embeds    []Embed `json:"embeds,omitempty"`
}

type Assignee struct {
	Fullname string
	Email    string
	Phone    string
}

type Template struct {
	ProjectName    string
	GroupName      string
	ParentName     string
	TaskName       string
	TaskType       string
	SubTaskName    string
	CurrentStatus  string
	PreviousStatus string
	CommentContent string
	CommentAuthor  string
	EntityType     string
	Assignees      []Assignee
}

func SendMessageBunch(conf config.Config, data []kitsu.MessagePayload, webHookURL string) {
	payload := Payload{}
	payload.Content = ""

	for _, elem := range data {
		var placeholders Template

		placeholders.ProjectName = elem.Project.Name
		placeholders.GroupName = elem.EntityType.Name
		placeholders.ParentName = elem.Parent.Name
		placeholders.TaskName = elem.Entity.Name
		placeholders.TaskType = elem.TaskType.Name
		placeholders.CurrentStatus = elem.TaskStatus.ShortName
		placeholders.PreviousStatus = elem.PreviousStatusName
		placeholders.CommentContent = elem.LatestComment.Comment.Text
		placeholders.CommentAuthor = elem.LatestComment.Author.FullName
		placeholders.EntityType = elem.EntityType.EntityType.Name

		placeholders.Assignees = make([]Assignee, len(elem.Assignees))
		for i := 0; i < len(elem.Assignees); i++ {
			placeholders.Assignees[i].Fullname = elem.Assignees[i].FullName
			placeholders.Assignees[i].Email = elem.Assignees[i].Email
			placeholders.Assignees[i].Phone = elem.Assignees[i].Phone
		}

		hexColor := strings.ReplaceAll(elem.TaskStatus.Color, "#", "")
		intColor, err := strconv.ParseInt(hexColor, 16, 64)
		if err != nil {
			fmt.Printf("Conversion failed: %s\n", err)
		}

		// Title template
		author := parseTaskTemplate("tpl/author.tpl", placeholders)
		title := parseTaskTemplate("tpl/title.tpl", placeholders)
		description := parseTaskTemplate("tpl/description.tpl", placeholders)
		footer := parseTaskTemplate("tpl/footer.tpl", placeholders)

		embed := Embed{}
		embed.Color = int(intColor)
		embed.Author.Name = truncate.TruncateString(author, 256)
		embed.Title = truncate.TruncateString(title, 256)

		// Form URL with appropriate filtering
		url := "assets"
		args := elem.EntityType.EntityType.Name + "%20" + strings.Replace(elem.Entity.Name, "_", "%20", -1)
		if elem.EntityType.Name == "Shot" {
			url = "shots"
			args = elem.Parent.Name + "%20" + strings.Replace(elem.Entity.Name, "_", "%20", -1)
		}
		embed.Url = conf.Kitsu.Hostname + "productions/" + elem.Project.ID + "/" + url + "/task-types/" + elem.Task.TaskTypeID + "?search=" + args

		embed.Description = truncate.TruncateString(description, 4096)
		embed.Footer.Text = truncate.TruncateString(footer, 2048)

		payload.Embeds = append(payload.Embeds, embed)
	}

	request.Do("", http.MethodPost, webHookURL, payload, nil)
}

func parseTaskTemplate(tplFilePath string, data Template) string {
	tpl, err := ioutil.ReadFile(tplFilePath)
	if err != nil {
		fmt.Print(err)
	}

	// Create a new template and parse template file
	t := template.Must(template.New("template").Parse(string(tpl)))

	output := new(bytes.Buffer)
	t.Execute(output, data)
	if err != nil {
		log.Println("executing template:", err)
	}

	return output.String()
}
