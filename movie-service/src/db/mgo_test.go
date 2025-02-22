package db

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	conn = MongoReplicaSet{}
	u, uok := os.LookupEnv("DB_USER")
	p, pok := os.LookupEnv("DB_PASS")
	s, sok := os.LookupEnv("DB_SERVERS")
	n, nok := os.LookupEnv("DB_NAME")
	r, rok := os.LookupEnv("DB_REPLICA")

	if !uok {
		fmt.Println("No DB user specified")
	}
	if !pok {
		fmt.Println("No DB pass specified")
	}
	if !sok {
		fmt.Println("No DB url specified")
	}
	if !nok {
		fmt.Println("No DB name specified")
	}
	if !rok {
		fmt.Println("No DB replica specified")
	}

	conn.User = u
	conn.Pass = p
	conn.Servers = s
	conn.Db = n
	conn.ReplicaSet = r
	conn.AuthSource = "authSource=admin"
}

// TestURLMany ...
func TestURLMany(t *testing.T) {
	conn := make(chan *MongoConnection)

	go MongoDB(conn)

	session := <-conn

	assert.NoError(t, session.Err)

	defer session.Session.Close()

	result := struct{ Ok int }{}
	err := session.Session.Run("ping", &result)

	assert.NoError(t, err)
	assert.Equal(t, result.Ok, 1)
}
