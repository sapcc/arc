package updater

import (
	log "github.com/Sirupsen/logrus"
	"github.com/inconshreveable/go-update"
	"github.com/inconshreveable/go-update/check"
	"runtime"
)

type Updater struct {
	params    check.Params
	updateUri string
}

/**
 * Basic options to initialize the struct
 * options["version"] = if version is 0, it will be set to 1 (check.go)
 * options["appName"] = identifier of the application to update
 * options["updateUri"] = update server uri
 */
func New(options map[string]string) *Updater {
	log.Infof("Updater setup with version %q, app name %q and update uri %q", options["version"], options["appName"], options["updateUri"])
	return &Updater{
		params: check.Params{
			AppVersion: options["version"],
			AppId:      options["appName"],
			OS:         runtime.GOOS,
		},
		updateUri: options["updateUri"],
	}
}

func (u *Updater) Update() (bool, error) {
	// update obj
	up := update.New()

	// check for the update
	r, err := u.params.CheckForUpdate(u.updateUri, up)
	if err == check.NoUpdateAvailable {
		// no content means no available update, http 204
		log.Errorf("No update available")
		return false, err
	} else if err != nil {
		log.Errorf("Error while checking for update: %q", err.Error())
		return false, err
	}
	log.Infof("Updated version %q for app %q available ", r.Version, u.params.AppId)

	// apply the update
	err = applyUpdate(r)
	if err != nil {
		return false, err
	}

	return true, nil
}

// private

var applyUpdate = apply_update

func apply_update(r *check.Result) error {
	err, _ := r.Update()
	if err != nil {
		log.Errorf("Failed to update: %q", err.Error())
		return err
	}
	log.Infof("Updated to version %q", r.Version)

	return nil
}
