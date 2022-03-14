// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package cfg

import (
	"encoding/json"
	"errors"
	"flag"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigyaml"
	"gopkg.in/yaml.v2"
)

const (
	AppName   = "omesrv"
	CfgFile   = "omesrv.yaml"
	EnvPrefix = "OME"
)

// ServerConfig represents global OME server configuration.
type ServerConfig struct {
	Domain string     `env:"APP_DOMAIN" usage:"application’s primary domain"`
	Appurl string     `env:"APP_URL" usage:"application’s primary domain in http format (e.g. https://my-domain.com)"`
	Apiver string     `env:"API_VER" default:"v1" usage:"current API version"`
	Http   HttpConfig //nolint
	Grpc   GrpcConfig
	Pgdb   PgdbConfig
	Redis  RedisConfig
	Log    LogConfig
	flags  *flag.FlagSet   // command line arguments
	field  []aconfig.Field // configured fields
	files  []string        // configuration files
	dump   func()
	exit   bool // must exit
	penv   bool // print environment flag
}

var _cfg ServerConfig

var (
	Server = &_cfg
	Http   = &_cfg.Http //nolint
	Grpc   = &_cfg.Grpc
	Redis  = &_cfg.Redis
)

// args mirrors os.Args sans program name.
//
// Used for testing.
var args = os.Args[1:]

const (
	confFlag = "conf"
	envFlag  = "env"
)

func init() {
	for i := range args {
		if strings.HasPrefix(args[i], "-test") {
			args = args[:0]
			return
		}
	}
}

// Log is a shortcut for the ServerConfig.Log.
func Log() *LogConfig {
	return &_cfg.Log
}

// Pgdb is a shortcut for the ServerConfig.Pgdb.
func Pgdb() *PgdbConfig {
	return &_cfg.Pgdb
}

// Load performs loading of ServerConfig from a file,
// environment, and command line arguments.
func Load() error {
	return _cfg.Load()
}

// Load performs loading of ServerConfig from a file,
// environment, and command line arguments.
func (c *ServerConfig) Load() error {
	loader, err := c.createConfigLoader()
	if err != nil {
		return err
	}
	err = loader.Load()
	if err == nil {
		if c.penv {
			c.printEnviron()
			return nil
		}
		if c.dump != nil {
			c.exit = true
			c.dump()
			return nil
		}
		return c.Check()
	}
	if errors.Is(err, flag.ErrHelp) {
		c.printHelp()
		return nil
	}
	if cause := errors.Unwrap(err); cause != nil {
		return cause
	}
	return err
}

// Check validates ServerConfig fields.
func (c *ServerConfig) Check() (err error) {
	if len(c.Apiver) == 0 {
		return EmptyError("apiver")
	}
	var check checkList
	check.add("http", &c.Http)
	check.add("grpc", &c.Grpc)
	check.add("pgdb", &c.Pgdb)
	check.add("redis", &c.Redis)
	check.add("log", &c.Log)
	return check.run()
}

// AddConfigFile adds custom file path to the configuration search
// locations.
func (c *ServerConfig) AddConfigFile(file string) {
	c.files = append(c.files, file)
}

// AddConfigDir adds a directory to the configuration file search
// locations. Appended with default CfgFile name.
func (c *ServerConfig) AddConfigDir(dir string) {
	c.files = append(c.files, filepath.Join(dir, CfgFile))
}

// MustExit indicates that the process should exit. It is toggled in the
// Load function.
//
// Most often it is true when help requested from the
// command line arguments.
func (c *ServerConfig) MustExit() bool {
	return c.exit
}

// ReleaseMemory clears internal fields after configuration has
// been processed.
func (c *ServerConfig) ReleaseMemory() {
	c.flags = nil
	c.field = nil
	c.files = nil
}

//-----------------------------------------------------------------------------

func (c *ServerConfig) createConfigLoader() (*aconfig.Loader, error) {
	files, err := configSearch(CfgFile, c.files)
	if err != nil {
		return nil, err
	}
	loader := aconfig.LoaderFor(c, aconfig.Config{
		SkipDefaults:       false,
		SkipFiles:          false,
		SkipEnv:            false,
		SkipFlags:          false,
		EnvPrefix:          EnvPrefix,
		FlagPrefix:         "",
		FlagDelimiter:      ".",
		AllFieldRequired:   false,
		AllowDuplicates:    false,
		AllowUnknownFields: false,
		AllowUnknownEnvs:   true,
		AllowUnknownFlags:  true,
		DontGenerateTags:   false,
		FailOnFileNotFound: false,
		MergeFiles:         false,
		Args:               args,
		FileFlag:           confFlag,
		Files:              files,
		FileDecoders: map[string]aconfig.FileDecoder{
			".yaml": aconfigyaml.New(),
		},
	})
	c.flags = loader.Flags()
	c.flags.Usage = func() {}
	c.flags.SetOutput(io.Discard)
	c.flags.BoolVar(&c.penv, envFlag, false, "print environment variables")
	c.flags.Func("dump", "print configuration before validation check [#yaml|json#]", c.dumpFlag)
	c.field = make([]aconfig.Field, 0, 32)
	loader.WalkFields(func(f aconfig.Field) bool {
		c.field = append(c.field, f)
		return true
	})
	return loader, nil
}

//-----------------------------------------------------------------------------

func (c *ServerConfig) dumpFlag(kind string) (err error) {
	const unknownFormat = ConstError("unknown dump format")
	switch kind {
	case "yaml":
		c.dump = c.PrintYAML
	case "json":
		c.dump = c.PrintJSON
	default:
		err = unknownFormat
	}
	return
}

//-----------------------------------------------------------------------------
// JSON
//-----------------------------------------------------------------------------

// WriteJSON writes JSON representation of ServerConfig to w.
//
// The resulting JSON is formatted with indents and new lines.
func (c *ServerConfig) WriteJSON(w io.Writer) error {
	b, err := json.MarshalIndent(c, "", "  ")
	if err == nil {
		_, err = w.Write(b)
	}
	return err
}

// PrintJSON outputs ServerConfig JSON to os.Stdout.
func (c *ServerConfig) PrintJSON() {
	_ = c.WriteJSON(os.Stdout)
}

//-----------------------------------------------------------------------------
// YAML
//-----------------------------------------------------------------------------

// WriteYAML writes YAML representation of ServerConfig to w.
func (c *ServerConfig) WriteYAML(w io.Writer) error {
	b, err := yaml.Marshal(c)
	if err == nil {
		_, err = w.Write(b)
	}
	return err
}

// PrintYAML outputs ServerConfig YAML to os.Stdout.
func (c *ServerConfig) PrintYAML() {
	_ = c.WriteYAML(os.Stdout)
}

// GetEnv lookups environment variable prefixed with EnvPrefix.
//
// Returns ok if variable found and has non-zero length.
func GetEnv(key string) (val string, ok bool) {
	const pfx = EnvPrefix + "_"
	if !strings.HasPrefix(key, pfx) {
		key = pfx + key
	}
	val, ok = os.LookupEnv(key)
	ok = ok && len(val) > 0
	return
}
