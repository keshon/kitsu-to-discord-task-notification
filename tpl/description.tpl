{{- if .OldStatusName}}
    {{- if ne .OldStatusName .StatusName}}
        Обновился статус на **{{.StatusName}}** (было *{{.OldStatusName}}*)
    {{- end}}
{{- else}}
    Новый статус **{{.StatusName}}**
{{- end}}
{{- if .CommentMessage}}
    Новый комментарий от **{{.CommentAuthor}}**:
    ```{{.CommentMessage}}```
{{- end}}