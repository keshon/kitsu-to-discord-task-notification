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
	"text/template"
	"time"
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

func SendMessageBunch(conf config.Config, message []kitsu.TaskResponse, webHookURL string) {
	payload := Payload{}
	payload.Content = ""

	for _, elem := range message {
		// Title template
		author := parseTaskTemplate(conf.TemplatePath+"/author.tpl", elem)
		title := parseTaskTemplate(conf.TemplatePath+"/title.tpl", elem)
		description := parseTaskTemplate(conf.TemplatePath+"/description.tpl", elem)
		footer := parseTaskTemplate(conf.TemplatePath+"/footer.tpl", elem)

		embed := Embed{}
		embed.Author.Name = truncate.TruncateString(author, 256)
		embed.Title = truncate.TruncateString(title, 256)

		if conf.Kitsu.SkipProject != true {
			embed.Url = conf.Kitsu.Hostname + "productions/" + elem.ProjectID + "/assets?search=" + elem.TaskName
		}

		embed.Description = truncate.TruncateString(description, 4096)
		embed.Footer.Text = truncate.TruncateString(footer, 2048)

		payload.Embeds = append(payload.Embeds, embed)
	}

	request.Do("", http.MethodPost, webHookURL, payload, nil)
	time.Sleep(time.Duration(conf.PostDelay) * time.Second)
}

func SendMessage(conf config.Config, message kitsu.TaskResponse, webHookURL string) {
	// Title template
	author := parseTaskTemplate(conf.TemplatePath+"/author.tpl", message)
	title := parseTaskTemplate(conf.TemplatePath+"/title.tpl", message)
	description := parseTaskTemplate(conf.TemplatePath+"/description.tpl", message)
	footer := parseTaskTemplate(conf.TemplatePath+"/footer.tpl", message)

	embed := Embed{}
	embed.Author.Name = author
	embed.Title = title

	if conf.Kitsu.SkipProject != true && message.TaskType != "" {
		typeURL := "assets"
		if message.TaskType == "Shot" {
			typeURL = "shots"
		}
		embed.Url = conf.Kitsu.Hostname + "productions/" + message.ProjectID + "/" + typeURL + "?search=" + message.TaskName
	}

	embed.Description = description
	embed.Footer.Text = footer

	payload := Payload{}
	payload.Content = ""
	payload.Embeds = append(payload.Embeds, embed)

	request.Do("", http.MethodPost, webHookURL, payload, nil)
	time.Sleep(time.Duration(conf.PostDelay) * time.Second)
}

func parseTaskTemplate(tplFilePath string, data kitsu.TaskResponse) string {
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
