package lichess

import (
	"context"
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	req, _ := http.NewRequest("GET", "https://lichess.org/api/board/game/stream/"+gameId, nil)
	req.Header.Add("Authorization", "Bearer "+l.token)

	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return "GetBoard error: " + err.Error()
	}
	data, _ := io.ReadAll(resp.Body)
	return string(data)
}
