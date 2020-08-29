package common

import (
	"context"
	"fmt"

	cloudflare "github.com/cloudflare/cloudflare-go"
)



func ReadCmdKV(request CmdKVRRequest) ([]byte, error) {
	api, err := cloudflare.New(ApiKey, EmailUser, cloudflare.UsingAccount(AccountID))
	if err != nil {
		fmt.Println(err)
	}

	cmdPayload, err := api.ReadWorkersKV(context.Background(), request.NameSpace.UUID, request.CmdKName)
	if err != nil {
		return nil, err
	}

	return cmdPayload, nil
}

func WriteOutKV(request OutKVWRequest) (cloudflare.Response, error) {
	api, err := cloudflare.New(ApiKey, EmailUser, cloudflare.UsingAccount(AccountID))
	if err != nil {
		fmt.Println(err)
	}

	resp, err := api.WriteWorkersKV(context.Background(), request.NameSpace.UUID, request.OutKName, request.OutPayload)
	return resp, err
}

func WriteCmdKV(request CmdKVWRequest) (cloudflare.Response, error) {
	api, err := cloudflare.New(ApiKey, EmailUser, cloudflare.UsingAccount(AccountID))
	if err != nil {
		return cloudflare.Response{}, err
	}

	// Write command payload
	resp, err := api.WriteWorkersKV(context.Background(), request.NameSpace.UUID, request.CmdKName, request.CmdPayload)
	if err != nil {
		return cloudflare.Response{}, err
	}

	return resp, err
}


func ReadOutKV(request OutKVRRequest) ([]byte, error) {
	api, err := cloudflare.New(ApiKey, EmailUser, cloudflare.UsingAccount(AccountID))
	if err != nil {
		fmt.Println(err)
	}

	resp, err := api.ReadWorkersKV(context.Background(), request.NameSpace.UUID, request.OutKName)
	if err != nil {
		return nil, err
	}
	return resp, nil
}


func ListWorkersKVs(namespace string)( cloudflare.ListStorageKeysResponse, error) {
	api, err := cloudflare.New(ApiKey, EmailUser, cloudflare.UsingAccount(AccountID))
	if err != nil {
		return cloudflare.ListStorageKeysResponse{}, err
	}

	resp, err := api.ListWorkersKVs(context.Background(), namespace)
	if err != nil {
		return cloudflare.ListStorageKeysResponse{}, err
	}

	return resp, nil
}
