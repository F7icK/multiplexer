package service

import (
	"context"
	"encoding/json"
	"github.com/F7icK/multiplexer/internal/multiplexer/types"
	"github.com/F7icK/multiplexer/pkg/infrastruct"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type Service struct {
	connection *types.Connection
	client     *http.Client
}

func NewService(limitConnection uint32, timeoutOutgoing time.Duration) *Service {

	client := &http.Client{
		Timeout: timeoutOutgoing * time.Second,
	}

	conn := &types.Connection{
		Connected:       0,
		LimitConnection: limitConnection,
		Mut:             sync.Mutex{},
	}

	return &Service{
		connection: conn,
		client:     client,
	}
}

func (s *Service) GetUrls(ctx context.Context, urls []string) (*types.ResultBody, error) {

	if !s.NewConnection() {
		return nil, infrastruct.ErrorLimitConnection
	}

	result := new(types.ResultBody)
	for _, urlInUrls := range urls {

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlInUrls, nil)
		if err != nil {
			log.Printf("NewRequestWithContext in GetUrls: %s", err)
			return nil, infrastruct.ErrorBadUrl
		}

		resp, err := s.client.Do(req)
		if err != nil {
			log.Printf("Do in GetUrls: %s", err)
			return nil, infrastruct.ErrorBadUrl
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			log.Printf("StatusCode in GetUrls")
			return nil, infrastruct.ErrorBadUrl
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return nil, infrastruct.ErrorInternalServerError
		}

		data := types.DataWrite{Url: urlInUrls}
		if err = json.Unmarshal(body, &data.Result); err != nil {
			data.Result = string(body)
		}

		result.Results = append(result.Results, data)
	}

	s.CloseConnection()
	return result, nil
}

func (s *Service) NewConnection() bool {
	s.connection.Mut.Lock()
	defer s.connection.Mut.Unlock()

	if s.connection.Connected >= s.connection.LimitConnection {
		return false
	}

	s.connection.Connected++

	return true
}

func (s *Service) CloseConnection() {
	s.connection.Mut.Lock()
	defer s.connection.Mut.Unlock()

	if s.connection.Connected > 0 {
		s.connection.Connected--
	}
}
