// tool.go
package serpapi

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// 定义搜索引擎类型
const (
	EngineGoogle = "google"
	EngineBing   = "bing"
)

// SerpAPITool 实现 Eino Tool 接口
type SerpAPITool struct {
	apiKey string
}

func NewSerpAPITool(apiKey string) tool.InvokableTool {
	return &SerpAPITool{apiKey: apiKey}
}

func (s *SerpAPITool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "search_internet",
		Desc: "Perform internet search using SerpAPI with Google or Bing engines.",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"query": {
				Desc:     "The search query text",
				Type:     "string",
				Required: true,
			},
			"engine": {
				Desc: "The search engine to use (google or bing)",
				Type: "string",
			},
		}),
	}, nil
}

func (s *SerpAPITool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	var args map[string]interface{}
	if err := json.Unmarshal([]byte(argumentsInJSON), &args); err != nil {
		return "", fmt.Errorf("failed to parse arguments: %v", err)
	}

	query, ok := args["query"].(string)
	if !ok || query == "" {
		return "", fmt.Errorf("invalid or missing 'query' parameter")
	}

	engine := EngineGoogle
	if e, exists := args["engine"].(string); exists {
		if e != EngineGoogle && e != EngineBing {
			return "", fmt.Errorf("invalid engine: %s", e)
		}
		engine = e
	}

	body, err := PerformSearch(SearchRequest{
		Query:  query,
		Engine: engine,
		APIKey: s.apiKey,
	})
	if err != nil {
		return "", err
	}

	return body, nil
}
