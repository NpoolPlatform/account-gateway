package platform

import (
	"context"
	"fmt"

	constant "github.com/NpoolPlatform/account-gateway/pkg/const"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	"github.com/google/uuid"
)

type Handler struct {
	ID         *string
	CoinTypeID *string
	Address    *string
	UsedFor    *basetypes.AccountUsedFor
	Backup     *bool
	Active     *bool
	Blocked    *bool
	Locked     *bool
	Offset     int32
	Limit      int32
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

func WithCoinTypeID(id *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			return nil
		}
		// TODO: check coin exist
		if _, err := uuid.Parse(*id); err != nil {
			return err
		}
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
		case basetypes.AccountUsedFor_UserBenefitHot:
		case basetypes.AccountUsedFor_GasProvider:
		case basetypes.AccountUsedFor_PaymentCollector:
		case basetypes.AccountUsedFor_UserBenefitCold:
		case basetypes.AccountUsedFor_PlatformBenefitCold:
		default:
			return fmt.Errorf("invalid usedfor")
		}
		h.UsedFor = usedFor
		return nil
	}
}

func WithBackup(backup *bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.Backup = backup
		return nil
	}
}

func WithActive(active *bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.Active = active
		return nil
	}
}

func WithBlocked(blocked *bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.Blocked = blocked
		return nil
	}
}

func WithLocked(locked *bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.Locked = locked
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
