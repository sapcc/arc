package main

import (
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/gorilla/handlers"

	"gitHub.***REMOVED***/monsoon/arc/update-server/storage"
	"gitHub.***REMOVED***/monsoon/arc/version"
)

var cliCommands = []cli.Command{
	{
		Name:   "local",
		Usage:  "Local storage",
		Action: localStorage,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "path,p",
				Usage:  "Directory containig update artifacts",
				Value:  "/var/lib/arc-update-site",
				EnvVar: "ARTIFACTS_PATH",
			},
		},
	},
	{
		Name:   "swift",
		Usage:  "Swift storage",
		Action: swiftStorage,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "username,u",
				Usage:  "User name for the swift authentication",
				EnvVar: "OS_USERNAME",
			},
			cli.StringFlag{
				Name:   "password,p",
				Usage:  "Password for the swift authentication",
				EnvVar: "OS_PASSWORD",
			},
			cli.StringFlag{
				Name:   "domain,d",
				Usage:  "Domain for the swift authentication",
				EnvVar: "OS_USER_DOMAIN_NAME",
			},
			cli.StringFlag{
				Name:   "auth_url,a",
				Usage:  "Authentication URL for the swift authentication",
				EnvVar: "OS_AUTH_URL",
			},
			cli.StringFlag{
				Name:   "container,c",
				Usage:  "The Swift container",
				EnvVar: "OS_CONTAINER",
			},
		},
	},
}

func localStorage(c *cli.Context) {
	var err error

	// check mandatory params
	buildsRootPath := c.String("path")
	if buildsRootPath == "" {
		log.Fatal("No path to update artifacts given.")
		return
	}
	if err = os.MkdirAll(buildsRootPath, 0755); err != nil {
		log.Fatalf("Path to artificats %s does not exist and can't be created: %s", buildsRootPath, err)
		return
	}
	log.Infof("Serving artifacts from %s.", buildsRootPath)

	// set the storage
	st, err = storage.New(storage.Local, c)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	// run the server
	runServer(c, storage.Local)
}

func swiftStorage(c *cli.Context) {
	var err error

	// check mandatory params
	if c.String("username") == "" || c.String("password") == "" || c.String("domain") == "" || c.String("auth_url") == "" || c.String("container") == "" {
		log.Fatal("Not enough arguments in call swift command")
		return
	}
	log.Infof("Serving artifacts from Swift container %s.", c.String("container"))

	// set the storage
	st, err = storage.New(storage.Swift, c)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	// run the server
	runServer(c, storage.Swift)
}

func runServer(c *cli.Context, storageType storage.StorageType) {
	log.Infof("Starting update server version %s.", version.String())

	// cache the templates
	templates = getTemplates()

	// get the router
	router := newRouter(storageType)

	// run server
	log.Infof("Listening on %q...", c.GlobalString("bind-address"))
	accessLogger := handlers.CombinedLoggingHandler(os.Stdout, router)
	if err := http.ListenAndServe(c.GlobalString("bind-address"), accessLogger); err != nil {
		log.Fatalf("Failed to bind on %s: %s", c.GlobalString("bind-address"), err)
	}
}
