package watcher

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGettingEvents(t *testing.T) {
	assert := assert.New(t)
	var envVersion = os.Getenv("POSTGRES_PASSWORD")
	fmt.Printf("envVersion: %s\n", envVersion)
	fmt.Printf("postgresPassword: %s\n", postgresPassword)

	assert.Equal(postgresPassword, envVersion)
}
