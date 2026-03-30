package gui

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/niuxh/sts2-go-agent/pkg/agent"
	"github.com/niuxh/sts2-go-agent/pkg/api"
	"github.com/niuxh/sts2-go-agent/pkg/config"
)

var (
	currentCancel context.CancelFunc
	logBuffer     []string
	logMu         sync.Mutex
)

func Launch() {
	// 静态页面 HTML
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, indexHTML)
	})

	// 获取配置接口
	http.HandleFunc("/api/config", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(config.Load())
	})

	// 启动 Agent 接口
	http.HandleFunc("/api/start", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost { return }
		
		var newCfg config.AppConfig
		json.NewDecoder(r.Body).Decode(&newCfg)
		newCfg.Save()

		// 停止旧的
		if currentCancel != nil { currentCancel() }

		logMu.Lock()
		logBuffer = []string{"[System] Initializing Agent...\n"}
		logMu.Unlock()

		gameClient := api.NewClient(newCfg.GameURL)
		llmProvider := agent.NewOpenAIProvider(newCfg.APIKey, newCfg.BaseURL, newCfg.Model)
		
		stsAgent := agent.NewAgent(gameClient, llmProvider, func(msg string) {
			logMu.Lock()
			logBuffer = append(logBuffer, msg)
			if len(logBuffer) > 200 { logBuffer = logBuffer[1:] }
			logMu.Unlock()
		})

		var ctx context.Context
		ctx, currentCancel = context.WithCancel(context.Background())
		go stsAgent.Run(ctx)
		
		w.WriteHeader(http.StatusOK)
	})

	// 停止接口
	http.HandleFunc("/api/stop", func(w http.ResponseWriter, r *http.Request) {
		if currentCancel != nil {
			currentCancel()
			currentCancel = nil
		}
		w.WriteHeader(http.StatusOK)
	})

	// 日志流接口 (SSE)
	http.HandleFunc("/api/logs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		for {
			logMu.Lock()
			if len(logBuffer) > 0 {
				data, _ := json.Marshal(logBuffer)
				fmt.Fprintf(w, "data: %s\n\n", string(data))
				w.(http.Flusher).Flush()
			}
			logMu.Unlock()
			time.Sleep(1 * time.Second)
		}
	})

	fmt.Println("-------------------------------------------")
	fmt.Println("STS2 AI Agent 控制中心已启动！")
	fmt.Println("请在浏览器中打开: http://localhost:8090")
	fmt.Println("-------------------------------------------")

	http.ListenAndServe(":8090", nil)
}

const indexHTML = `
<!DOCTYPE html>
<html>
<head>
    <title>STS2 AI Agent Control Center</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif; background: #121212; color: #e0e0e0; margin: 0; display: flex; flex-direction: column; height: 100vh; }
        .container { max-width: 900px; margin: 20px auto; padding: 20px; background: #1e1e1e; border-radius: 8px; box-shadow: 0 4px 6px rgba(0,0,0,0.3); }
        .header { text-align: center; margin-bottom: 20px; color: #4CAF50; }
        .form-group { margin-bottom: 15px; }
        label { display: block; margin-bottom: 5px; font-weight: bold; }
        input, select { width: 100%; padding: 10px; background: #2c2c2c; border: 1px solid #444; color: white; border-radius: 4px; box-sizing: border-box; }
        .controls { display: flex; gap: 10px; margin-top: 20px; }
        button { flex: 1; padding: 12px; border: none; border-radius: 4px; cursor: pointer; font-weight: bold; font-size: 16px; }
        .btn-start { background: #4CAF50; color: white; }
        .btn-stop { background: #f44336; color: white; }
        #logs { margin-top: 20px; background: #000; padding: 15px; height: 350px; overflow-y: auto; border-radius: 4px; font-family: 'Courier New', Courier, monospace; font-size: 14px; line-height: 1.4; border: 1px solid #333; }
        .status { margin-left: 10px; font-size: 14px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header"><h1>Slay the Spire 2 AI Agent</h1></div>
        
        <div class="form-group">
            <label>LLM API Key (OpenAI/DeepSeek)</label>
            <input type="password" id="apiKey" placeholder="sk-xxxx...">
        </div>
        <div class="form-group">
            <label>Base URL</label>
            <input type="text" id="baseURL" value="https://api.openai.com/v1">
        </div>
        <div class="form-group">
            <label>Model Name</label>
            <input type="text" id="model" value="gpt-4o">
        </div>
        <div class="form-group">
            <label>Game API URL</label>
            <input type="text" id="gameURL" value="http://127.0.0.1:8080">
        </div>

        <div class="controls">
            <button class="btn-start" onclick="startAgent()">启动 AI Agent</button>
            <button class="btn-stop" onclick="stopAgent()">停止</button>
        </div>

        <div id="logs"></div>
    </div>

    <script>
        // 加载配置
        fetch('/api/config').then(r => r.json()).then(cfg => {
            document.getElementById('apiKey').value = cfg.api_key || "";
            document.getElementById('baseURL').value = cfg.base_url || "https://api.openai.com/v1";
            document.getElementById('model').value = cfg.model || "gpt-4o";
            document.getElementById('gameURL').value = cfg.game_url || "http://127.0.0.1:8080";
        });

        function startAgent() {
            const body = {
                api_key: document.getElementById('apiKey').value,
                base_url: document.getElementById('baseURL').value,
                model: document.getElementById('model').value,
                game_url: document.getElementById('gameURL').value
            };
            fetch('/api/start', { method: 'POST', body: JSON.stringify(body) });
        }

        function stopAgent() {
            fetch('/api/stop');
        }

        // 监听日志流
        const eventSource = new EventSource('/api/logs');
        eventSource.onmessage = function(event) {
            const logs = JSON.parse(event.data);
            const logDiv = document.getElementById('logs');
            logDiv.innerHTML = logs.join('').replace(/\n/g, '<br>');
            logDiv.scrollTop = logDiv.scrollHeight;
        };
    </script>
</body>
</html>
`
