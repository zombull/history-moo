package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/zombull/choo-choo/bug"

	"gopkg.in/yaml.v2"
)

type Config struct {
	// Config is the location of the user's config file.  This value is not
	// saved/exposed as the location is defined via an environment variable
	// to avoid a Catch-22.
	Config string `yaml:"-"`

	// Database is the path to the directory where the SQLite database
	// exists (or is created).
	// Env Var: CHOO_CHOO_DATABASE
	// Default: $HOME/.db
	Database string `yaml:"database"`

	// Server is the path to the root directory of the web server.
	// Env Var: CHOO_CHOO_DATABASE
	// Default: $HOME/Development/go/src/github.com/zombull/choo-choo/server
	Server string `yaml:"server"`

	// Specify the Moonboard set, i.e. year and holds combination.
	// Env Var: CHOO_CHOO_MOONBOARD_SET
	// Default: "2016 A+B+O"
	MoonboardSet string `yaml:"moonboard_set"`

	// Specify the Moonboard user (display name), e.g. Sean Christopherson
	// Env Var: CHOO_CHOO_MOONBOARD_USER
	// Default: "Sean Christopherson"
	MoonboardUser string `yaml:"moonboard_user"`

	// Github is the URL of the Github repository for the database.  This
	// is required for using the sync command, but is optional otherwise.
	// Env Var: CHOO_CHOO_GITHUB
	// Default: ""
	Github string `yaml:"github"`
}

func loadEnvVar(name, def string) string {
	name = "CHOO_CHOO_" + name
	if os.Getenv(name) != "" {
		def = os.Getenv(name)
	}
	return os.ExpandEnv(def)
}

// LoadConfig reads the configuration from the config path; if the path does
// not exist, it returns a default configuration.
func loadConfig() *Config {
	// Use a default config if a user-defined file does not exist.
	// Basic Windows (not MinGW or MSysGit) may not have $HOME set,
	// look for HOMEDRIVE and HOMEPATH.
	dir := "$HOME"
	if os.Getenv("HOME") == "" && os.Getenv("HOMEDRIVE") != "" && os.Getenv("HOMEPATH") != "" {
		dir = path.Join(os.Getenv("HOMEDRIVE"), os.Getenv("HOMEPATH"))
	}
	c := Config{
		Config:        path.Join(dir, ".config", "choo-choo", "config.yml"),
		Database:      path.Join(dir, ".db"),
		Server:        path.Join(dir, "Development", "go", "src", "github.com", "zombull", "choo-choo", "server"),
		MoonboardSet:  "2016 A+B+O",
		MoonboardUser: "Sean Christopherson",
		Github:        "",
	}

	path := loadEnvVar("CONFIG", c.Config)
	data, err := ioutil.ReadFile(path)
	if err == nil {
		err = yaml.Unmarshal(data, &c)
		bug.OnError(err)
	} else {
		bug.On(!os.IsNotExist(err), fmt.Sprintf("cannot read config file: %v", err))
	}

	c.Database = loadEnvVar("DATABASE", c.Database)
	c.Server = loadEnvVar("SERVER", c.Server)
	c.MoonboardSet = loadEnvVar("MOONBOARD_SET", c.MoonboardSet)
	c.MoonboardUser = loadEnvVar("MOONBOARD_USER", c.MoonboardUser)
	c.Github = loadEnvVar("GITHUB", c.Github)

	bug.On(len(c.Database) == 0, "database must be a non-empty string")
	bug.On(len(c.Server) == 0, "database must be a non-empty string")

	err = os.MkdirAll(c.Database, 0770)
	bug.OnError(err)

	return &c
}
