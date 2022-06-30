package testing

import (
	"database/sql"
	"log"
	"sync"
	"sync/atomic"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"

	"github.com/migratemgr8/mgr8/global"
	"github.com/migratemgr8/mgr8/testing/databaseconfigs"
)

type dockerManager struct {
	pool      *dockertest.Pool
	configs   map[global.Database]DatabaseConfig
	resources map[global.Database]*dockertest.Resource
	calls     int64
}

var m *dockerManager

var initializeDockerManagerOnce sync.Once

func NewDockerManager() *dockerManager {
	initializeDockerManagerOnce.Do(initializeDockerManager)
	atomic.AddInt64(&m.calls, 1)
	return m
}

func initializeDockerManager() {
	m = &dockerManager{calls: 0}
	var err error
	m.pool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker: %v", err)
	}
	m.pool.MaxWait = 120 * time.Second

	m.configs = make(map[global.Database]DatabaseConfig)
	m.configs[global.Postgres] = databaseconfigs.NewPostgresConfig()
	m.configs[global.MySql] = databaseconfigs.NewMySqlConfig()

	m.resources = make(map[global.Database]*dockertest.Resource)
	for database, databaseConfig := range m.configs {
		m.resources[database], err = m.pool.RunWithOptions(databaseConfig.DockerOptions(), func(config *docker.HostConfig) {
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{Name: "no"}
		})
		if err != nil {
			log.Fatalf("Could not start resource: %v", err)
		}
		err = m.resources[database].Expire(120)
		if err != nil {
			log.Fatalf("Docker expiration config failed: %v", err)
		}

		var db *sql.DB
		if err = m.pool.Retry(func() error {
			db, err = sql.Open(database.String(), databaseConfig.DatabaseUrl(m.resources[database]))
			if err != nil {
				return err
			}
			return db.Ping()
		}); err != nil {
			log.Fatalf("Could not connect to docker: %s", err)
		}
		err = db.Close()
		if err != nil {
			log.Fatalf("Could not close health check connection: %v", err)
		}
	}
}

func (m *dockerManager) GetConnectionString(d global.Database) string {
	return m.configs[d].DatabaseUrl(m.resources[d])
}

func (m *dockerManager) CloseAll() error {
	atomic.AddInt64(&m.calls, -1)
	if m.calls == 0 {
		for _, r := range m.resources {
			if err := m.pool.Purge(r); err != nil {
				return err
			}
		}
	}
	return nil
}
