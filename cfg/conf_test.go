// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package cfg

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

//-----------------------------------------------------------------------------

func TestServerConfig_Load(t *testing.T) {
	loadConfig(t)
}

//-----------------------------------------------------------------------------

func TestServerConfig_LoadEnv(t *testing.T) {
	t.Setenv("OME_HTTP_PORT", "80")
	t.Setenv("OME_GRPC_PORT", "90")
	t.Setenv("OME_REDIS_STORE_POOL", "32")
	t.Setenv("OME_REDIS_STORE_PASS", "SECRET")
	loadConfig(t)
}

//-----------------------------------------------------------------------------

func TestServerConfig_LoadArg(t *testing.T) {
	setArg("-http.port", "100", "-grpc.port", "110", "-redis.store.pass", "secret")
	loadConfig(t)
}

//-----------------------------------------------------------------------------

func TestServerConfig_Help(t *testing.T) {
	if !testing.Verbose() {
		return
	}
	setArg("-h")
	err := Server.Load()
	require.NoError(t, err)
}

//-----------------------------------------------------------------------------

func TestServerConfig_Env(t *testing.T) {
	if !testing.Verbose() {
		return
	}
	setArg("-env")
	err := Server.Load()
	require.NoError(t, err)
}

//-----------------------------------------------------------------------------

func TestServerConfig_PrintJson(t *testing.T) {
	logJSON(t, sampleConfig())
}

func TestServerConfig_PrintYaml(t *testing.T) {
	logYAML(t, sampleConfig())
}

//-----------------------------------------------------------------------------

func TestConfigSearch(t *testing.T) {
	all, err := configSearch(CfgFile, nil)
	if err != nil {
		t.Fatal(err)
	}
	logJSON(t, all)
}

//-----------------------------------------------------------------------------

func TestParseBytes(t *testing.T) {
	n, err := strconv.Atoi("100")
	require.NoError(t, err)
	t.Logf("n = %d", n)
}

//-----------------------------------------------------------------------------
// Helpers
//-----------------------------------------------------------------------------

func loadConfig(t testing.TB) {
	err := Server.Load()
	require.NoError(t, err)
	if Server.MustExit() {
		return
	}
	logYAML(t, Server)
}

func sampleConfig() *ServerConfig {
	var c ServerConfig
	c.Http.Port = 8080
	c.Grpc.Port = 8090
	r := &c.Redis.Store
	r.Pool = 10
	r.Addr = "localhost:6379"
	r.Pass = "store pass"
	r = &c.Redis.Pubsub
	r.Pool = 20
	r.Addr = "127.0.0.1:6379"
	r.Pass = "pub sub pass"
	return &c
}

func logYAML(t testing.TB, v interface{}) {
	if testing.Verbose() {
		t.Logf("YAML:\n%s", yamlBytes(t, v))
	}
}

func logJSON(t testing.TB, v interface{}) {
	if testing.Verbose() {
		t.Logf("JSON:\n%s", jsonBytes(t, v))
	}
}

func yamlBytes(t testing.TB, v interface{}) []byte {
	b, err := yaml.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func jsonBytes(t testing.TB, v interface{}) []byte {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func setArg(a ...string) {
	if a == nil {
		a = []string{}
	}
	args = a
}
