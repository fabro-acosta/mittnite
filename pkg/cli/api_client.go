package cli

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/unitedsymme/mittnite/pkg/proc"
	"net/http"
)

const (
	ApiActionJobStart   = "start"
	ApiActionJobRestart = "restart"
	ApiActionJobStop    = "stop"
	ApiActionJobStatus  = "status"
	ApiActionJobLogs    = "logs"
)

type APIClient struct {
	apiAddress string
}

func NewApiClient(apiAddress string) *APIClient {
	return &APIClient{
		apiAddress: apiAddress,
	}
}

func (api *APIClient) CallAction(job, action string) APIResponse {
	switch action {
	case ApiActionJobStart:
		return api.JobStart(job)
	case ApiActionJobRestart:
		return api.JobRestart(job)
	case ApiActionJobStop:
		return api.JobStop(job)
	default:
		return &CommonAPIResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "",
			Error:      fmt.Errorf("unknown action %s", action),
		}
	}
}

func (api *APIClient) JobStart(job string) APIResponse {
	client, url, err := api.buildHTTPClientAndURL()
	if err != nil {
		return &CommonAPIResponse{Error: err}
	}

	url.Path = fmt.Sprintf("/v1/job/%s/start", job)
	return NewAPIResponse(client.Post(url.String(), "application/json", nil))
}

func (api *APIClient) JobRestart(job string) APIResponse {
	client, url, err := api.buildHTTPClientAndURL()
	if err != nil {
		return &CommonAPIResponse{Error: err}
	}
	url.Path = fmt.Sprintf("/v1/job/%s/restart", job)
	return NewAPIResponse(client.Post(url.String(), "application/json", nil))
}

func (api *APIClient) JobStop(job string) APIResponse {
	client, url, err := api.buildHTTPClientAndURL()
	if err != nil {
		return &CommonAPIResponse{Error: err}
	}

	url.Path = fmt.Sprintf("/v1/job/%s/stop", job)
	return NewAPIResponse(client.Post(url.String(), "application/json", nil))
}

func (api *APIClient) JobStatus(job string) TypedAPIResponse[proc.CommonJobStatus] {
	client, url, err := api.buildHTTPClientAndURL()
	if err != nil {
		return TypedAPIResponse[proc.CommonJobStatus]{Error: err}
	}

	url.Path = fmt.Sprintf("/v1/job/%s/status", job)
	return *NewTypedAPIResponse(proc.CommonJobStatus{})(client.Get(url.String()))
}

func (api *APIClient) JobList() TypedAPIResponse[[]string] {
	client, url, err := api.buildHTTPClientAndURL()
	if err != nil {
		return TypedAPIResponse[[]string]{Error: err}
	}

	url.Path = "/v1/jobs"
	return *NewTypedAPIResponse(make([]string, 0))(client.Get(url.String()))
}

func (api *APIClient) JobLogs(job string, follow bool, tailLen int) APIResponse {
	dialer, url, err := api.buildWebsocketURL()
	if err != nil {
		return &CommonAPIResponse{Error: fmt.Errorf("error building websocket url: %w", err)}
	}

	qryValues := url.Query()
	qryValues.Add("taillen", fmt.Sprintf("%d", tailLen))
	if follow {
		qryValues.Add("follow", "true")
	}

	url.RawQuery = qryValues.Encode()
	url.Path = fmt.Sprintf("/v1/job/%s/logs", job)

	handler := func(ctx context.Context, conn *websocket.Conn, msgChan chan []byte, errChan chan error) {
		for {
			select {
			default:
				_, msg, err := conn.ReadMessage()
				if err != nil {
					errChan <- err
					return
				}
				msgChan <- msg
			case <-ctx.Done():
				return
			}
		}
	}
	return NewStreamingAPIResponse(url, dialer, handler)
}
