package config

import (
	"log"
	"time"

	"github.com/gocql/gocql"
)

func ConnectToCassandra(keyspace string) *gocql.Session {
	cluster := gocql.NewCluster("localhost:9042")
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.One // Change from QUORUM to One
	cluster.ProtoVersion = 4
	cluster.ConnectTimeout = time.Second * 10

	session, err := cluster.CreateSession()
	if err != nil {
		log.Printf("Error creating session: %v", err)
		return nil
	}

	return session
}
