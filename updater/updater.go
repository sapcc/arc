package updater

import (
	"fmt"
	"runtime"

	log "github.com/Sirupsen/logrus"
	"github.com/inconshreveable/go-update"
	"github.com/inconshreveable/go-update/check"
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
	return &Updater{
		params: check.Params{
			AppVersion: options["version"],
			AppId:      options["appName"],
			OS:         runtime.GOOS,
		},
		updateUri: options["updateUri"],
	}
}

/*
 * Check for new updates and replace binary
 */
func (u *Updater) CheckAndUpdate() (bool, error) {
	r, err := u.Check()
	if err == check.NoUpdateAvailable {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("Error while checking for update: %q", err.Error())
	}
	log.Infof("Updated version %q for app %q available ", r.Version, u.params.AppId)

	// replace binary
	err = u.Update(r)
	if err != nil {
		return false, fmt.Errorf("Failed to update: %q", err.Error())
	}
	log.Infof("Updated to version %q", r.Version)

	return true, nil
}

/*
 * Replace binary
 */
var ApplyUpdate = apply_update

func (u *Updater) Update(r *check.Result) error {
	err := ApplyUpdate(r)
	if err != nil {
		return err
	}
	return nil
}

/*
 * Check last version available
 */
func (u *Updater) Check() (*check.Result, error) {
	up := update.New()
	r, err := u.params.CheckForUpdate(u.updateUri, up)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// private

func apply_update(r *check.Result) error {
	err, _ := r.Update()
	if err != nil {
		return err
	}
	return nil
}
