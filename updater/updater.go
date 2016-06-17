package updater

import (
	"encoding/hex"
	"fmt"
	"runtime"
	"sync/atomic"

	log "github.com/Sirupsen/logrus"
	update "github.com/inconshreveable/go-update"
)

type Updater struct {
	Params CheckParams
	client *Client
}

/**
 * Basic options to initialize the struct
 * options["version"] = if version is 0, it will be set to 1 (check.go)
 * options["appName"] = identifier of the application to update
 * options["updateUri"] = update server uri
 */
func New(options map[string]string) *Updater {
	client := NewClient(options["updateUri"])
	return &Updater{
		Params: CheckParams{
			AppVersion: options["version"],
			AppId:      options["appName"],
			OS:         runtime.GOOS,
			Arch:       runtime.GOARCH,
		},
		client: client,
	}
}

var checkAndUpdateRunning int32 = 0

/*
 * Check for new updates and replace binary
 */
func (u *Updater) CheckAndUpdate() (bool, error) {
	//Ensure only one CheckAndUpdate call is running at any given point in time
	if !atomic.CompareAndSwapInt32(&checkAndUpdateRunning, 0, 1) {
		return false, nil //Already running, bail
	}
	defer atomic.SwapInt32(&checkAndUpdateRunning, 0)

	r, err := u.Check()
	if err == NoUpdateAvailable {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("Error while checking for update: %q", err.Error())
	}
	log.Infof("Updated version %s for app %s available ", r.Version, u.Params.AppId)

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

func (u *Updater) Update(r *CheckResult) error {
	err := ApplyUpdate(u, r)
	return err
}

/*
 * Check last version available
 */
func (u *Updater) Check() (*CheckResult, error) {
	return u.client.CheckForUpdate(u.Params)
}

// private

func apply_update(u *Updater, r *CheckResult) error {
	reader, err := u.client.GetUpdate(r.Url)
	if err != nil {
		return err
	}
	defer (*reader).Close()

	//decode checksum
	checksum, err := hex.DecodeString(r.Checksum)
	if err != nil {
		return err
	}

	err = update.Apply(*reader, update.Options{Checksum: checksum})
	return err
}
