package commands

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	"github.com/oklog/run"
	"github.com/sapcc/arc/api-server/pki"
	arc_config "github.com/sapcc/arc/config"
	"github.com/sapcc/arc/server"
	"github.com/sapcc/arc/transport"
	"github.com/sapcc/arc/updater"
	"github.com/sapcc/arc/version"
)

var (
	errShutdown         = errors.New("shutdown")
	errGracefulShutdown = errors.New("gracefulShutdown")
)

// CmdServer starts a node server
func CmdServer(c *cli.Context, cfg arc_config.Config, appName string) (int, error) {
	// init broker transport
	tp, err := transport.New(cfg, true)
	if err != nil {
		return 1, err
	}
	if err = tp.Connect(); err != nil {
		return 1, fmt.Errorf("failed to connect to broker: %s", err)
	}

	// init server
	server := server.New(cfg, tp)

	// init routine handler
	var runner run.Group
	{
		// Server Actor
		runner.Add(func() error {
			defer logend(logstart("arc server"))
			log.Infof("running arc server with version %s. identity: %s, project: %s and organization: %s", version.Version, cfg.Identity, cfg.Project, cfg.Organization)
			return server.Run()
		}, func(err error) {
			log.Infof("server actor was interrupted with: %v", err)
			if err == errGracefulShutdown {
				server.GracefulShutdown()
			} else {
				server.Stop()
			}
		})
	}
	{
		// Set-up our signal Actor
		cancelSignalHandler := make(chan struct{})
		runner.Add(func() error {
			defer logend(logstart("signal handler"))
			return signalHandler(cancelSignalHandler)
		}, func(err error) {
			log.Infof("Signal actor was interrupted with: %v", err)
			close(cancelSignalHandler)
		})
	}

	// update binary Actor
	if c.String("update-uri") != "" {
		cancelVersionUpdaterChan := make(chan struct{})
		runner.Add(func() error {
			defer logend(logstart("version updater"))
			log.Infof("runing version updater with interval %v, version %q, app name %q and update uri %q", c.Int("update-interval"), version.Version, appName, c.String("update-uri"))
			return runVersionUpdater(c.Int("update-interval"), appName, c.String("update-uri"), cancelVersionUpdaterChan)
		}, func(err error) {
			log.Infof("update binary was interrupted with: %v", err)
			close(cancelVersionUpdaterChan)
		})
	}

	// update cert Actor
	renewCertURI, err := RenewCertURI(c)
	if err != nil {
		log.Errorf("failed to get renew cert URI: %s \n", err)
	} else {
		defer logend(logstart("cert updater"))
		log.Infof("running cert updater with URI %s, interval %v minutes and threshold %v hours", renewCertURI, c.Int("cert-update-interval"), c.Int("cert-update-threshold"))
		cancelCertHandler := make(chan struct{})
		runner.Add(func() error {
			return runCertUpdater(renewCertURI, c.Int("cert-update-interval"), c.Int("cert-update-threshold"), cfg, cancelCertHandler)
		}, func(err error) {
			log.Infof("cert actor was interrupted with: %v", err)
			close(cancelCertHandler)
		})
	}

	return 1, runner.Run()
}

func runCertUpdater(renewCertURI string, renewCertInterval int, renewCertThreshold int, cfg arc_config.Config, cancel <-chan struct{}) error {
	updaterTickChan := time.NewTicker(time.Minute * time.Duration(renewCertInterval))
	defer updaterTickChan.Stop()

	for {
		select {
		case <-updaterTickChan.C:
			hoursLeft, err := pki.CheckAndRenewCert(&cfg, renewCertURI, int64(renewCertThreshold), false)
			if err != nil {
				log.Error(err)
			} else {
				// when hours left is 0 the new cert is downlaoded
				if hoursLeft == 0 {
					log.Infof("cert updater downloaded new cert, gracefulShutdown")
					return errGracefulShutdown
				}
				log.Infof("cert updater skipped, %v hours to expiration", hoursLeft)
			}
		case <-cancel:
			return nil
		}
	}
}

func runVersionUpdater(interval int, appName string, updateURI string, cancel <-chan struct{}) error {
	up := updater.New(map[string]string{
		"version":   version.Version,
		"appName":   appName,
		"updateUri": updateURI,
	})
	updaterTickChan := time.NewTicker(time.Second * time.Duration(interval))
	defer updaterTickChan.Stop()

	for {
		select {
		case <-updaterTickChan.C:
			success, err := up.CheckAndUpdate()
			if err != nil {
				log.Error(err)
			}
			if success {
				return errGracefulShutdown
			}
		case <-cancel:
			return nil
		}
	}
}

func signalHandler(cancel <-chan struct{}) error {
	gracefulChan := make(chan os.Signal, 1)
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(gracefulChan, syscall.SIGTERM)
	signal.Notify(shutdownChan, syscall.SIGINT)

	select {
	case sig := <-shutdownChan:
		log.Infof("received signal %s", sig)
		return errShutdown
	case sig := <-gracefulChan:
		log.Infof("received signal %s", sig)
		return errGracefulShutdown
	case <-cancel:
		return nil
	}
}

func logstart(what string) string {
	log.Println("Starting ", what)
	return what
}
func logend(what string) {
	log.Println("Stopped ", what)
}
