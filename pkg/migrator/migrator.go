package migrator

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	constant "github.com/NpoolPlatform/go-service-framework/pkg/mysql/const"

	billingent "github.com/NpoolPlatform/cloud-hashing-billing/pkg/db/ent"
	coinaccountinfoent "github.com/NpoolPlatform/cloud-hashing-billing/pkg/db/ent/coinaccountinfo"
	billingconst "github.com/NpoolPlatform/cloud-hashing-billing/pkg/message/const"

	accountmgrpb "github.com/NpoolPlatform/message/npool/account/mgr/v1/account"

	"github.com/NpoolPlatform/account-manager/pkg/db"
	"github.com/NpoolPlatform/account-manager/pkg/db/ent"
)

const (
	keyUsername = "username"
	keyPassword = "password"
	keyDBName   = "database_name"
	maxOpen     = 10
	maxIdle     = 10
	MaxLife     = 3
)

func dsn(hostname string) (string, error) {
	username := config.GetStringValueWithNameSpace(constant.MysqlServiceName, keyUsername)
	password := config.GetStringValueWithNameSpace(constant.MysqlServiceName, keyPassword)
	dbname := config.GetStringValueWithNameSpace(hostname, keyDBName)

	svc, err := config.PeekService(constant.MysqlServiceName)
	if err != nil {
		logger.Sugar().Warnw("dsb", "error", err)
		return "", err
	}

	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true&interpolateParams=true",
		username, password,
		svc.Address,
		svc.Port,
		dbname,
	), nil
}

func open(hostname string) (conn *sql.DB, err error) {
	hdsn, err := dsn(hostname)
	if err != nil {
		return nil, err
	}

	logger.Sugar().Infow("open", "hdsn", hdsn)

	conn, err = sql.Open("mysql", hdsn)
	if err != nil {
		return nil, err
	}

	// https://github.com/go-sql-driver/mysql
	// See "Important settings" section.

	conn.SetConnMaxLifetime(time.Minute * MaxLife)
	conn.SetMaxOpenConns(maxOpen)
	conn.SetMaxIdleConns(maxIdle)

	return conn, nil
}

var goodBenefits []*billingent.GoodBenefit
var goodPayments []*billingent.GoodPayment
var userWithdraws []*billingent.UserWithdraw
var coinSettings []*billingent.CoinSetting
var accounts []*billingent.CoinAccountInfo

//nolint
func accountUsedFor(ctx context.Context, id string, cli *billingent.Client) (accountmgrpb.AccountUsedFor, error) {
	var err error

	if len(goodBenefits) == 0 {
		goodBenefits, err = cli.
			GoodBenefit.
			Query().
			All(ctx)
		if err != nil {
			return accountmgrpb.AccountUsedFor_DefaultAccountUsedFor, err
		}
	}
	for _, b := range goodBenefits {
		if b.BenefitAccountID.String() == id {
			return accountmgrpb.AccountUsedFor_GoodBenefit, nil
		}
	}

	if len(goodPayments) == 0 {
		goodPayments, err = cli.
			GoodPayment.
			Query().
			All(ctx)
		if err != nil {
			return accountmgrpb.AccountUsedFor_DefaultAccountUsedFor, err
		}
	}
	for _, b := range goodPayments {
		if b.AccountID.String() == id {
			return accountmgrpb.AccountUsedFor_GoodPayment, nil
		}
	}

	if len(userWithdraws) == 0 {
		userWithdraws, err = cli.
			UserWithdraw.
			Query().
			All(ctx)
		if err != nil {
			return accountmgrpb.AccountUsedFor_DefaultAccountUsedFor, err
		}
	}
	for _, b := range userWithdraws {
		if b.AccountID.String() == id {
			return accountmgrpb.AccountUsedFor_UserWithdraw, nil
		}
	}

	if len(coinSettings) == 0 {
		coinSettings, err = cli.
			CoinSetting.
			Query().
			All(ctx)
		if err != nil {
			return accountmgrpb.AccountUsedFor_DefaultAccountUsedFor, err
		}
	}
	for _, b := range coinSettings {
		if b.PlatformOfflineAccountID.String() == id {
			return accountmgrpb.AccountUsedFor_PlatformBenefitCold, nil
		}
		if b.UserOfflineAccountID.String() == id {
			return accountmgrpb.AccountUsedFor_UserBenefitCold, nil
		}
		if b.UserOnlineAccountID.String() == id {
			return accountmgrpb.AccountUsedFor_UserBenefitHot, nil
		}
		if b.GoodIncomingAccountID.String() == id {
			return accountmgrpb.AccountUsedFor_PaymentCollector, nil
		}
		if b.GasProviderAccountID.String() == id {
			return accountmgrpb.AccountUsedFor_GasProvider, nil
		}
	}

	return accountmgrpb.AccountUsedFor_DefaultAccountUsedFor, nil
}

func migrateAccount(ctx context.Context, conn *sql.DB) error {
	cli, err := db.Client()
	if err != nil {
		return err
	}

	accs, err := cli.
		Account.
		Query().
		Limit(1).
		All(ctx)
	if err != nil {
		return err
	}
	if len(accs) > 0 {
		return nil
	}

	cli1 := billingent.NewClient(billingent.Driver(entsql.OpenDB(dialect.MySQL, conn)))
	accounts, err = cli1.
		CoinAccountInfo.
		Query().
		Where(
			coinaccountinfoent.DeleteAt(0),
		).
		All(ctx)
	if err != nil {
		logger.Sugar().Errorw("migrateAccount", "error", err)
		return err
	}

	err = db.WithTx(ctx, func(_ctx context.Context, tx *ent.Tx) error {
		for _, info := range accounts {
			usedFor, err := accountUsedFor(ctx, info.ID.String(), cli1)
			if err != nil {
				return err
			}
			if usedFor == accountmgrpb.AccountUsedFor_DefaultAccountUsedFor {
				continue
			}

			_, err = tx.
				Account.
				Create().
				SetID(info.ID).
				SetCoinTypeID(info.CoinTypeID).
				SetPlatformHoldPrivateKey(info.PlatformHoldPrivateKey).
				SetUsedFor(usedFor.String()).
				Save(_ctx)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func migrateGoodBenefit(ctx context.Context, conn *sql.DB) error {
	return nil
}

func migrateGoodPayment(ctx context.Context, conn *sql.DB) error {
	return nil
}

func migrateUserWithdraw(ctx context.Context, conn *sql.DB) error {
	return nil
}

func migrateCoinSetting(ctx context.Context, conn *sql.DB) error {
	return nil
}

func Migrate(ctx context.Context) error {
	if err := db.Init(); err != nil {
		logger.Sugar().Errorw("migrateAccount", "error", err)
		return err
	}

	billingConn, err := open(billingconst.ServiceName)
	if err != nil {
		logger.Sugar().Errorw("migrateAccount", "error", err)
		return err
	}
	defer billingConn.Close()

	if err := migrateAccount(ctx, billingConn); err != nil {
		return err
	}
	if err := migrateGoodBenefit(ctx, billingConn); err != nil {
		return err
	}
	if err := migrateGoodPayment(ctx, billingConn); err != nil {
		return err
	}
	if err := migrateUserWithdraw(ctx, billingConn); err != nil {
		return err
	}
	if err := migrateCoinSetting(ctx, billingConn); err != nil {
		return err
	}

	return nil
}
