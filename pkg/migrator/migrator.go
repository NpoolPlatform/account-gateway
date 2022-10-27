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
	infos, err := cli1.
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
		for _, info := range infos {
			_, err := tx.
				Account.
				Create().
				SetID(info.ID).
				SetCoinTypeID(info.CoinTypeID).
				SetPlatformHoldPrivateKey(info.PlatformHoldPrivateKey).
				Save(_ctx)
			if err != nil {
				return err
			}
		}

		return nil
	})
	for _, info := range infos {
		logger.Sugar().Infow("migrateAccount", "info", info)
	}

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

	return nil
}
