package transfer

import (
	"context"
	"fmt"

	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/transfer"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"

	transfermwcli "github.com/NpoolPlatform/account-middleware/pkg/client/transfer"
	transfermwpb "github.com/NpoolPlatform/message/npool/account/mw/v1/transfer"
)

type deleteHandler struct {
	*Handler
	infos []*transfermwpb.Transfer
	users map[string]*usermwpb.User
	accs  []*npool.Transfer
}

func (h *deleteHandler) formalize() {
	for _, val := range h.infos {
		userInfo, ok := h.users[val.TargetUserID]
		if !ok {
			continue
		}

		h.accs = append(h.accs, &npool.Transfer{
			ID:                 val.ID,
			AppID:              val.AppID,
			UserID:             val.UserID,
			TargetUserID:       val.TargetUserID,
			TargetEmailAddress: userInfo.EmailAddress,
			TargetPhoneNO:      userInfo.PhoneNO,
			CreatedAt:          val.CreatedAt,
			TargetUsername:     userInfo.Username,
			TargetFirstName:    userInfo.FirstName,
			TargetLastName:     userInfo.LastName,
		})
	}
}

func (h *deleteHandler) getUsers(ctx context.Context) error {
	targetUserIDs := []string{}
	for _, val := range h.infos {
		targetUserIDs = append(targetUserIDs, val.TargetUserID)
	}

	users, _, err := usermwcli.GetUsers(ctx, &usermwpb.Conds{
		EntIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: targetUserIDs},
	}, 0, int32(len(targetUserIDs)))
	if err != nil {
		return err
	}

	for _, val := range users {
		h.users[val.EntID] = val
	}
	return nil
}

func (h *Handler) DeleteTransfer(ctx context.Context) (*npool.Transfer, error) {
	if h.ID == nil {
		return nil, fmt.Errorf("invalid id")
	}
	if h.AppID == nil {
		return nil, fmt.Errorf("invalid appID")
	}
	if h.UserID == nil {
		return nil, fmt.Errorf("invalid userID")
	}
	exist, err := transfermwcli.ExistTransferConds(ctx, &transfermwpb.Conds{
		ID:     &basetypes.StringVal{Op: cruder.EQ, Value: *h.ID},
		AppID:  &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		UserID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.UserID},
	})
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("transfer not exist")
	}
	info, err := transfermwcli.DeleteTransfer(ctx, &transfermwpb.TransferReq{
		ID: h.ID,
	})
	if err != nil {
		return nil, err
	}

	handler := &deleteHandler{
		Handler: h,
		infos:   []*transfermwpb.Transfer{info},
		users:   map[string]*usermwpb.User{},
	}
	if err := handler.getUsers(ctx); err != nil {
		return nil, err
	}

	handler.formalize()

	return handler.accs[0], nil
}
