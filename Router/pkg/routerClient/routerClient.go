package routerClient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	chatv1 "github.com/QutaqKicker/ChatParser/Protos/gen/go/chat"
	commonv1 "github.com/QutaqKicker/ChatParser/Protos/gen/go/common"
	userv1 "github.com/QutaqKicker/ChatParser/Protos/gen/go/user"
	"io"
	"net/http"
	"strconv"
	"time"
)

type RouterClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewRouterClient(baseURL string) *RouterClient {
	return &RouterClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second, // Установите тайм-аут на 10 секунд (можно изменить)
		},
	}
}

func (c *RouterClient) BuildUrl(endpoint string) string {
	return fmt.Sprintf("%s%s", c.baseURL, endpoint)
}

func (c *RouterClient) GetMessagesCount(ctx context.Context, filter *commonv1.MessagesFilter) (int64, error) {
	request, err := c.getRequest(ctx, "/chat/messages/count", http.MethodPost, filter)
	if err != nil {
		return 0, err
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return 0, err
	}

	if response.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("received from server incorrect status code: %d", response.StatusCode)
	}

	if buf, err := io.ReadAll(response.Body); err == nil {
		messagesCount, err := strconv.ParseInt(string(buf), 10, 64)
		return messagesCount, err
	} else {
		return 0, err
	}
}

func (c *RouterClient) SearchMessages(ctx context.Context, messagesRequest *chatv1.SearchMessagesRequest) (*chatv1.SearchMessagesResponse, error) {
	request, err := c.getRequest(ctx, "/chat/messages/search", http.MethodPost, messagesRequest)
	if err != nil {
		return nil, err
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received from server incorrect status code: %d", response.StatusCode)
	}

	var messagesResponse *chatv1.SearchMessagesResponse
	err = json.NewDecoder(response.Body).Decode(messagesResponse)
	return messagesResponse, err
}

func (c *RouterClient) ParseFromDir(ctx context.Context, dirPath string) (bool, error) {
	request, err := c.getRequest(ctx, "/chat/parse-from-dir", http.MethodGet, nil)
	if err != nil {
		return false, err
	}

	request.URL.Query().Set("dir-path", dirPath)
	response, err := c.httpClient.Do(request)
	if err != nil {
		return false, err
	}

	if response.StatusCode != http.StatusOK {
		return false, fmt.Errorf("received from server incorrect status code: %d", response.StatusCode)
	} else {
		return true, nil
	}
}

func (c *RouterClient) ExportToDir(ctx context.Context, filter *commonv1.MessagesFilter, exportType string) (bool, error) {
	request, err := c.getRequest(ctx, "/backup/export-to-dir", http.MethodPost, nil)
	if err != nil {
		return false, err
	}

	request.URL.Query().Set("export-type", exportType)
	response, err := c.httpClient.Do(request)
	if err != nil {
		return false, err
	}

	if response.StatusCode != http.StatusOK {
		return false, fmt.Errorf("received from server incorrect status code: %d", response.StatusCode)
	} else {
		return true, nil
	}
}

func (c *RouterClient) GetUsersWithMessagesCount(ctx context.Context) (*userv1.GetUsersResponse, error) {
	request, err := c.getRequest(ctx, "/user/messages-count", http.MethodPost, nil)
	if err != nil {
		return nil, err
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received from server incorrect status code: %d", response.StatusCode)
	}

	var userResponse *userv1.GetUsersResponse
	err = json.NewDecoder(response.Body).Decode(userResponse)
	return userResponse, err

}

func (c *RouterClient) getRequest(ctx context.Context, url, method string, payload any) (*http.Request, error) {
	payloadBuffer := new(bytes.Buffer)
	if payload != nil {
		err := json.NewDecoder(payloadBuffer).Decode(payload)
		if err != nil {
			return nil, err
		}
	}

	request, err := http.NewRequestWithContext(ctx, method, c.BuildUrl(url), payloadBuffer)
	request.Header.Set("Content-Type", "application/json")
	return request, err
}
