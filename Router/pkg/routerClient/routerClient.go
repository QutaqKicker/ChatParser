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
)

func GetMessagesCount(ctx context.Context, client *http.Client, filter *commonv1.MessagesFilter) (int64, error) {
	request, err := getRequest(ctx, "/chat/messages/count", "POST", filter)
	if err != nil {
		return 0, err
	}

	response, err := client.Do(request)
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

func SearchMessages(ctx context.Context, client *http.Client, messagesRequest *chatv1.SearchMessagesRequest) (*chatv1.SearchMessagesResponse, error) {
	request, err := getRequest(ctx, "/chat/messages/search", "POST", messagesRequest)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)
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

func ParseFromDir(ctx context.Context, client *http.Client, dirPath string) (bool, error) {
	request, err := getRequest(ctx, "/chat/parse-from-dir", "GET", nil)
	if err != nil {
		return false, err
	}

	request.URL.Query().Set("dir-path", dirPath)
	response, err := client.Do(request)
	if err != nil {
		return false, err
	}

	if response.StatusCode != http.StatusOK {
		return false, fmt.Errorf("received from server incorrect status code: %d", response.StatusCode)
	} else {
		return true, nil
	}
}

func ExportToDir(ctx context.Context, client *http.Client, filter *commonv1.MessagesFilter, exportType string) (bool, error) {
	request, err := getRequest(ctx, "/backup/export-to-dir", "POST", nil)
	if err != nil {
		return false, err
	}

	request.URL.Query().Set("export-type", exportType)
	response, err := client.Do(request)
	if err != nil {
		return false, err
	}

	if response.StatusCode != http.StatusOK {
		return false, fmt.Errorf("received from server incorrect status code: %d", response.StatusCode)
	} else {
		return true, nil
	}
}

func GetUsersWithMessagesCount(ctx context.Context, client *http.Client) (*userv1.GetUsersResponse, error) {
	request, err := getRequest(ctx, "/user/messages-count", "POST", nil)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)
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

func getRequest(ctx context.Context, url, method string, payload any) (*http.Request, error) {
	payloadBuffer := new(bytes.Buffer)
	err := json.NewDecoder(payloadBuffer).Decode(payload)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, method, url, payloadBuffer)
	return request, err
}
