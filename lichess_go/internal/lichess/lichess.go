package lichess

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Lichess struct {
	token string
}

func NewLichess(token string) *Lichess {
	return &Lichess{
		token: token,
	}
}

func (l Lichess) GetBoard(gameId string) string {
	client := http.Client{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	req, err := http.NewRequest("GET", "https://lichess.org/api/board/game/stream/"+gameId, nil)
	if err != nil {
		return "GetBoard error creating request: " + err.Error()
	}
	req.Header.Add("Authorization", "Bearer "+l.token)

	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return "GetBoard error: " + err.Error()
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	return string(data)
}

func (l Lichess) MakeMove(gameId string, move string, draw bool) string {
	client := http.Client{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	url := fmt.Sprintf("https://lichess.org/api/board/game/%s/move/%s", gameId, move)
	if draw {
		url += "?offeringDraw=true"
	}

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "MakeMove error creating request: " + err.Error()
	}
	req.Header.Add("Authorization", "Bearer "+l.token)

	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return "MakeMove error: " + err.Error()
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "MakeMove error reading response: " + err.Error()
	}
	return string(data)
}
