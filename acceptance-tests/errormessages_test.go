package acceptance_test

import (
	"csbbrokerpakgcp/acceptance-tests/helpers/brokers"
	"csbbrokerpakgcp/acceptance-tests/helpers/cf"
	"csbbrokerpakgcp/acceptance-tests/helpers/random"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Error Messages", Label("errormessages"), func() {
	When("the create-service command fails immediately", func() {
		It("prints a useful error message", func() {
			name := random.Name(random.WithPrefix("error"))
			defer cf.Run("delete-service", "-f", name)

			session := cf.Start("create-service", "csb-google-storage-bucket", "private", name, "-b", brokers.DefaultBrokerName(), "-c", `{"storage_class":"bogus"}`)
			Eventually(session, time.Minute).Should(Exit(1))
			Expect(session.Out).To(Say(`FAILED\n`))
			Expect(session.Err).To(Say(`Service broker error: 1 error\(s\) occurred: storage_class: storage_class must be one of the following: "COLDLINE", "MULTI_REGIONAL", "NEARLINE", "REGIONAL", "STANDARD"`))
		})
	})

	When("the service creation fail asynchronously", func() {
		It("puts a useful error message in the service description", func() {
			name := random.Name(random.WithPrefix("error"))
			defer cf.Run("delete-service", "-f", name)

			session := cf.Start("create-service", "csb-google-storage-bucket", "private", name, "-b", brokers.DefaultBrokerName(), "-c", `{"project":"not-real-project"}`)
			Eventually(session, time.Minute).Should(Exit(0))

			Eventually(func() string {
				stdout, _ := cf.Run("service", name)
				return stdout
			}, 10*time.Minute, 10*time.Second).Should(MatchRegexp(`status:\s+create failed`))

			stdout, _ := cf.Run("service", name)
			Expect(stdout).To(MatchRegexp(`message:\s+Error: googleapi: Error 400: Unknown project id: not-real-project, invalid  with google_storage_bucket.bucket,  on main.tf line 1, in resource "google_storage_bucket" "bucket":   1: resource "google_storage_bucket" "bucket"`))
		})
	})
})