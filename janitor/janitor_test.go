package janitor

import (
	"encoding/json"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bamzi/jobrunner"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Janitor", func() {

	var (
		janitor = &Janitor{}
	)

	JustBeforeEach(func() {
		SCHEDULE_TIME = "* * * * * ?"

		env := os.Getenv("ARC_ENV")
		if env == "" {
			env = "test"
		}

		// init janitor
		conf := JanitorConf{
			DbConfigFile: "../api-server/db/dbconf.yml",
			Environment:  env,
		}
		janitor = InitJanitor(conf)
	})

	AfterEach(func() {
		janitor.Stop()
	})

	It("should run all clean jobs", func() {
		routines := "FailQueuedJobs FailExpiredJobs PruneJobs AggregateLogs PruneLocks PruneCertificates"
		janitor.InitScheduler()
		time.Sleep(1 * time.Second)

		Expect(len(jobrunner.Entries())).To(Equal(6))

		for _, element := range jobrunner.Entries() {
			// convert data
			data, err := json.Marshal(element.Job)
			Expect(err).NotTo(HaveOccurred())
			jobrunnerJob := jobrunner.Job{}
			err = json.Unmarshal([]byte(data), &jobrunnerJob)
			Expect(err).NotTo(HaveOccurred())
			Expect(parseValue(jobrunnerJob.Latency)).To(BeNumerically(">", 0.0))
			Expect(strings.Contains(routines, jobrunnerJob.Name)).To(BeTrue())
		}

	})

})

func parseValue(payload string) float64 {
	var validValue = regexp.MustCompile(`\d{1,}[.]\d{1,}|\d{1,}`)
	// get the first value of the string
	strArray := validValue.FindAllString(payload, 1)
	if len(strArray) > 0 {
		// parse to float
		value, err := strconv.ParseFloat(strArray[0], 64)
		if err == nil {
			return value
		}
	}
	return 0
}
