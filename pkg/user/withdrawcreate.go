package user

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/user"
)

func CreateWithdrawAccount(
	ctx context.Context,
	appID, userID, coinTypeID, address string,
	labels []string,
) (
	*npool.Account, error,
) {
	return nil, nil
}
