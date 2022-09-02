package services

import (
	"errors"
	"fmt"

	"github.com/cjlapao/postfixcli-backend-api/ioc"
	"github.com/cjlapao/postfixcli-backend-api/models"
)

func (svc *AzureService) GetDnsZones(resourceGroup string) (*models.AzureDNSZones, error) {
	endpoint := "/providers/Microsoft.Network/dnsZones?api-version=2018-05-01"
	baseUrl := fmt.Sprintf("%v/%v", svc.getBaseUrl(resourceGroup), endpoint)

	request, err := svc.get(baseUrl)

	if err != nil {
		return nil, err
	}

	response, err := svc.Client.Do(request)

	if err != nil {
		return nil, err
	}

	var result models.AzureDNSZones
	err = svc.processResponse(response, &result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (svc *AzureService) GetDnsZone(resourceGroup string, zoneName string) (*models.AzureDNSZone, error) {
	endpoint := fmt.Sprintf("/providers/Microsoft.Network/dnsZones/%v?api-version=2018-05-01", zoneName)
	baseUrl := fmt.Sprintf("%v/%v", svc.getBaseUrl(resourceGroup), endpoint)

	request, err := svc.get(baseUrl)

	if err != nil {
		return nil, err
	}

	response, err := svc.Client.Do(request)

	if err != nil {
		return nil, err
	}

	var result models.AzureDNSZone
	err = svc.processResponse(response, &result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (svc *AzureService) GetDnsRecord(resourceGroup string, zoneName string, recordType string, recordName string) (*models.AzureDNSRecord, error) {
	endpoint := fmt.Sprintf("/providers/Microsoft.Network/dnsZones/%v/%v/%v?api-version=2018-05-01", zoneName, recordType, recordName)
	baseUrl := fmt.Sprintf("%v/%v", svc.getBaseUrl(resourceGroup), endpoint)

	request, err := svc.get(baseUrl)

	if err != nil {
		return nil, err
	}

	response, err := svc.Client.Do(request)

	if err != nil {
		return nil, err
	}

	var result models.AzureDNSRecord
	err = svc.processResponse(response, &result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (svc *AzureService) UpsertDnsRecord(resourceGroup string, zoneName string, dnsRecord models.AzureDNSRecord) (*models.AzureDNSRecord, error) {
	endpoint := fmt.Sprintf("/providers/Microsoft.Network/dnsZones/%v/%v/%v?api-version=2018-05-01", zoneName, dnsRecord.Type, dnsRecord.Name)
	baseUrl := fmt.Sprintf("%v/%v", svc.getBaseUrl(resourceGroup), endpoint)

	putBody := models.AzureDNSRecord{
		Properties: models.AzureDnsRecordProperties{
			Metadata: &models.AzureDnsMetadata{
				Key1: dnsRecord.Name,
			},
			TTL:         dnsRecord.Properties.TTL,
			FQDN:        dnsRecord.Properties.FQDN,
			ARecords:    dnsRecord.Properties.ARecords,
			AAAARecords: dnsRecord.Properties.AAAARecords,
			MXRecords:   dnsRecord.Properties.MXRecords,
			SRVRecords:  dnsRecord.Properties.SRVRecords,
			CNAMERecord: dnsRecord.Properties.CNAMERecord,
		},
	}

	request, err := svc.put(baseUrl, putBody)

	if err != nil {
		return nil, err
	}

	response, err := svc.Client.Do(request)

	if err != nil {
		return nil, err
	}

	var result models.AzureDNSRecord
	err = svc.processResponse(response, &result)

	if err != nil {
		return nil, err
	}

	operation := "created"
	if response.StatusCode == 200 {
		operation = "updated"
	}

	ioc.Log.Info("Dns %v Record %v was %v correctly for zone %v", dnsRecord.Type, dnsRecord.Name, operation, zoneName)

	return &result, nil
}

func (svc *AzureService) UpsertADnsRecord(resourceGroup string, zoneName string, name string, ttl int, ipv4s ...string) error {
	if len(ipv4s) == 0 {
		return errors.New("you need at least one ipv4")
	}

	putModel := models.AzureDNSRecord{
		Name: name,
		Type: "A",
		Properties: models.AzureDnsRecordProperties{
			TTL: int64(ttl),
		},
	}

	for _, ipv4 := range ipv4s {
		r := models.AzureDnsARecord{
			Ipv4Address: ipv4,
		}
		putModel.Properties.ARecords = append(putModel.Properties.ARecords, r)
	}

	_, err := svc.UpsertDnsRecord(resourceGroup, zoneName, putModel)

	return err
}

func (svc *AzureService) UpsertAAAADnsRecord(resourceGroup string, zoneName string, name string, ttl int, ipv6s ...string) error {
	if len(ipv6s) == 0 {
		return errors.New("you need at least one ipv6")
	}

	putModel := models.AzureDNSRecord{
		Name: name,
		Type: "AAAA",
		Properties: models.AzureDnsRecordProperties{
			TTL: int64(ttl),
		},
	}

	for _, ipv6 := range ipv6s {
		r := models.AzureDnsAAAARecord{
			Ipv6Address: ipv6,
		}
		putModel.Properties.AAAARecords = append(putModel.Properties.AAAARecords, r)
	}

	_, err := svc.UpsertDnsRecord(resourceGroup, zoneName, putModel)

	return err
}

func (svc *AzureService) UpsertMXDnsRecord(resourceGroup string, zoneName string, ttl int, exchange ...models.AzureDnsMXRecord) error {
	if len(exchange) == 0 {
		return errors.New("you need at least one exchange server")
	}

	putModel := models.AzureDNSRecord{
		Name: "@",
		Type: "MX",
		Properties: models.AzureDnsRecordProperties{
			TTL: int64(ttl),
		},
	}

	for _, srv := range exchange {
		putModel.Properties.MXRecords = append(putModel.Properties.MXRecords, srv)
	}

	_, err := svc.UpsertDnsRecord(resourceGroup, zoneName, putModel)

	return err
}

func (svc *AzureService) UpsertCNAMEDnsRecord(resourceGroup string, zoneName string, name string, ttl int, alias string) error {
	if alias == "" {
		return errors.New("you need at least one alias")
	}

	putModel := models.AzureDNSRecord{
		Name: name,
		Type: "CNAME",
		Properties: models.AzureDnsRecordProperties{
			TTL: int64(ttl),
			CNAMERecord: &models.AzureDnsCNAMERecord{
				Cname: alias,
			},
		},
	}

	_, err := svc.UpsertDnsRecord(resourceGroup, zoneName, putModel)

	return err
}

func (svc *AzureService) UpsertSRVDnsRecord(resourceGroup string, zoneName string, name string, ttl int, srvs ...models.AzureDnsSRVRecord) error {
	if len(srvs) == 0 {
		return errors.New("you need at least one server record")
	}

	putModel := models.AzureDNSRecord{
		Name: name,
		Type: "SRV",
		Properties: models.AzureDnsRecordProperties{
			TTL: int64(ttl),
		},
	}

	for _, srv := range srvs {
		putModel.Properties.SRVRecords = append(putModel.Properties.SRVRecords, srv)
	}

	_, err := svc.UpsertDnsRecord(resourceGroup, zoneName, putModel)

	return err
}

func (svc *AzureService) UpsertTXTDnsRecord(resourceGroup string, zoneName string, name string, ttl int, value string) error {
	if value == "" {
		return errors.New("TXT record cannot have a empty value")
	}

	putModel := models.AzureDNSRecord{
		Name: name,
		Type: "TXT",
		Properties: models.AzureDnsRecordProperties{
			TTL: int64(ttl),
		},
	}

	if len(value) <= 254 {
		putModel.Properties.TXTRecords = append(putModel.Properties.TXTRecords, models.AzureDnsTXTRecord{
			Value: []string{value},
		})
	} else {
		txt := models.AzureDnsTXTRecord{
			Value: make([]string, 0),
		}

		sub := ""
		for i, v := range value {
			sub = fmt.Sprintf("%v%v", sub, v)
			if i > 254 {
				txt.Value = append(txt.Value, sub)
				sub = ""
			}
		}
		putModel.Properties.TXTRecords = txt
	}

	_, err := svc.UpsertDnsRecord(resourceGroup, zoneName, putModel)

	return err
}
