package postgres

import (
	"context"
	"fmt"
	"strings"
	"test/config"
	"test/pkg/logger"
	"test/storage"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/golang-migrate/migrate/v4/database"          //database is needed for migration
	_ "github.com/golang-migrate/migrate/v4/database/postgres" //postgres is used for database
	_ "github.com/golang-migrate/migrate/v4/source/file"       //file is needed for migration url
	_ "github.com/lib/pq"
)

type Store struct {
	pool  *pgxpool.Pool
	log   logger.ILogger
	cfg   config.Config
	redis storage.IRedisStorage
}

func New(ctx context.Context, cfg config.Config, log logger.ILogger, redis storage.IRedisStorage) (storage.IStorage, error) {
	url := fmt.Sprintf(
		`postgres://%s:%s@%s:%s/%s?sslmode=disable`,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDB,
	)

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Error("error while parsing config", logger.Error(err))
		return nil, err
	}

	poolConfig.MaxConns = 100

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Error("error while connecting to db", logger.Error(err))
		return nil, err
	}

	//migration
	m, err := migrate.New("file://migrations/postgres/", url)
	if err != nil {
		log.Error("error while migrating", logger.Error(err))
		return nil, err
	}

	log.Info("???? came")

	if err = m.Up(); err != nil {
		log.Warning("migration up", logger.Error(err))
		if !strings.Contains(err.Error(), "no change") {
			fmt.Println("entered")
			version, dirty, err := m.Version()
			log.Info("version and dirty", logger.Any("version", version), logger.Any("dirty", dirty))
			if err != nil {
				log.Error("err in checking version and dirty", logger.Error(err))
				return nil, err
			}

			if dirty {
				version--
				if err = m.Force(int(version)); err != nil {
					log.Error("ERR in making force", logger.Error(err))
					return nil, err
				}
			}
			log.Warning("WARNING in migrating", logger.Error(err))
			return nil, err
		}
	}

	log.Info("!!!!! came here")

	return Store{
		pool:  pool,
		log:   log,
		cfg:   cfg,
		redis: redis,
	}, nil
}

func (s Store) Close() {
	s.pool.Close()
}

func (s Store) User() storage.IUserStorage {
	return NewUserRepo(s.pool, s.log, s.redis)
}

func (s Store) Category() storage.ICategoryStorage {
	return NewCategoryRepo(s.pool, s.log, s.redis)
}

func (s Store) Product() storage.IProductStorage {
	return NewProductRepo(s.pool, s.log, s.redis)
}

func (s Store) Basket() storage.IBasketStorage {
	return NewBasketRepo(s.pool, s.log, s.redis)

}

func (s Store) BasketProduct() storage.IBasketProductStorage {
	return NewBasketProductRepo(s.pool, s.log, s.redis)
}

func (s Store) Store() storage.IStoreStorage {
	return NewStoreRepo(s.pool, s.log, s.redis)
}

func (s Store) Branch() storage.IBranchStorage {
	return NewBranchRepo(s.pool, s.log, s.redis)
}

func (s Store) Dealer() storage.IDealerStorage {
	return NewDealerRepo(s.pool, s.log, s.redis)
}

func (s Store) Income() storage.IIncomeStorage {
	return NewIncomeRepo(s.pool, s.log, s.redis)
}

func (s Store) IncomeProduct() storage.IIncomeProductStorage {
	return NewIncomeProductRepo(s.pool, s.log, s.redis)
}

func (s Store) Redis() storage.IRedisStorage {
	return s.redis
}

func (s Store) Report() storage.IReportStorage {
	return NewReport(s.pool, s.log)
}