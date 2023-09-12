package goutil

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/golangci/golangci-lint/pkg/logutils"
)

type EnvKey string

const (
	EnvGoCmd string = "GOLANGCI_LINT_GOCMD"

	EnvGoCache EnvKey = "GOCACHE"
	EnvGoRoot  EnvKey = "GOROOT"
)

type Env struct {
	vars   map[string]string
	log    logutils.Log
	debugf logutils.DebugFunc
}

func NewEnv(log logutils.Log) *Env {
	return &Env{
		vars:   map[string]string{},
		log:    log,
		debugf: logutils.Debug(logutils.DebugKeyEnv),
	}
}

func (e *Env) Discover(ctx context.Context) error {
	goCmd := os.Getenv(EnvGoCmd)
	if goCmd == "" {
		fmt.Println("Using default goCmd")
		goCmd = "go"
	} else {
		fmt.Println("Using " + goCmd)
	}

	startedAt := time.Now()
	args := []string{"env", "-json"}
	args = append(args, string(EnvGoCache), string(EnvGoRoot))
	out, err := exec.CommandContext(ctx, goCmd, args...).Output()
	if err != nil {
		return fmt.Errorf("failed to run 'go env': %w", err)
	}

	if err = json.Unmarshal(out, &e.vars); err != nil {
		return fmt.Errorf("failed to parse 'go %s' json: %w", strings.Join(args, " "), err)
	}

	e.debugf("Read go env for %s: %#v", time.Since(startedAt), e.vars)
	return nil
}

func (e Env) Get(k EnvKey) string {
	envValue := os.Getenv(string(k))
	if envValue != "" {
		return envValue
	}

	return e.vars[string(k)]
}
