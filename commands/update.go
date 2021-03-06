package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/sapcc/arc/updater"
	"github.com/sapcc/arc/version"
)

func Update(c *cli.Context, options map[string]interface{}) (int, error) {
	if c.String("update-uri") == "" {
		return 1, errors.New("not enough arguments in call update command. Option update-uri is missing")
	}

	up := updater.New(map[string]string{
		"version":   version.Version,
		"appName":   options["appName"].(string),
		"updateUri": c.String("update-uri"),
	})

	r, err := up.Check()
	if err == updater.ErrorNoUpdateAvailable {
		fmt.Println("No update available")
		return 0, nil
	} else if err != nil {
		return 1, err
	}

	if r == nil {
		fmt.Println("No update response available.")
		return 0, nil
	}

	fmt.Printf("Available update version %s\n", r.Version)

	if !c.Bool("no-update") {
		if !c.Bool("force") {
			// ask for update
			fmt.Printf("Would you like to update to the version %q (yes/no):\n", r.Version)
			confirm, err := askForConfirmation()
			if err != nil {
				return 1, err
			}
			if !confirm {
				return 0, nil
			}
		}
		// update
		err = up.Update(r)
		if err != nil {
			return 1, err
		}
	}
	return 0, nil
}

// private

/*
 * askForConfirmation uses Scanln to parse user input
 */
var confirmInput = os.Stdin

func askForConfirmation() (bool, error) {
	var response string
	_, err := fmt.Fscanf(confirmInput, "%s", &response)
	if err != nil {
		return false, err
	}
	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	if containsString(okayResponses, response) {
		return true, nil
	} else if containsString(nokayResponses, response) {
		return false, nil
	} else {
		fmt.Println("Please type yes or no and then press enter:")
		return askForConfirmation()
	}
}

/*
 * posString returns the first index of element in slice.
 * If slice does not contain element, returns -1.
 */
func posString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}

/*
 * containsString returns true iff slice contains element
 */
func containsString(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
}
