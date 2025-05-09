package serpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

// TestSerpAPITool 测试 SerpAPITool 是否能正常工作
func TestSerpAPITool(t *testing.T) {
	apiKey := os.Getenv("SERPAPI_API_KEY")
	fmt.Println("SERPAPI_API_KEY:", apiKey)
	if apiKey == "" {
		t.Skip("SERPAPI_API_KEY not set, skipping test.")
	}

	tool := NewSerpAPITool(apiKey)

	tests := []struct {
		name    string
		query   string
		engine  string
		wantErr bool
	}{
		{
			name:    "Google Search for '广州塔'",
			query:   "广州塔",
			engine:  EngineGoogle,
			wantErr: false,
		},
		{
			name:    "Bing Search for '成都武侯祠'",
			query:   "成都武侯祠",
			engine:  EngineBing,
			wantErr: false,
		},
		{
			name:    "Invalid engine type",
			query:   "example",
			engine:  "invalid_engine",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := map[string]interface{}{
				"query":  tt.query,
				"engine": tt.engine,
			}
			argsJSON, _ := json.Marshal(args)

			result, err := tool.InvokableRun(context.Background(), string(argsJSON))
			if (err != nil) != tt.wantErr {
				t.Errorf("InvokableRun() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && result == "" {
				t.Error("Expected non-empty result, got empty")
			}

			fmt.Printf("Query: %s\nEngine: %s\nResult:\n%s\n\n", tt.query, tt.engine, result)
		})
	}
}
