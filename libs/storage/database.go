package storage

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Database struct {
	DB                    *sqlx.DB
	host                  string
	port                  int
	user                  string
	password              string
	name                  string
	maxOpenConnections    int
	maxIdleConnections    int
	connectionMaxLifetime time.Duration
	connectionMaxIdleTime time.Duration
	thumbnailsPrefix      string
}

// NewDatabaseOptions for NewDatabase.
type NewDatabaseOptions struct {
	ThumbnailPerfix       string
	Host                  string
	Port                  int
	User                  string
	Password              string
	Name                  string
	MaxOpenConnections    int
	MaxIdleConnections    int
	ConnectionMaxLifetime time.Duration
	ConnectionMaxIdleTime time.Duration
}

// NewDatabase with the given options.
// If no logger is provided, logs are discarded.
func NewDatabase(opts NewDatabaseOptions) *Database {
	return &Database{
		host:                  opts.Host,
		port:                  opts.Port,
		user:                  opts.User,
		password:              opts.Password,
		name:                  opts.Name,
		maxOpenConnections:    opts.MaxOpenConnections,
		maxIdleConnections:    opts.MaxIdleConnections,
		connectionMaxLifetime: opts.ConnectionMaxLifetime,
		connectionMaxIdleTime: opts.ConnectionMaxIdleTime,
		thumbnailsPrefix:      opts.ThumbnailPerfix,
	}
}

// Connect to the database.
func (d *Database) Connect() error {
	slog.Info("Connecting to database", slog.String("url", d.createDataSourceName(false)))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	d.DB, err = sqlx.ConnectContext(ctx, "pgx", d.createDataSourceName(true))
	if err != nil {
		return err
	}

	slog.Debug("Setting connection pool options",
		slog.Int("max open connections", d.maxOpenConnections),
		slog.Int("max idle connections", d.maxIdleConnections),
		slog.Duration("connection max lifetime", d.connectionMaxLifetime),
		slog.Duration("connection max idle time", d.connectionMaxIdleTime))
	d.DB.SetMaxOpenConns(d.maxOpenConnections)
	d.DB.SetMaxIdleConns(d.maxIdleConnections)
	d.DB.SetConnMaxLifetime(d.connectionMaxLifetime)
	d.DB.SetConnMaxIdleTime(d.connectionMaxIdleTime)

	return nil
}

func (d *Database) createDataSourceName(withPassword bool) string {
	password := d.password
	if !withPassword {
		password = "xxx"
	}
	return fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=disable", d.user, password, d.host, d.port, d.name)
}

func (d *Database) StorePath(ctx context.Context, path string, timestamp int64) error {
	slog.Info("Storing metadata Path")
	splitedString := strings.Split(path, "/")
	objName := splitedString[1]
	username := splitedString[0]

	insertUserIfNotExists := `
		insert into users (username)
		values ($1) ON CONFLICT (username) DO NOTHING
		RETURNING id  
		`
	var userID int
	d.DB.QueryRowxContext(ctx, insertUserIfNotExists, username).Scan(&userID)
	query := `
		insert into images (name, user_id)
		values ($1, $2)`
	_, err := d.DB.ExecContext(ctx, query, objName, userID)
	return err
}

func (d *Database) StorePathWithUser(ctx context.Context, user, path string, timestamp int64) error {
	return nil
}

func (d *Database) StoreLatest(ctx context.Context, index uint32, latest, objName string) error {
	return nil
}

func (d *Database) GetThumbnails(
	ctx context.Context,
	thumbnailCount int,
	keyword string) ([]string, error) {
	var urls []string
	var thumbnailURLs = make([]string, len(urls))
	err := d.DB.Select(&urls, "SELECT name FROM images limit $1", thumbnailCount)
	for i := range urls {
		thumbnailURLs[i] = d.thumbnailsPrefix + urls[i]
	}
	fmt.Println(urls)
	return thumbnailURLs, err
}
