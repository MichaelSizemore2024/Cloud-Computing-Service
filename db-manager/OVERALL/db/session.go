package db

import (
	configs "dbmanager/OVERALL/config"
	"dbmanager/OVERALL/log"
	"sync"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"go.uber.org/zap"
)

// Global Vars
var (
	// Lazy initialization
	once sync.Once
	// Interface
	sessionWrapperService SessionWrapperService
)

// Defines interface above
type SessionWrapperService interface {
	Query(stmt string, names []string) QueryXService
}

// Wrapper for session
// maybe can be removed? seems redudant but maybe needed for lazy
type SessionWrapperImpl struct {
	Session *gocqlx.Session
}

// Defines Query used by interface
func (s *SessionWrapperImpl) Query(stmt string, names []string) QueryXService {
	queryX := s.Session.Query(stmt, names)
	return &QueryXServiceImpl{
		QueryX: queryX,
	}
}

func GetSession(configs configs.DbConfigs) SessionWrapperService {
	//Initializes db session using things from configs
	once.Do(func() {
		scyllaDbCluster := gocql.NewCluster(configs.DBHosts...)
		// I think this needs to be adjusted to allow for different user provided keyspaces?
		// Not sure how well that will work considering this is in a once.Do
		scyllaDbCluster.Keyspace = configs.DBKeyspace
		scyllaDbCluster.ConnectTimeout = configs.DBConnectTimeout
		scyllaDbCluster.WriteTimeout = configs.DBWriteTimeout
		scyllaDbCluster.Timeout = configs.DBReadTimeout
		scyllaDbCluster.NumConns = configs.DBConnectionsPerHost
		scyllaDbCluster.PoolConfig = gocql.PoolConfig{
			HostSelectionPolicy: gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy()),
		}
		scyllaDbCluster.ReconnectionPolicy = &gocql.ConstantReconnectionPolicy{
			MaxRetries: configs.DBReConnectionMaxRetries,
			Interval:   configs.DBReConnectionInterval,
		}
		scyllaDbCluster.SocketKeepalive = configs.DBKeepAliveTime
		scyllaDbCluster.Authenticator = gocql.PasswordAuthenticator{
			Username: configs.DBUsername,
			Password: configs.DBPassword,
		}
		// Creates Session
		s, err := gocqlx.WrapSession(scyllaDbCluster.CreateSession())
		if err != nil {
			log.Logger.Fatal("failed to create scylla db session ", zap.Error(err))
		}
		// Wraps and returns session
		sessionWrapperService = &SessionWrapperImpl{
			Session: &s,
		}
	})
	return sessionWrapperService
}
