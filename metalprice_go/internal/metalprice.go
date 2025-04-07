package metalprice

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type MetalPrice struct {
	Server *server.MCPServer
	ApiKey string
}

type ToolHandler = func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)

const BASE_URL = "https://api.metalpriceapi.com/v1"

func NewMetalPriceMcp(server *server.MCPServer, apiKey string) *MetalPrice {
	mcp := &MetalPrice{
		Server: server,
		ApiKey: apiKey,
	}

	mcp.registerTools()
	return mcp
}

func (m *MetalPrice) registerTools() {
	m.Server.AddTool(m.ListSymbols())
}

func (m *MetalPrice) ListSymbols() (mcp.Tool, ToolHandler) {
	tool := mcp.NewTool("list symbols",
		mcp.WithDescription("List all available symbols"),
	)
	return tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client := http.Client{}

		req, err := http.NewRequest("GET", BASE_URL+"/symbols", nil)
		if err != nil {
		}

		req.Header = http.Header{
			"X-API-KEY":    {m.ApiKey},
			"Content-Type": {"application/json"},
		}

		resp, err := client.Do(req)
		if err != nil {
		}

		if err := handleApiError(resp); err != nil {
			return nil, err
		}

		var data map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, fmt.Errorf("error decoding response: %v", err)
		}

		prettyJson, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("error formatting response: %v", err)
		}

		return mcp.NewToolResultText(string(prettyJson)), nil
	}
}

func handleApiError(response *http.Response) error {
	if response.StatusCode != http.StatusOK {
		switch response.StatusCode {
		case http.StatusNotFound:
			return fmt.Errorf("API error: User requested a non-existent API function")
		case 101:
			return fmt.Errorf("API error: User did not supply an API Key")
		case 102:
			return fmt.Errorf("API error: User did not supply an access key or supplied an invalid access key")
		case 103:
			return fmt.Errorf("API error: The user's account is not active. User will be prompted to get in touch with Customer Support")
		case 104:
			return fmt.Errorf("API error: Too Many Requests")
		case 105:
			return fmt.Errorf("API error: User has reached or exceeded his subscription plan's monthly API request allowance")
		case 201:
			return fmt.Errorf("API error: User entered an invalid Base Currency [ latest, historical, timeframe, change ]")
		case 202:
			return fmt.Errorf("API error: User entered an invalid from Currency [ convert ]")
		case 203:
			return fmt.Errorf("API error: User entered invalid to currency [ convert ]")
		case 204:
			return fmt.Errorf("API error: User entered invalid amount [ convert ]")
		case 205:
			return fmt.Errorf("API error: User entered invalid date [ historical, convert, timeframe, change ]")
		case 206:
			return fmt.Errorf("API error: Invalid timeframe [ timeframe, change ]")
		case 207:
			return fmt.Errorf("API error: Timeframe exceeded 365 days [ timeframe ]")
		case 300:
			return fmt.Errorf("API error:  	The user's query did not return any results [ latest, historical, convert, timeframe, change ]")
		default:
			return fmt.Errorf("API error: %s", response.Status)
		}
	}

	return nil
}
