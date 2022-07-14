package service

import (
	"context"
	"encoding/json"
	"fmt"
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
	test            int
}

func NewService(cfg *types.Config) *Service {

	client := &http.Client{
		Timeout: cfg.TimeoutOutgoing * time.Second,
	}

	conn := &types.Connection{
		Connected:       0,
		LimitConnection: cfg.LimitConnection,
		Mut:             sync.Mutex{},
	}

	return &Service{
		connection:      conn,
		client:          client,
		limitGoRoutines: cfg.LimitGoRoutines,
		test:            0,
	}
}

func (s *Service) SrvMultiplexer(ctx context.Context, urls []string) (*types.ResultBody, error) {

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

	chanUrl := make(chan string, len(urls))
	chanResult := make(chan types.DataWriteChan, limitRoutines)
	result := new(types.ResultBody)

	for i := 1; i <= s.limitGoRoutines; i++ {
		go s.executorGet(ctx, chanUrl, chanResult)
	}

	for _, url := range urls {
		chanUrl <- url
	}
	close(chanUrl)

	for i := 1; i <= len(urls); i++ {

		output := <-chanResult

		if output.Error != nil {
			msg := fmt.Sprintf("one of the urls returned an error. please check the url: %s and repeat the request", output.Url)
			return nil, infrastruct.NewError(msg, http.StatusFailedDependency)
		} else {
			result.Results = append(result.Results, types.DataWrite{Url: output.Url, Result: output.Result})
		}
	}

	return result, nil
}

func (s *Service) executorGet(ctx context.Context, chanUrl <-chan string, chanResult chan<- types.DataWriteChan) {

	for urlInChanUrl := range chanUrl {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlInChanUrl, nil)
		if err != nil {
			log.Printf("err with NewRequestWithContext in executorGet: %s", err)
			chanResult <- types.DataWriteChan{Url: urlInChanUrl, Error: err}
			close(chanResult)
			return
		}

		resp, err := s.client.Do(req)
		if err != nil {
			log.Printf("err with Do in executorGet: %s", err)
			chanResult <- types.DataWriteChan{Url: urlInChanUrl, Error: err}
			close(chanResult)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			log.Printf("err with StatusCode in executorGet: StatusCode %d", resp.StatusCode)
			chanResult <- types.DataWriteChan{Url: urlInChanUrl, Error: infrastruct.NewError(resp.Status, resp.StatusCode)}
			close(chanResult)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("err with ReadAll in executorGet: %s", err)
			return
		}

		data := types.DataWriteChan{Url: urlInChanUrl}
		if err = json.Unmarshal(body, &data.Result); err != nil {
			data.Result = string(body)
		}

		chanResult <- data
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
