package updater

import (
	"time"
	log "github.com/Sirupsen/logrus"
	"github.com/inconshreveable/go-update"
	"github.com/inconshreveable/go-update/check"
)

type Updater struct {
	params    check.Params
	updateUri string
}

/**
 * options["version"] = if version is 0, it will be set to 1 (check.go)
 * options["appName"] = identifier of the application to update
 * options["updateUri"] = update server uri
 */
func New(options map[string]string) *Updater {
	log.Infof("Updater setup with version '%s', app name '%s' and update uri '%s'", options["version"], options["appName"], options["updateUri"])
	return &Updater{
		params: check.Params{
			AppVersion: options["version"],
			AppId:      options["appName"],
		},
		updateUri: options["updateUri"],
	}
}

func (u *Updater) Update(tickChan *time.Ticker) error {
	// update obj
	up := update.New()

	// check for the update
	r, err := u.params.CheckForUpdate(u.updateUri, up)
	if err == check.NoUpdateAvailable {
		// no content means no available update, http 204
		log.Errorf("No update available")
		return err
	} else if err != nil {
		log.Errorf("Error while checking for update: %v\n", err)
		return err
	}
	log.Infof("Updated version '%s' for app '%s' available ", r.Version, u.params.AppId)

	// apply the update
	err, _ = r.Update()
	if err != nil {
		log.Errorf("Failed to update: %v\n", err)
		return err
	}
	log.Infof("Updated to version %s!\n", r.Version)

	// Stop the ticker to not apply another update until the app is restarted
	tickChan.Stop()

	return nil
}