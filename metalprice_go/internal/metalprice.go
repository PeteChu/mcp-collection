package metalprice

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type MetalPrice struct {
	Server *server.MCPServer
	ApiKey string
}

type BaseApiResponse struct {
	Success bool `json:"success"`
	Error   struct {
		StatusCode int    `json:"statusCode"`
		Message    string `json:"message"`
	} `json:"error"`
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
	m.Server.AddTool(m.LiveRates())
}

func (m *MetalPrice) ListSymbols() (mcp.Tool, ToolHandler) {
	tool := mcp.NewTool("metalprice_list_symbols",
		mcp.WithDescription("Get list of all supported currencies"),
	)
	return tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client := http.Client{}

		req, err := http.NewRequest("GET", BASE_URL+"/symbols", nil)
		if err != nil {
			return nil, fmt.Errorf("error creating request: %v", err)
		}

		req.Header = http.Header{
			"X-API-KEY":    {m.ApiKey},
			"Content-Type": {"application/json"},
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error making request: %v", err)
		}
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)

		if err := handleApiError(b); err != nil {
			return nil, err
		}

		if err != nil {
			return nil, fmt.Errorf("error reading response: %v", err)
		}

		return mcp.NewToolResultText(string(b)), nil
	}
}

func (m *MetalPrice) LiveRates() (mcp.Tool, ToolHandler) {
	tool := mcp.NewTool("metalprice_live_rates",
		mcp.WithDescription("Get real-time exchange rate data for all available/specific currencies"),
		mcp.WithString(
			"base",
			mcp.Description("Specify a base currency. Base Currency will default to USD if this parameter is not defined."),
			mcp.DefaultString("usd"),
		),
		mcp.WithString(
			"currencies",
			mcp.Description("Specify a comma-separated list of currency codes to limit API responses to specified currencies. If this parameter is not defined, the API will return all supported currencies."),
		),
		mcp.WithString(
			"unit",
			mcp.Description("(Paid plan) Specify troy_oz or gram or kilogram. If not defined, the API will return metals in troy ounce."),
			mcp.Enum("troy_oz", "gram", "kilogram"),
			mcp.DefaultString("troy_oz"),
		),
		mcp.WithString("math",
			mcp.Description("Specify math operators to perform on the result. Use value to refer to the rates. Specify one or more of the mathematical operators add, subtract, multiply, and/or divide."),
		),
	)
	return tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		params := make([]string, 0)
		params = append(params, fmt.Sprintf("api_key=%s", m.ApiKey))

		if base, ok := request.Params.Arguments["base"]; ok {
			params = append(params, fmt.Sprintf("base=%s", base.(string)))
		}
		if currencies, ok := request.Params.Arguments["currencies"]; ok {
			params = append(params, fmt.Sprintf("currencies=%s", currencies.(string)))
		}
		if unit, ok := request.Params.Arguments["unit"]; ok {
			params = append(params, fmt.Sprintf("unit=%s", unit.(string)))
		}
		if math, ok := request.Params.Arguments["math"]; ok {
			params = append(params, fmt.Sprintf("math=%s", math.(string)))
		}

		query := ""
		for i, val := range params {
			if i == 0 {
				query += "?" + val
			} else {
				query += "&" + val
			}
		}

		url := BASE_URL + "/latest" + query

		client := http.Client{}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating request: %v", err)
		}

		req.Header = http.Header{
			"Content-Type": {"application/json"},
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error making request: %v", err)
		}
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)

		if err := handleApiError(b); err != nil {
			return nil, err
		}

		if err != nil {
			return nil, fmt.Errorf("error reading response: %v", err)
		}

		return mcp.NewToolResultText(string(b)), nil
	}
}
