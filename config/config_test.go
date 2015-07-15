package config_test

import (
	. "gitHub.***REMOVED***/monsoon/arc/config"

	"github.com/codegangsta/cli"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"flag"
)

var _ = Describe("Config", func() {

	It("should retrieve variables from the context", func() {
		// prepare flags
		globalSet := flag.NewFlagSet("test", 0)
		globalSet.String("transport", "mqtt", "test")
		globalSet.String("tls-client-cert", "", "test")
		globalSet.String("tls-client-key", "", "test")
		globalSet.String("tls-ca-cert", "", "test")
		globalSet.String("log-level", "info", "test")

		stringSlice := cli.StringSlice{}
		stringSlice.Set("tcp://localhost:1883")
		flag := cli.StringSliceFlag{Name: "endpoint", Value: &stringSlice}
		flag.Apply(globalSet)
		ctx := cli.NewContext(nil, nil, globalSet)

		// load context to the config
		conf := Config{}
		err := conf.Load(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(ctx.GlobalString("transport")).To(Equal("mqtt"))
		Expect(ctx.GlobalString("tls-client-cert")).To(Equal(""))
		Expect(ctx.GlobalString("tls-client-key")).To(Equal(""))
		Expect(ctx.GlobalString("tls-ca-cert")).To(Equal(""))
		Expect(ctx.GlobalString("log-level")).To(Equal("info"))
		Expect(ctx.GlobalStringSlice("endpoint")).To(Equal([]string{"tcp://localhost:1883"}))
	})

})
