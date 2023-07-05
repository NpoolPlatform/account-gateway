package user

import (
	"context"
	"fmt"

	constant "github.com/NpoolPlatform/account-gateway/pkg/const"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	"github.com/google/uuid"
)

type Handler struct {
	ID               *string
	AppID            *string
	UserID           *string
	CoinTypeID       *string
	Address          *string
	UsedFor          *basetypes.AccountUsedFor
	Labels           *[]string
	Account          *string
	AccountType      *basetypes.SignMethod
	VerificationCode *string
	Memo             *string
	Backup           *bool
	Active           *bool
	Blocked          *bool
	Offset           int32
	Limit            int32
}

func NewHandler(ctx context.Context, options ...func(context.Context, *Handler) error) (*Handler, error) {
	handler := &Handler{}
	for _, opt := range options {
		if err := opt(ctx, handler); err != nil {
			return nil, err
		}
	}
	return handler, nil
}

func WithID(id *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		h.ID = id
		return nil
	}
}

func WithAppID(id *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		h.AppID = id
		return nil
	}
}

func WithUserID(id *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		h.UserID = id
		return nil
	}
}

func WithCoinTypeID(id *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		// TODO: check coin exist
		h.CoinTypeID = id
		return nil
	}
}

func WithAddress(addr *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if addr == nil {
			return nil
		}
		if *addr == "" {
			return fmt.Errorf("invalid address")
		}
		h.Address = addr
		return nil
	}
}

func WithUsedFor(usedFor *basetypes.AccountUsedFor) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if usedFor == nil {
			return nil
		}
		switch *usedFor {
		case basetypes.AccountUsedFor_UserWithdraw:
		case basetypes.AccountUsedFor_UserDirectBenefit:
		default:
			return fmt.Errorf("invalid usedfor")
		}
		h.UsedFor = usedFor
		return nil
	}
}

func WithLabels(labels *[]string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if labels == nil {
			return nil
		}
		h.Labels = labels
		return nil
	}
}

func WithAccount(account *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if account == nil {
			return nil
		}
		if *account == "" {
			return fmt.Errorf("invalid account")
		}
		h.Account = account
		return nil
	}
}

func WithAccountType(accountType *basetypes.SignMethod) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if accountType == nil {
			return nil
		}
		switch *accountType {
		case basetypes.SignMethod_Email:
		case basetypes.SignMethod_Mobile:
		case basetypes.SignMethod_Google:
		default:
			return fmt.Errorf("invalid accountType")
		}
		h.AccountType = accountType
		return nil
	}
}

func WithVerificationCode(verifCode *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if verifCode == nil {
			return nil
		}
		if *verifCode == "" {
			return fmt.Errorf("invalid verificationCode")
		}
		h.VerificationCode = verifCode
		return nil
	}
}

func WithMemo(meno *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if meno == nil {
			return nil
		}
		if *meno == "" {
			return fmt.Errorf("invalid meno")
		}
		h.Memo = meno
		return nil
	}
}

func WithBackup(backup *bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.Backup = backup
		return nil
	}
}

func WithBlocked(blocked *bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.Blocked = blocked
		return nil
	}
}

func WithActive(active *bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.Active = active
		return nil
	}
}

func WithOffset(offset int32) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.Offset = offset
		return nil
	}
}

func WithLimit(limit int32) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if limit == 0 {
			limit = constant.DefaultLimit
		}
		h.Limit = limit
		return nil
	}
}
