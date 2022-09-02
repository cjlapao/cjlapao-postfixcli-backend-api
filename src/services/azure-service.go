package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cjlapao/postfixcli-backend-api/ioc"
	"github.com/cjlapao/postfixcli-backend-api/models"
	"github.com/pascaldekloe/jwt"
)

var globalAzureService *AzureService

type AzureService struct {
	Context        context.Context
	TenantId       string
	SubscriptionId string
	ClientId       string
	ClientSecret   string
	Resource       string
	Client         *http.Client
	baseUrl        string
	baseLoginUrl   string
	token          string
}

func GetAzureService() *AzureService {
	if globalAzureService != nil {
		return globalAzureService
	}

	return NewAzureService()
}

func NewAzureService() *AzureService {
	globalAzureService = &AzureService{
		baseUrl:      "https://management.azure.com/subscriptions/${AZURE_SUBSCRIPTION_ID}/resourceGroups/${AZURE_RESOURCE_GROUP=}",
		baseLoginUrl: "https://login.microsoftonline.com/${AZURE_TENANT_ID}/oauth2/token",
		Resource:     "https://management.azure.com",
	}

	globalAzureService.Context = context.Background()

	globalAzureService.TenantId = ioc.Config.GetString("azure__tenant__id")
	globalAzureService.SubscriptionId = ioc.Config.GetString("azure__subscription__id")
	globalAzureService.ClientId = ioc.Config.GetString("azure__client__id")
	globalAzureService.ClientSecret = ioc.Config.GetString("azure__client__secret")

	globalAzureService.Client = globalAzureService.getHttpClient()
	return globalAzureService
}

func (svc *AzureService) getHttpClient() *http.Client {
	client := http.Client{}

	return &client
}

func (svc *AzureService) readBody(resp *http.Response, dest any) error {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &dest)

	if err != nil {
		return err
	}

	return nil
}

func (svc *AzureService) processResponse(resp *http.Response, dest any) error {
	if resp.StatusCode <= 199 || resp.StatusCode >= 300 {
		var azureError models.AzureErrorResponse
		err := svc.readBody(resp, &azureError)

		if azureError.Error.Code == "" {
			body, _ := io.ReadAll(resp.Body)
			azureError.Error.Code = fmt.Sprintf("EMPTY_%v", resp.StatusCode)
			if string(body) == "" {
				azureError.Error.Message = "Empty error"
			} else {
				azureError.Error.Message = string(body)
			}
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("%v: %v", azureError.Error.Code, azureError.Error.Message)
	}

	err := svc.readBody(resp, &dest)

	if err != nil {
		return err
	}

	return nil
}

func (svc *AzureService) getBaseUrl(resourceGroup string) string {
	baseUrl := strings.ReplaceAll(svc.baseUrl, "${AZURE_SUBSCRIPTION_ID}", svc.SubscriptionId)
	baseUrl = strings.ReplaceAll(baseUrl, "${AZURE_RESOURCE_GROUP=}", resourceGroup)
	return baseUrl
}

func (svc *AzureService) getBaseLoginUrl() string {
	return strings.ReplaceAll(svc.baseLoginUrl, "${AZURE_TENANT_ID}", svc.TenantId)
}

func (svc *AzureService) getToken() string {
	baseUrl := svc.getBaseLoginUrl()

	form := url.Values{}
	form.Add("grant_type", "client_credentials")
	form.Add("client_id", svc.ClientId)
	form.Add("client_secret", svc.ClientSecret)
	form.Add("resource", svc.Resource)

	request, err := http.NewRequest("POST", baseUrl, strings.NewReader(form.Encode()))
	if err != nil {
		return ""
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := svc.Client.Do(request)

	if err != nil {
		return ""
	}

	var tokenResponse models.AzureTokenResponse
	err = svc.readBody(response, &tokenResponse)

	if err != nil {
		return ""
	}

	svc.token = tokenResponse.AccessToken

	return svc.token
}

func (svc *AzureService) get(endpointUrl string) (*http.Request, error) {
	request, err := http.NewRequest("GET", endpointUrl, nil)

	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %v", svc.Token()))
	request.Header.Add("Content-type", "application/json")

	return request, nil
}

func (svc *AzureService) put(endpointUrl string, body any) (*http.Request, error) {
	jsonBody, err := json.MarshalIndent(body, "", "  ")

	ioc.Log.Info(string(jsonBody))
	if err != nil {
		return nil, err
	}

	if string(jsonBody) == "" {
		return nil, errors.New("body cannot be empty")
	}

	request, err := http.NewRequest("PUT", endpointUrl, strings.NewReader(string(jsonBody)))

	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %v", svc.Token()))
	request.Header.Add("Content-type", "application/json")

	return request, nil
}

func (svc *AzureService) Token() string {
	if svc.token != "" {
		expired := false
		claims, err := jwt.ParseWithoutCheck([]byte(svc.token))
		if claims.Expires.Time().Before(time.Now()) || err != nil {
			expired = true
		}

		if !expired {
			return svc.token
		}
	}

	return svc.getToken()
}
