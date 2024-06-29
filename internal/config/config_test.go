package config_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"go-load-balancer/internal/config"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}

var _ = Describe("Config", func() {
	Context("LoadConfig", func() {
		It("should load valid configuration", func() {
			file, err := os.CreateTemp("", "config.json")
			Expect(err).NotTo(HaveOccurred())
			defer os.Remove(file.Name())

			_, err = file.WriteString(`{
                "servers": ["http://localhost:8081", "http://localhost:8082"],
                "strategy": "round-robin"
            }`)
			Expect(err).NotTo(HaveOccurred())
			file.Close()

			cfg, err := config.LoadConfig(file.Name())
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.Servers).To(ConsistOf("http://localhost:8081", "http://localhost:8082"))
			Expect(cfg.Strategy).To(Equal("round-robin"))
		})

		It("should return an error for invalid configuration", func() {
			file, err := os.CreateTemp("", "invalid_config.json")
			Expect(err).NotTo(HaveOccurred())
			defer os.Remove(file.Name())

			_, err = file.WriteString(`{invalid_json}`)
			Expect(err).NotTo(HaveOccurred())
			file.Close()

			_, err = config.LoadConfig(file.Name())
			Expect(err).To(HaveOccurred())
		})
	})
})
