package broker

import (
	"encoding/json"
	"errors"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	yaml "gopkg.in/yaml.v2"

	"bytes"

	"github.com/Peripli/service-manager-cli/internal/cmd"
	"github.com/Peripli/service-manager-cli/pkg/smclient/smclientfakes"
	"github.com/Peripli/service-manager-cli/pkg/types"
)

func TestListBrokersCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}

var _ = Describe("List brokers command test", func() {

	var client *smclientfakes.FakeClient
	var command *ListBrokersCmd
	var buffer *bytes.Buffer
	broker := types.Broker{
		Name: "broker1",
		ID:   "id1",
		URL:  "http://broker1.com",
	}
	broker2 := types.Broker{
		Name: "broker2",
		ID:   "id2",
		URL:  "http://broker2.com",
	}

	BeforeEach(func() {
		buffer = &bytes.Buffer{}
		client = &smclientfakes.FakeClient{}
		context := &cmd.Context{Output: buffer, Client: client}
		command = NewListBrokersCmd(context)
	})

	executeWithArgs := func(args []string) error {
		commandToRun := command.Prepare(cmd.SmPrepare)
		commandToRun.SetArgs(args)

		return commandToRun.Execute()
	}

	Context("when no brokers are registered", func() {
		It("should list empty brokers", func() {
			client.ListBrokersReturns(&types.Brokers{Brokers: []types.Broker{}}, nil)
			err := executeWithArgs([]string{})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring("No brokers registered"))
		})
	})

	Context("when brokers are registered", func() {
		It("should list 1 broker", func() {
			result := &types.Brokers{Brokers: []types.Broker{broker}}
			client.ListBrokersReturns(result, nil)
			err := executeWithArgs([]string{})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring(result.TableData().String()))
		})

		It("should list more brokers", func() {
			result := &types.Brokers{Brokers: []types.Broker{broker, broker2}}
			client.ListBrokersReturns(result, nil)
			err := executeWithArgs([]string{})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).To(ContainSubstring(result.Message()))
			Expect(buffer.String()).To(ContainSubstring(result.TableData().String()))
		})
	})

	Context("when format flag is used", func() {
		It("should print in json", func() {
			result := &types.Brokers{Brokers: []types.Broker{broker}}
			client.ListBrokersReturns(result, nil)

			executeWithArgs([]string{"-o", "json"})

			jsonByte, _ := json.MarshalIndent(result, "", "  ")
			jsonOutputExpected := string(jsonByte) + "\n"
			Expect(buffer.String()).To(ContainSubstring(jsonOutputExpected))
		})

		It("should print in yaml", func() {
			result := &types.Brokers{Brokers: []types.Broker{broker}}
			client.ListBrokersReturns(result, nil)

			executeWithArgs([]string{"-o", "yaml"})

			yamlByte, _ := yaml.Marshal(result)
			yamlOutputExpected := string(yamlByte) + "\n"
			Expect(buffer.String()).To(ContainSubstring(yamlOutputExpected))
		})
	})

	Context("when invalid flag is used", func() {
		It("should handle cobra error", func() {
			err := executeWithArgs([]string{"--ooutput", "json"})
			Expect(err).Should(HaveOccurred())
			Expect(err).To(MatchError("unknown flag: --ooutput"))
		})

		It("should handle wrong value", func() {
			err := executeWithArgs([]string{"--output", "invalid"})

			Expect(err).Should(HaveOccurred())
			Expect(err).To(MatchError("unknown output: invalid"))
		})
	})

	Context("when error is returned by Service manager", func() {
		It("should handle error", func() {
			client.ListBrokersReturns(nil, errors.New("Http Client Error"))
			err := executeWithArgs([]string{})

			Expect(err).Should(HaveOccurred())
			Expect(err).To(MatchError("Http Client Error"))
		})
	})

})
