{{- if .PreviousStatus}}
    {{- if ne .PreviousStatus .CurrentStatus}}
        Status changed to **{{.CurrentStatus}}** (was *{{.PreviousStatus}}*)
    {{- else}}
        Status repeated - **{{.CurrentStatus}}**
    {{- end}}
{{- else}}
    New status **{{.CurrentStatus}}**
{{- end}}
{{- if .CommentContent}}
    New comment from **{{.CommentAuthor}}**:
    ```{{.CommentContent}}```
{{- end}}