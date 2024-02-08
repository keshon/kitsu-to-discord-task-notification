package debug

import (
	"app/src/utils/config"
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"runtime"
	"time"

	"github.com/hokaccha/go-prettyjson"
)

func Info(resp *http.Response, respBody []byte) {

	if config.Read().Debug == true {
		// Func name and path
		pc := make([]uintptr, 10) // at least 1 entry needed
		runtime.Callers(2, pc)
		f := runtime.FuncForPC(pc[0])
		file, line := f.FileLine(pc[0])
		slog.Info("%s:%d %s\n", file, line, f.Name())

		// Headers
		prettyResp, _ := prettyjson.Marshal(resp)

		// Body
		var arrayMap []map[string]interface{}
		var objectMap map[string]interface{}
		var prettyBody []byte

		x := bytes.TrimLeft(respBody, " \t\r\n")
		isArray := len(x) > 0 && x[0] == '['
		isObject := len(x) > 0 && x[0] == '{'

		if isArray == true {
			err := json.Unmarshal(respBody, &arrayMap)
			if err != nil {
				panic(err)
			}
			prettyBody, _ = prettyjson.Marshal(arrayMap)
		}
		if isObject == true {
			err := json.Unmarshal(respBody, &objectMap)
			if err != nil {
				panic(err)
			}

			prettyBody, _ = prettyjson.Marshal(objectMap)
		}

		dt := time.Now()

		// Pretty print
		slog.Info("--start--")
		slog.Info(dt.String())
		slog.Info("Status : ", resp.StatusCode)
		slog.Info("Headers : ", string(prettyResp))
		slog.Info("Body : ", string(prettyBody))
		slog.Info("--end--\n")
	}
}
