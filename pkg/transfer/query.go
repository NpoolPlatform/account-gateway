package transfer

import (
	"context"
	"fmt"

	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npool "github.com/NpoolPlatform/message/npool/account/gw/v1/transfer"

	"github.com/NpoolPlatform/message/npool/account/gw/v1/transfer"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"

	transfermwcli "github.com/NpoolPlatform/account-middleware/pkg/client/transfer"
	transfermwpb "github.com/NpoolPlatform/message/npool/account/mw/v1/transfer"
)

type queryHandler struct {
	*Handler
	infos []*transfermwpb.Transfer
	users map[string]*usermwpb.User
	accs  []*npool.Transfer
}

func (h *queryHandler) getUsers(ctx context.Context) error {
	targetUserIDs := []string{}
	for _, val := range h.infos {
		targetUserIDs = append(targetUserIDs, val.TargetUserID)
	}

	users, _, err := usermwcli.GetUsers(ctx, &usermwpb.Conds{
		IDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: targetUserIDs},
	}, 0, int32(len(targetUserIDs)))
	if err != nil {
		return err
	}

	for _, val := range users {
		h.users[val.ID] = val
	}
	return nil
}

func (h *queryHandler) formalize() {
	for _, val := range h.infos {
		userInfo, ok := h.users[val.TargetUserID]
		if !ok {
			continue
		}

		h.accs = append(h.accs, &transfer.Transfer{
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

func (h *Handler) GetTransfers(ctx context.Context) ([]*transfer.Transfer, uint32, error) {
	if h.AppID == nil {
		return nil, 0, fmt.Errorf("invalid appID")
	}
	if h.UserID == nil {
		return nil, 0, fmt.Errorf("invalid userID")
	}
	infos, total, err := transfermwcli.GetTransfers(
		ctx,
		&transfermwpb.Conds{
			AppID:  &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
			UserID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.UserID},
		},
		h.Offset,
		h.Limit,
	)
	if err != nil {
		return nil, 0, err
	}
	if len(infos) == 0 {
		return []*transfer.Transfer{}, 0, nil
	}

	handler := &queryHandler{
		Handler: h,
		infos:   infos,
		users:   map[string]*usermwpb.User{},
	}

	if err := handler.getUsers(ctx); err != nil {
		return nil, 0, err
	}

	handler.formalize()
	return handler.accs, total, nil
}

func (h *Handler) GetAppTransfers(ctx context.Context) ([]*transfer.Transfer, uint32, error) {
	if h.AppID == nil {
		return nil, 0, fmt.Errorf("invalid appID")
	}
	infos, total, err := transfermwcli.GetTransfers(
		ctx,
		&transfermwpb.Conds{
			AppID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID},
		},
		h.Offset,
		h.Limit,
	)
	if err != nil {
		return nil, 0, err
	}
	if len(infos) == 0 {
		return []*transfer.Transfer{}, 0, nil
	}

	handler := &queryHandler{
		Handler: h,
		infos:   infos,
		users:   map[string]*usermwpb.User{},
	}

	if err := handler.getUsers(ctx); err != nil {
		return nil, 0, err
	}

	handler.formalize()
	return handler.accs, total, nil
}