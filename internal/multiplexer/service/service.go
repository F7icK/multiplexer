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
	connection      *types.Connection
	client          *http.Client
	limitGoRoutines int
}

func NewService(limitConnection uint32, timeoutOutgoing time.Duration, limitGoRoutines int) *Service {

	client := &http.Client{
		Timeout: timeoutOutgoing * time.Second,
	}

	conn := &types.Connection{
		Connected:       0,
		LimitConnection: limitConnection,
		Mut:             sync.Mutex{},
	}

	return &Service{
		connection:      conn,
		client:          client,
		limitGoRoutines: limitGoRoutines,
	}
}

func (s *Service) GetUrls(ctx context.Context, urls []string) (*types.ResultBody, error) {

	if !s.newConnection() {
		return nil, infrastruct.ErrorLimitConnection
	}
	defer s.closeConnection()

	var limitRoutines int
	if len(urls) <= s.limitGoRoutines {
		limitRoutines = len(urls)
	} else {
		limitRoutines = s.limitGoRoutines
	}

	chanUrl := make(chan string, limitRoutines)
	chanResult := make(chan types.DataWrite, limitRoutines)
	result := new(types.ResultBody)

	for i := 1; i <= s.limitGoRoutines && i <= len(urls); i++ {
		go s.executorGet(ctx, chanUrl, chanResult)
	}

	for _, url := range urls {
		chanUrl <- url
	}
	close(chanUrl)

	for i := 1; i <= len(urls); i++ {

		output := <-chanResult
		if output.Url != "" {
			result.Results = append(result.Results, output)
		} else {
			return nil, infrastruct.ErrorBadUrl
		}
	}

	return result, nil
}

func (s *Service) executorGet(ctx context.Context, jobs <-chan string, results chan<- types.DataWrite) {

	for urlInUrls := range jobs {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlInUrls, nil)
		if err != nil {
			log.Printf("NewRequestWithContext in GetUrls: %s", err)
			close(results)
			return
		}

		resp, err := s.client.Do(req)
		if err != nil {
			log.Printf("err with Do in executorGet: %s", err)
			close(results)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			log.Printf("err with StatusCode in executorGet: StatusCode %d", resp.StatusCode)
			close(results)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("err with ReadAll in executorGet: %s", err)
			return
		}

		data := types.DataWrite{Url: urlInUrls}
		if err = json.Unmarshal(body, &data.Result); err != nil {
			data.Result = string(body)
		}

		results <- data
	}
}

func (s *Service) newConnection() bool {
	s.connection.Mut.Lock()
	defer s.connection.Mut.Unlock()

	if s.connection.Connected >= s.connection.LimitConnection {
		return false
	}

	s.connection.Connected++

	return true
}

func (s *Service) closeConnection() {
	s.connection.Mut.Lock()
	defer s.connection.Mut.Unlock()

	if s.connection.Connected > 0 {
		s.connection.Connected--
	}
}
