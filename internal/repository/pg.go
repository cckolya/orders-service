package repository

import (
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"orders-service/config"
	"time"
)

var schema = `
create table if not exists orders (
    added_at        TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    order_uid       TEXT,
    raw_message     TEXT,
    invalid_message TEXT
    );
`

type Pg struct {
	log *zerolog.Logger
	db  *sqlx.DB
}

func NewPg(log *zerolog.Logger, cfg config.Postgres) (*Pg, error) {
	connectionURL := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)

	open, err := sqlx.Open("pgx", connectionURL)
	if err != nil {
		return nil, err
	}
	err = open.Ping()
	if err != nil {
		return nil, err
	}

	open.SetMaxOpenConns(cfg.Settings.MaxOpenConns)
	open.SetConnMaxLifetime(cfg.Settings.ConnMaxLifeTime * time.Second)
	open.SetMaxIdleConns(cfg.Settings.MaxIdleConns)
	open.SetConnMaxIdleTime(cfg.Settings.MaxIdleLifeTime * time.Second)

	pg := &Pg{db: open, log: log}

	pg.mustInitSchema()
	if err != nil {
		return nil, err
	}
	return pg, nil
}

func (p *Pg) mustInitSchema() {
	p.db.MustExec(schema)
}

func (p *Pg) CreateValidMsg(orderUID string, data []byte) error {
	_, err := p.db.Exec(`insert into orders (order_uid, raw_message) values ($1, $2)`, orderUID, data)
	return err
}

func (p *Pg) CreateInvalidMsg(data []byte) error {
	_, err := p.db.Exec(`insert into orders (invalid_message) values ($1)`, data)
	return err
}

func (p *Pg) GetOrderByID(uid string) ([]byte, error) {
	var arr []byte
	err := p.db.Get(&arr, `select raw_message from orders where order_uid = $1`, uid)
	if err != nil {
		return nil, err
	}
	return arr, nil
}

func (p *Pg) GetValidOrders(limit int) (map[string]interface{}, error) {
	var orders []struct {
		OrderUid   string `db:"order_uid"`
		RawMessage []byte `db:"raw_message"`
	}
	err := p.db.Select(&orders, `select order_uid, raw_message from orders where order_uid is not null limit $1`, limit)
	if err != nil {
		return nil, err
	}

	ret := make(map[string]interface{}, len(orders))
	for _, order := range orders {
		ret[order.OrderUid] = order.RawMessage
	}

	return ret, nil
}
