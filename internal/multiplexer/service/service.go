package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/F7icK/multiplexer/internal/multiplexer/types"
	"github.com/F7icK/multiplexer/pkg/infrastruct"
)

type Service struct {
	connection      *types.Connection
	client          *http.Client
	limitGoRoutines int
	test            int
}

func NewService(cfg *types.Config) (*Service, error) {
	switch {
	case !validationPort(cfg.Port):
		return nil, infrastruct.ErrValidationPort
	case cfg.LimitGoRoutines < 1:
		return nil, infrastruct.ErrLimitGoRoutines
	case cfg.LimitConnection < 1:
		return nil, infrastruct.ErrLimitConnection
	case cfg.TimeoutIncoming < 1:
		return nil, infrastruct.ErrTimeoutIncoming
	case cfg.TimeoutOutgoing < 1:
		return nil, infrastruct.ErrTimeoutOutgoing
	}

	client := &http.Client{
		CheckRedirect: nil,
		Transport:     http.DefaultTransport,
		Jar:           &cookiejar.Jar{},
		Timeout:       cfg.TimeoutOutgoing * time.Second,
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
	}, nil
}

/*
	Для реализации: "для каждого входящего запроса должно быть не больше 4 одновременных исходящих",
	я бы хотел использовать errgroups,
	но так как требования по ТЗ "использовать можно только компоненты стандартной библиотеки Go",
	реализовать получилось не очень красиво.
*/

func (s *Service) SrvMultiplexer(ctx context.Context, arrURLs []string) (*types.ResultBody, error) {
	if !s.newConnection() {
		return nil, infrastruct.ErrHTTPLimitConnection
	}
	defer s.closeConnection()

	var limitRoutines int
	if len(arrURLs) <= s.limitGoRoutines {
		limitRoutines = len(arrURLs)
	} else {
		limitRoutines = s.limitGoRoutines
	}

	chanURL := make(chan string, len(arrURLs))
	chanResult := make(chan types.DataWriteChan, limitRoutines)
	defer close(chanResult)
	result := new(types.ResultBody)

	for i := 1; i <= s.limitGoRoutines; i++ {
		go s.executorGet(ctx, chanURL, chanResult)
	}

	for _, url := range arrURLs {
		chanURL <- url
	}
	close(chanURL)

	for i := 1; i <= len(arrURLs); i++ {
		output := <-chanResult

		if output.Err != nil {
			msg := fmt.Sprintf("one of the urls returned an error. please check the url: %s and repeat the request", output.URL)
			return nil, infrastruct.NewError(msg, http.StatusFailedDependency)
		}

		result.Results = append(result.Results, types.DataWrite{URL: output.URL, Result: output.Result})
	}

	return result, nil
}

func (s *Service) executorGet(ctx context.Context, chanURL <-chan string, chanResult chan<- types.DataWriteChan) {
	for urlInChanURL := range chanURL {
		data := types.DataWriteChan{URL: urlInChanURL, Result: nil, Err: nil}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlInChanURL, nil)
		if err != nil {
			log.Printf("err with NewRequestWithContext in executorGet: %s", err)
			data.Err = err
			chanResult <- data
			return
		}

		resp, err := s.client.Do(req)
		if err != nil {
			log.Printf("err with Do in executorGet: %s", err)
			data.Err = err
			chanResult <- data
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			log.Printf("err with StatusCode in executorGet: StatusCode %d", resp.StatusCode)
			data.Err = infrastruct.NewError(resp.Status, resp.StatusCode)
			chanResult <- data
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("err with ReadAll in executorGet: %s", err)
			return
		}

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

func validationPort(port string) bool {
	if !strings.HasPrefix(port, ":") {
		return false
	}

	if len(port) < 2 {
		return false
	}

	for _, value := range strings.Trim(port, ":") {
		if !unicode.IsDigit(value) {
			return false
		}
	}

	return true
}
