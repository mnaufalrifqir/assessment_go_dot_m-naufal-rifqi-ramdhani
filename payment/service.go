package payment

import (
	"api-dot/user"
	"api-dot/utils"

	"strconv"

	midtrans "github.com/veritrans/go-midtrans"
)

type Service interface {
	GetPaymentURL(transaction Transaction, user user.User) (string, error)
}

type service struct {
}

func NewService() *service {
	return  &service{}
}

func (s *service) GetPaymentURL(transaction Transaction, user user.User) (string, error){
	midclient := midtrans.NewClient()
    midclient.ServerKey = utils.GetConfig("MIDTRANS_SERVER_KEY")
    midclient.ClientKey = utils.GetConfig("MIDTRANS_CLIENT_KEY")
    midclient.APIEnvType = midtrans.Sandbox

    snapGateway := midtrans.SnapGateway{
        Client: midclient,
    }

	snapReq := &midtrans.SnapReq{
		CustomerDetail: &midtrans.CustDetail{
			Email: user.Email,
			FName: user.Name,
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID: strconv.Itoa(int(transaction.ID)),
			GrossAmt: int64(transaction.Amount),
		},
	}

	snapTokenResp, err := snapGateway.GetToken(snapReq)
	if err != nil{
		return "", err
	}

	return snapTokenResp.RedirectURL, nil
}


