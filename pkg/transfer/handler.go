package transfer

import (
	"context"
	"fmt"

	constant "github.com/NpoolPlatform/account-gateway/pkg/const"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	"github.com/google/uuid"
)

type Handler struct {
	ID                *uint32
	EntID             *string
	AppID             *string
	UserID            *string
	Account           *string
	AccountType       *basetypes.SignMethod
	VerificationCode  *string
	TargetUserID      *string
	TargetAccount     *string
	TargetAccountType *basetypes.SignMethod
	Offset            int32
	Limit             int32
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

func WithID(id *uint32, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid id")
			}
			return nil
		}
		h.ID = id
		return nil
	}
}

func WithEntID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid entid")
			}
			return nil
		}
		_, err := uuid.Parse(*id)
		if err != nil {
			return err
		}
		h.EntID = id
		return nil
	}
}

func WithAppID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid appid")
			}
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		h.AppID = id
		return nil
	}
}

func WithUserID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid userid")
			}
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		h.UserID = id
		return nil
	}
}

func WithAccount(account *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if account == nil {
			if must {
				return fmt.Errorf("invalid account")
			}
			return nil
		}
		if *account == "" {
			return fmt.Errorf("invalid account")
		}
		h.Account = account
		return nil
	}
}

func WithAccountType(accountType *basetypes.SignMethod, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if accountType == nil {
			if must {
				return fmt.Errorf("invalid accounttype")
			}
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

func WithVerificationCode(verifCode *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if verifCode == nil {
			if must {
				return fmt.Errorf("invalid verificationcode")
			}
			return nil
		}
		if *verifCode == "" {
			return fmt.Errorf("invalid verificationCode")
		}
		h.VerificationCode = verifCode
		return nil
	}
}

func WithTargetUserID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid targetuserid")
			}
			return nil
		}
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
		h.TargetUserID = id
		return nil
	}
}

func WithTargetAccount(account *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if account == nil {
			if must {
				return fmt.Errorf("invalid targetaccount")
			}
			return nil
		}
		if *account == "" {
			return fmt.Errorf("invalid targetAccount")
		}
		h.TargetAccount = account
		return nil
	}
}

func WithTargetAccountType(accountType *basetypes.SignMethod, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if accountType == nil {
			if must {
				return fmt.Errorf("invalid targetaccounttype")
			}
			return nil
		}
		switch *accountType {
		case basetypes.SignMethod_Email:
		case basetypes.SignMethod_Mobile:
		case basetypes.SignMethod_Google:
		default:
			return fmt.Errorf("invalid targetAccountType")
		}
		h.TargetAccountType = accountType
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
