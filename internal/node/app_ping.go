package node

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

func (app *App) Ping(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	goroutineCount := runtime.NumGoroutine()
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	bs, _ := json.Marshal(app.stat())

	lines := []string{
		"BUILT HASH:  https://github.com/unchainese/unchain/tree/" + app.cfg.GitHash,
		"BUILT TIME:  " + app.cfg.BuildTime,
		"RUN_AT:     " + app.cfg.RunAt,
		fmt.Sprintf("GOROUTINE: %d", goroutineCount),
		fmt.Sprintf("MEMORY.Alloc:    %.2fMB", float64(memStats.Alloc)/1024/1024),
		fmt.Sprintf("MEMORY.TotalAlloc:    %.2fMB", float64(memStats.TotalAlloc)/1024/1024),
		"\n\n",
		string(bs),
	}
	w.Write([]byte(strings.Join(lines, "\n\n")))
}
