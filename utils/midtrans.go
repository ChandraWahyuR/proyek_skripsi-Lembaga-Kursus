package utils

import "github.com/veritrans/go-midtrans"

func NewMidtransClient(serverKey, clientKey string) midtrans.Client {
	midclient := midtrans.NewClient()
	midclient.ServerKey = serverKey
	midclient.ClientKey = clientKey
	midclient.APIEnvType = midtrans.Sandbox

	return midclient
}
