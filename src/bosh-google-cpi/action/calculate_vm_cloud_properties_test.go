package action_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	. "bosh-google-cpi/action"

	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("CalculateVMCloudProperties", func() {
	var subject CalculateVMCloudProperties

	BeforeEach(func() {
		subject = NewCalculateVMCloudProperties()
	})

	DescribeTable("valid machine specification return matching specs", func(cpu, ram, disk int) {
		res, err := subject.Run(DesiredVMSpec{CPU: cpu, RAM: ram, EphemeralDiskSize: disk})
		Expect(err).NotTo(HaveOccurred())

		Expect(res).To(MatchFields(IgnoreExtras, Fields{
			"CPU":            Equal(cpu),
			"RAM":            Equal(ram),
			"RootDiskSizeGb": Equal(disk),
		}))
	},
		Entry("1 core, 1024 memory, 1024 disk", 1, 1024, 1024),
		Entry("2 core, 2048 memory, 1024 disk", 2, 2048, 1024),
		Entry("4 core, 6144 memory, 1024 disk", 4, 6144, 1024),
	)
	DescribeTable("invalid machine specification returns error", func(cpu, ram, disk int, errs []string) {
		_, err := subject.Run(DesiredVMSpec{CPU: cpu, RAM: ram, EphemeralDiskSize: disk})
		Expect(err).To(HaveOccurred())

		for _, expected := range errs {
			Expect(err.Error()).To(ContainSubstring(expected))
		}
	},
		Entry("no cores", 0, 1024, 1024, []string{NoCPUErr}),
		Entry("no memory", 1, 0, 1024, []string{RamPerCPUErr}),
		Entry("bad memory multiple", 1, 922, 1024, []string{RamMultipleErr}),
		Entry("too little memory, bad multiple", 2, 1843, 1024, []string{RamPerCPUErr, RamMultipleErr}),
	)
	Context("valid machine specifications", func() {

	})
})
