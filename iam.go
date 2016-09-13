package main

import (
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iam/v1"
)

type IamAPI struct {
	api *iam.Service
}

func NewIamAPI() (*IamAPI, error) {
	client, err := google.DefaultClient(context.Background(), iam.CloudPlatformScope)
	if err != nil {
		return nil, err
	}
	iamService, err := iam.New(client)
	if err != nil {
		return nil, err
	}
	return &IamAPI{iamService}, nil
}

func (a *IamAPI) GetLatestKey(projectId, account string) ([]byte, error) {
	keySrv := iam.NewProjectsServiceAccountsKeysService(a.api)
	listRes, err := keySrv.List(account).Do()
	if err != nil {
		return []byte{}, err
	}
	keys := listRes.Keys
	if len(keys) != 0 {
		return keys[0].MarshalJSON()
		var validKey *iam.ServiceAccountKey
		now := time.Now()
		for _, key := range keys {
			valid := true
			if key.ValidAfterTime != "" {
				t, err := time.Parse(time.RFC3339Nano, key.ValidAfterTime)
				if err != nil {
					return nil, err
				}
				valid = t.After(now)
			}
			if key.ValidBeforeTime != "" {
				t, err := time.Parse(time.RFC3339Nano, key.ValidBeforeTime)
				if err != nil {
					return nil, err
				}
				valid = t.Before(now)
			}
			if valid {
				validKey = key
				break
			}
		}
		if validKey != nil {
			return validKey.MarshalJSON()
		}
	}
	req := &iam.CreateServiceAccountKeyRequest{
		PrivateKeyType: "TYPE_GOOGLE_CREDENTIALS_FILE",
	}
	createRes, err := keySrv.Create(account, req).Do()
	if err != nil {
		return []byte{}, err
	}
	return createRes.MarshalJSON()
}
