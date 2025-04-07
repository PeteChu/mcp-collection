package metalprice

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func handleApiError(body []byte) error {
	var responseBody BaseApiResponse
	json.Unmarshal(body, &responseBody)
	if !responseBody.Success {
		code := responseBody.Error.StatusCode
		switch code {
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
			return fmt.Errorf("API error: The user's query did not return any results [ latest, historical, convert, timeframe, change ]")
		default:
			return fmt.Errorf("API error: %+v", code)
		}
	}

	return nil
}
