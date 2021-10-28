{{- if .PreviousStatus}}
    {{- if ne .PreviousStatus .CurrentStatus}}
        **{{.CurrentStatus}}** - обновился статус (было *{{.PreviousStatus}}*)
    {{- else}}
        **{{.CurrentStatus}}** - статус выставлен повторно
    {{- end}}
{{- else}}
    Новый статус **{{.CurrentStatus}}**
{{- end}}
{{- if .CommentContent}}
    **{{.CommentAuthor}}** оставил комментарий:
    ```{{.CommentContent}}```
{{- end}}