package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/cjlapao/postfixcli-backend-api/ioc"
	"github.com/cjlapao/postfixcli-backend-api/models"
)

var globalLinodeService *LinodeService

type LinodeService struct {
	Context context.Context
	Token   string
	Client  *http.Client
	BaseUrl string
	Version string
}

func GetLinodeService() *LinodeService {
	if globalLinodeService != nil {
		return globalLinodeService
	}

	return NewLinodeService()
}

func NewLinodeService() *LinodeService {
	globalLinodeService = &LinodeService{
		BaseUrl: "https://api.linode.com",
		Version: "v4",
	}

	globalLinodeService.Context = context.Background()

	linodeToken := ioc.Config.GetString("linode__token")
	if linodeToken != "" {
		globalLinodeService.Token = linodeToken
	}

	globalLinodeService.Client = globalLinodeService.getHttpClient()
	return globalLinodeService
}

func (svc *LinodeService) generateFilterHeader(filters []models.LinodeFilter) (string, string) {
	var filterResult string

	if len(filters) > 0 {
		filterResult = "{ "
		for idx, filter := range filters {
			if idx > 0 {
				filterResult = fmt.Sprintf("%v,", filterResult)
			}
			filterResult = fmt.Sprintf("%v \"%v\":\"%v\"", filterResult, filter.Field, filter.Value)
		}
		filterResult = fmt.Sprintf("%v }", filterResult)
	} else {
		return "", ""
	}

	return "X-Filter", filterResult
}

func (svc *LinodeService) getBaseUrl() string {
	return fmt.Sprintf("%v/%v", svc.BaseUrl, svc.Version)
}

func (svc *LinodeService) getHttpClient() *http.Client {
	client := http.Client{}

	return &client
}

func (svc *LinodeService) get(path string, filter ...models.LinodeFilter) (*http.Request, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%v/%v", svc.getBaseUrl(), path), nil)

	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %v", svc.Token))

	if len(filter) > 0 {
		filterHeader, filterHeaderValue := svc.generateFilterHeader(filter)
		request.Header.Add(filterHeader, filterHeaderValue)
	}

	return request, nil
}

func (svc *LinodeService) readBody(resp *http.Response, dest any) error {
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

func (svc *LinodeService) processResponse(resp *http.Response, dest any) error {
	if resp.StatusCode <= 199 || resp.StatusCode >= 300 {
		var linodeError models.LinodeErrorResponse
		err := svc.readBody(resp, &linodeError)

		if err != nil {
			return err
		}

		return fmt.Errorf("%v: %v", linodeError.Errors[0].Field, linodeError.Errors[0].Reason)
	}

	err := svc.readBody(resp, &dest)

	if err != nil {
		return err
	}

	return nil
}

func (svc *LinodeService) GetNodeBalancerDetails(ipv4 string) (*models.NodeBalancerDetails, error) {
	filter := models.LinodeFilter{
		Field: "ipv4",
		Value: ipv4,
	}

	request, err := svc.get("/nodebalancers", filter)

	if err != nil {
		return nil, err
	}

	resp, err := svc.Client.Do(request)

	if err != nil {
		return nil, err
	}

	var response models.LinodeLoadBalancerResponse

	err = svc.processResponse(resp, &response)

	if err != nil {
		return nil, err
	}

	result := models.NodeBalancerDetails{
		Hostname: response.Data[0].Hostname,
		IPV4:     response.Data[0].Ipv4,
		IPV6:     response.Data[0].Ipv6,
	}

	ioc.Log.Info(fmt.Sprintf("%v", result))
	return &result, nil
}

func (svc *LinodeService) GetNodeInstances() (*[]models.LinodeInstanceDetail, error) {
	result := make([]models.LinodeInstanceDetail, 0)
	request, err := svc.get("/linode/instances")

	if err != nil {
		return nil, err
	}

	resp, err := svc.Client.Do(request)

	if err != nil {
		return nil, err
	}

	var response models.LinodeInstanceResponse

	err = svc.processResponse(resp, &response)

	if err != nil {
		return nil, err
	}

	if len(response.Data) > 0 {
		for _, responseInstance := range response.Data {
			instance := models.LinodeInstanceDetail{
				ID:   responseInstance.ID,
				Name: responseInstance.Label,
			}
			ipv6 := strings.ReplaceAll(responseInstance.Ipv6, "/128", "")
			instance.Ips = append(instance.Ips, models.Ip{Ip: ipv6, Type: models.IPV6})
			for _, ip := range responseInstance.Ipv4 {
				instance.Ips = append(instance.Ips, models.Ip{Ip: ip, Type: models.IPV4})
			}

			result = append(result, instance)
		}
	}

	ioc.Log.Info(fmt.Sprintf("%v", result))
	return &result, nil
}
