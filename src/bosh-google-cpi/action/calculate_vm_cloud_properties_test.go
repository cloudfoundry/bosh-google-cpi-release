package action_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

	. "bosh-google-cpi/action"
)

var _ = Describe("CalculateVMCloudProperties", func() {
	var subject CalculateVMCloudProperties

	BeforeEach(func() {
		subject = NewCalculateVMCloudProperties()
	})

	It("failed if no CPU is specified", func() {
		_, err := subject.Run(DesiredVMSpec{CPU: 0, RAM: 1024, EphemeralDiskSize: 1024})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(NoCPUErr))
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
	DescribeTable("memory is rounded up to the nearest valid value if necessary", func(cpu, ram, roundedRam int) {
		res, err := subject.Run(DesiredVMSpec{CPU: cpu, RAM: ram, EphemeralDiskSize: 1024})
		Expect(err).NotTo(HaveOccurred())

		Expect(res.RAM).To(Equal(roundedRam))
	},
		Entry("1 core, no memory", 1, 0, 1024),
		Entry("1 core, under memory", 1, 922, 1024),
		Entry("2 cores, under memory", 2, 1024, 2048),
		Entry("16 cores, over memory", 16, 16128, 16128),
		Entry("16 cores, under memory", 16, 1, 14848),
	)
})
