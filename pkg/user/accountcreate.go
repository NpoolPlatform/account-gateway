package user

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/go-service-framework/pkg/pubsub"
	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/user"
	eventmwpb "github.com/NpoolPlatform/message/npool/inspire/mw/v1/event"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	coininfocli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin"

	addresscheck "github.com/NpoolPlatform/account-gateway/pkg/addresscheck"

	sphinxproxypb "github.com/NpoolPlatform/message/npool/sphinxproxy"
	sphinxproxycli "github.com/NpoolPlatform/sphinx-proxy/pkg/client"

	useraccmwcli "github.com/NpoolPlatform/account-middleware/pkg/client/user"
	useraccmwpb "github.com/NpoolPlatform/message/npool/account/mw/v1/user"

	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	usercodemwcli "github.com/NpoolPlatform/basal-middleware/pkg/client/usercode"
	usercodemwpb "github.com/NpoolPlatform/message/npool/basal/mw/v1/usercode"
)

type createHandler struct {
	*Handler
	coinName            *string
	checkAddressBalance bool
}

func (h *createHandler) rewardSetAddress() {
	if err := pubsub.WithPublisher(func(publisher *pubsub.Publisher) error {
		req := &eventmwpb.CalcluateEventRewardsRequest{
			AppID:       *h.AppID,
			UserID:      *h.UserID,
			EventType:   basetypes.UsedFor_SetAddress,
			Consecutive: 1,
		}
		return publisher.Update(
			basetypes.MsgID_CalculateEventRewardReq.String(),
			nil,
			nil,
			nil,
			req,
		)
	}); err != nil {
		logger.Sugar().Errorw(
			"rewardSetAddress",
			"AppID", *h.AppID,
			"UserID", h.UserID,
			"Error", err,
		)
	}
}

func (h *createHandler) validate(ctx context.Context) error { //nolint
	switch *h.AccountType {
	case basetypes.SignMethod_Email:
		fallthrough //nolint
	case basetypes.SignMethod_Mobile:
		if h.Account == nil || *h.Account == "" {
			return fmt.Errorf("account is empty")
		}
	case basetypes.SignMethod_Google:
	default:
		return fmt.Errorf("accounttype %v invalid", *h.AccountType)
	}

	if *h.VerificationCode == "" {
		return fmt.Errorf("invalid verificationcode")
	}
	if *h.Address == "" {
		return fmt.Errorf("invalid address")
	}

	switch *h.UsedFor {
	case basetypes.AccountUsedFor_UserWithdraw:
	case basetypes.AccountUsedFor_UserDirectBenefit:
	default:
		return fmt.Errorf("usedfor %v invalid", *h.UsedFor)
	}

	return nil
}

func (h *createHandler) checkVerifyUserCode(ctx context.Context) error {
	user, err := usermwcli.GetUser(ctx, *h.AppID, *h.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("invalid user")
	}

	if *h.AccountType == basetypes.SignMethod_Google {
		h.Account = &user.GoogleSecret
	}

	if err := usercodemwcli.VerifyUserCode(ctx, &usercodemwpb.VerifyUserCodeRequest{
		Prefix:      basetypes.Prefix_PrefixUserCode.String(),
		AppID:       *h.AppID,
		Account:     *h.Account,
		AccountType: *h.AccountType,
		UsedFor:     basetypes.UsedFor_SetWithdrawAddress,
		Code:        *h.VerificationCode,
	}); err != nil {
		return err
	}

	return nil
}

func (h *createHandler) getCoinName(ctx context.Context) error {
	coin, err := coininfocli.GetCoin(ctx, *h.CoinTypeID)
	if err != nil {
		return err
	}
	if coin == nil {
		return fmt.Errorf("invlaid coin")
	}
	h.coinName = &coin.Name
	h.checkAddressBalance = coin.CheckNewAddressBalance

	return nil
}

func (h *createHandler) validateAddress(ctx context.Context) error {
	if h.coinName == nil {
		return fmt.Errorf("invalid coinname")
	}

	if !h.checkAddressBalance {
		err := addresscheck.ValidateAddress(*h.coinName, *h.Address)
		if err != nil {
			return fmt.Errorf("invalid %v address", *h.coinName)
		}
		return nil
	}

	bal, err := sphinxproxycli.GetBalance(ctx, &sphinxproxypb.GetBalanceRequest{
		Name:    *h.coinName,
		Address: *h.Address,
	})
	if err != nil {
		return err
	}
	if bal == nil {
		return fmt.Errorf("invalid address")
	}

	return nil
}

func (h *createHandler) createAccount(ctx context.Context) error {
	info, err := useraccmwcli.CreateAccount(ctx, &useraccmwpb.AccountReq{
		AppID:      h.AppID,
		UserID:     h.UserID,
		CoinTypeID: h.CoinTypeID,
		Address:    h.Address,
		UsedFor:    h.UsedFor,
		Labels:     h.Labels,
		Memo:       h.Memo,
	})
	if err != nil {
		return err
	}
	if info == nil {
		return fmt.Errorf("fail create account")
	}

	h.ID = &info.ID
	h.EntID = &info.EntID
	return nil
}

func (h *Handler) CreateAccount(ctx context.Context) (*npool.Account, error) {
	handler := &createHandler{
		Handler: h,
	}
	if err := handler.validate(ctx); err != nil {
		return nil, err
	}
	if err := handler.checkVerifyUserCode(ctx); err != nil {
		return nil, err
	}
	if err := handler.getCoinName(ctx); err != nil {
		return nil, err
	}
	if err := handler.validateAddress(ctx); err != nil {
		return nil, err
	}
	if err := handler.createAccount(ctx); err != nil {
		return nil, err
	}

	handler.rewardSetAddress()

	return h.GetAccount(ctx)
}
