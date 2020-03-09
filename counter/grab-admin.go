package counter

import (
	"database/sql"
	"os"
	"strings"
	"sync"

	"gitlab.ifchange.com/data/cordwood/log"
	"gitlab.ifchange.com/data/cordwood/mysql"
)

var (
	grabAdminOnceSetupMySQL sync.Once
	grabAdminMySQLConn      *sql.DB
	grabAdminWorkStatus     bool

	nilGrabAdminCounter = new(grabAdminRPC)
)

func NewGrabAdminRPCCounter(kind Kind, siteID int) Counter {
	grabAdminOnceSetupMySQL.Do(func() {
		if os.Getenv("ENV") != "prod" {
			return
		}
		cfg := mysql.Config{
			Username:     "ts",
			Password:     "sdfe232t9ddde3d",
			Address:      "192.168.8.238",
			Port:         "3306",
			DB:           "grabadmin",
			KeepAlive:    10,
			MaxOpenConns: 10,
			MaxIdleConns: 10,
		}
		conn, err := mysql.Connect(cfg)
		if err != nil {
			log.Warnf("cordwood counter try construct grab-admin database's connection failed, useful in PROD only. %v", err)
			return
		}
		grabAdminMySQLConn = conn
		grabAdminWorkStatus = true
	})
	if !grabAdminWorkStatus {
		return nilGrabAdminCounter
	}
	return &grabAdminRPC{
		kind:   kind,
		siteID: siteID,
	}
}

var (
	_ Counter = (*grabAdminRPC)(nil)
)

type grabAdminRPC struct {
	kind   Kind
	siteID int
}

func (counter *grabAdminRPC) Inc(status bool, msg ...string) {
	if !grabAdminWorkStatus {
		return
	}
	var (
		dbStatus int = statusFail
	)
	if status {
		dbStatus = statusSuccess
	}

	_, err := grabAdminMySQLConn.Exec(`INSERT INTO rpc_statistics
		(kind,site_id,status,msg) VALUES (?,?,?,?)`,
		int(counter.kind), counter.siteID, dbStatus, strings.Join(msg, " "))

	if err != nil {
		log.Errorf("cordwood counter insert rpc_statistics %v", err)
	}
}
