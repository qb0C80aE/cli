package v7action_test

import (
	"errors"

	"code.cloudfoundry.org/cli/actor/actionerror"
	. "code.cloudfoundry.org/cli/actor/v7action"
	"code.cloudfoundry.org/cli/actor/v7action/v7actionfakes"
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
	"code.cloudfoundry.org/cli/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Space Quota Actions", func() {
	var (
		actor                     *Actor
		fakeCloudControllerClient *v7actionfakes.FakeCloudControllerClient
		trueValue                 = true
		falseValue                = true
	)

	BeforeEach(func() {
		actor, fakeCloudControllerClient, _, _, _, _ = NewTestActor()
	})

	Describe("CreateSpaceQuota", func() {
		var (
			spaceQuotaName   string
			organizationGuid string
			warnings         Warnings
			executeErr       error
			limits           QuotaLimits
		)

		JustBeforeEach(func() {
			warnings, executeErr = actor.CreateSpaceQuota(spaceQuotaName, organizationGuid, limits)
		})

		When("creating a space quota with all values set", func() {
			BeforeEach(func() {
				fakeCloudControllerClient.CreateSpaceQuotaReturns(ccv3.SpaceQuota{}, ccv3.Warnings{"some-quota-warning"}, nil)
				limits = QuotaLimits{
					TotalMemoryInMB:       &types.NullInt{IsSet: true, Value: 2},
					PerProcessMemoryInMB:  &types.NullInt{IsSet: true, Value: 3},
					TotalInstances:        &types.NullInt{IsSet: true, Value: 4},
					PaidServicesAllowed:   &trueValue,
					TotalServiceInstances: &types.NullInt{IsSet: true, Value: 6},
					TotalRoutes:           &types.NullInt{IsSet: true, Value: 8},
					TotalReservedPorts:    &types.NullInt{IsSet: true, Value: 9},
				}
			})

			It("makes the space quota", func() {
				Expect(executeErr).ToNot(HaveOccurred())

				Expect(fakeCloudControllerClient.CreateSpaceQuotaCallCount()).To(Equal(1))
				givenSpaceQuota := fakeCloudControllerClient.CreateSpaceQuotaArgsForCall(0)

				Expect(givenSpaceQuota).To(Equal(ccv3.SpaceQuota{
					Quota: ccv3.Quota{
						Name: spaceQuotaName,
						Apps: ccv3.AppLimit{
							TotalMemory:       &types.NullInt{IsSet: true, Value: 2},
							InstanceMemory:    &types.NullInt{IsSet: true, Value: 3},
							TotalAppInstances: &types.NullInt{IsSet: true, Value: 4},
						},
						Services: ccv3.ServiceLimit{
							TotalServiceInstances: &types.NullInt{IsSet: true, Value: 6},
							PaidServicePlans:      &trueValue,
						},
						Routes: ccv3.RouteLimit{
							TotalRoutes:        &types.NullInt{IsSet: true, Value: 8},
							TotalReservedPorts: &types.NullInt{IsSet: true, Value: 9},
						},
					},
					OrgGUID:    organizationGuid,
					SpaceGUIDs: nil,
				}))

				Expect(warnings).To(ConsistOf("some-quota-warning"))
			})
		})

		When("creating a space quota with empty limits", func() {
			var (
				ccv3Quota ccv3.SpaceQuota
			)

			BeforeEach(func() {
				spaceQuotaName = "quota-name"
				limits = QuotaLimits{}

				ccv3Quota = ccv3.SpaceQuota{
					Quota: ccv3.Quota{
						Name: spaceQuotaName,
						Apps: ccv3.AppLimit{
							TotalMemory:       &types.NullInt{Value: 0, IsSet: true},
							InstanceMemory:    nil,
							TotalAppInstances: nil,
						},
						Services: ccv3.ServiceLimit{
							TotalServiceInstances: &types.NullInt{Value: 0, IsSet: true},
							PaidServicePlans:      nil,
						},
						Routes: ccv3.RouteLimit{
							TotalRoutes:        &types.NullInt{Value: 0, IsSet: true},
							TotalReservedPorts: &types.NullInt{Value: 0, IsSet: true},
						},
					},
				}
				fakeCloudControllerClient.CreateSpaceQuotaReturns(
					ccv3Quota,
					ccv3.Warnings{"some-quota-warning"},
					nil,
				)
			})

			It("call the create endpoint with the respective values and returns warnings", func() {
				Expect(fakeCloudControllerClient.CreateSpaceQuotaCallCount()).To(Equal(1))

				Expect(warnings).To(ConsistOf("some-quota-warning"))

				passedQuota := fakeCloudControllerClient.CreateSpaceQuotaArgsForCall(0)
				Expect(passedQuota).To(Equal(ccv3Quota))
			})
		})

		When("creating a quota with all values set to unlimited", func() {
			var (
				ccv3Quota ccv3.SpaceQuota
			)

			BeforeEach(func() {
				spaceQuotaName = "quota-name"
				limits = QuotaLimits{
					TotalMemoryInMB:       &types.NullInt{Value: -1, IsSet: true},
					PerProcessMemoryInMB:  &types.NullInt{Value: -1, IsSet: true},
					TotalInstances:        &types.NullInt{Value: -1, IsSet: true},
					PaidServicesAllowed:   &trueValue,
					TotalServiceInstances: &types.NullInt{Value: -1, IsSet: true},
					TotalRoutes:           &types.NullInt{Value: -1, IsSet: true},
					TotalReservedPorts:    &types.NullInt{Value: -1, IsSet: true},
				}
				ccv3Quota = ccv3.SpaceQuota{
					Quota: ccv3.Quota{
						Name: spaceQuotaName,
						Apps: ccv3.AppLimit{
							TotalMemory:       &types.NullInt{Value: 0, IsSet: false},
							InstanceMemory:    &types.NullInt{Value: 0, IsSet: false},
							TotalAppInstances: &types.NullInt{Value: 0, IsSet: false},
						},
						Services: ccv3.ServiceLimit{
							TotalServiceInstances: &types.NullInt{Value: 0, IsSet: false},
							PaidServicePlans:      &trueValue,
						},
						Routes: ccv3.RouteLimit{
							TotalRoutes:        &types.NullInt{Value: 0, IsSet: false},
							TotalReservedPorts: &types.NullInt{Value: 0, IsSet: false},
						},
					},
				}
				fakeCloudControllerClient.CreateSpaceQuotaReturns(
					ccv3Quota,
					ccv3.Warnings{"some-quota-warning"},
					nil,
				)
			})

			It("call the create endpoint with the respective values and returns warnings", func() {
				Expect(fakeCloudControllerClient.CreateSpaceQuotaCallCount()).To(Equal(1))

				Expect(warnings).To(ConsistOf("some-quota-warning"))

				passedQuota := fakeCloudControllerClient.CreateSpaceQuotaArgsForCall(0)
				Expect(passedQuota).To(Equal(ccv3Quota))
			})
		})

		When("creating a quota returns an error", func() {
			BeforeEach(func() {
				fakeCloudControllerClient.CreateSpaceQuotaReturns(
					ccv3.SpaceQuota{},
					ccv3.Warnings{"some-quota-warning"},
					errors.New("create-error"),
				)
			})

			It("returns the error and warnings", func() {
				Expect(fakeCloudControllerClient.CreateSpaceQuotaCallCount()).To(Equal(1))

				Expect(warnings).To(ConsistOf("some-quota-warning"))
				Expect(executeErr).To(MatchError("create-error"))
			})
		})
	})

	Describe("GetSpaceQuotaByName", func() {
		var (
			quotaName  string
			orgGUID    string
			quota      SpaceQuota
			warnings   Warnings
			executeErr error
		)

		BeforeEach(func() {
			quotaName = "quota-name"
			orgGUID = "some-org-guid"
		})

		JustBeforeEach(func() {
			quota, warnings, executeErr = actor.GetSpaceQuotaByName(quotaName, orgGUID)
		})

		When("when the API layer call returns an error", func() {
			BeforeEach(func() {
				fakeCloudControllerClient.GetSpaceQuotasReturns(
					[]ccv3.SpaceQuota{},
					ccv3.Warnings{"some-quota-warning"},
					errors.New("list-error"),
				)
			})

			It("returns the error and prints warnings", func() {
				Expect(fakeCloudControllerClient.GetSpaceQuotasCallCount()).To(Equal(1))

				Expect(warnings).To(ConsistOf("some-quota-warning"))
				Expect(executeErr).To(MatchError("list-error"))
				Expect(quota).To(Equal(SpaceQuota{}))
			})
		})

		When("when the space quota could not be found", func() {
			BeforeEach(func() {
				fakeCloudControllerClient.GetSpaceQuotasReturns(
					[]ccv3.SpaceQuota{},
					ccv3.Warnings{"some-quota-warning"},
					nil,
				)
			})

			It("returns the error and prints warnings", func() {
				Expect(fakeCloudControllerClient.GetSpaceQuotasCallCount()).To(Equal(1))

				Expect(warnings).To(ConsistOf("some-quota-warning"))
				Expect(executeErr).To(MatchError(actionerror.SpaceQuotaNotFoundByNameError{Name: quotaName}))
				Expect(quota).To(Equal(SpaceQuota{}))
			})
		})

		When("getting a single quota by name", func() {
			BeforeEach(func() {
				fakeCloudControllerClient.GetSpaceQuotasReturns(
					[]ccv3.SpaceQuota{
						{
							Quota: ccv3.Quota{
								GUID: "quota-guid",
								Name: quotaName,
							},
						},
					},
					ccv3.Warnings{"some-quota-warning"},
					nil,
				)
			})

			It("queries the API and returns the matching space quota", func() {
				Expect(executeErr).ToNot(HaveOccurred())

				Expect(fakeCloudControllerClient.GetSpaceQuotasCallCount()).To(Equal(1))
				query := fakeCloudControllerClient.GetSpaceQuotasArgsForCall(0)
				Expect(query).To(ConsistOf(
					ccv3.Query{Key: ccv3.OrganizationGUIDFilter, Values: []string{orgGUID}},
					ccv3.Query{Key: ccv3.NameFilter, Values: []string{quotaName}},
				))

				Expect(warnings).To(ConsistOf("some-quota-warning"))
				Expect(quota).To(Equal(SpaceQuota{
					Quota: ccv3.Quota{
						GUID: "quota-guid",
						Name: quotaName,
					},
				}))
			})
		})
	})

	Describe("GetSpaceQuotasByOrgGUID", func() {
		var (
			orgGUID    string
			quotas     []SpaceQuota
			warnings   Warnings
			executeErr error
		)

		BeforeEach(func() {
			orgGUID = "org-guid"
		})

		JustBeforeEach(func() {
			quotas, warnings, executeErr = actor.GetSpaceQuotasByOrgGUID(orgGUID)
		})

		When("when the API layer call returns an error", func() {
			BeforeEach(func() {
				fakeCloudControllerClient.GetSpaceQuotasReturns(
					[]ccv3.SpaceQuota{},
					ccv3.Warnings{"some-quota-warning"},
					errors.New("list-error"),
				)
			})

			It("returns the error and prints warnings", func() {
				Expect(fakeCloudControllerClient.GetSpaceQuotasCallCount()).To(Equal(1))

				Expect(warnings).To(ConsistOf("some-quota-warning"))
				Expect(executeErr).To(MatchError("list-error"))
				Expect(quotas).To(Equal([]SpaceQuota{}))
			})
		})

		When("getting all space quotas associated with the same organization", func() {
			BeforeEach(func() {
				fakeCloudControllerClient.GetSpaceQuotasReturns(
					[]ccv3.SpaceQuota{
						{
							Quota: ccv3.Quota{
								GUID: "quota-guid",
								Name: "quota-beluga",
							},
							OrgGUID: orgGUID,
						},
						{
							Quota: ccv3.Quota{
								GUID: "quota-2-guid",
								Name: "quota-manatee",
							},
							OrgGUID: orgGUID,
						},
					},
					ccv3.Warnings{"some-quota-warning"},
					nil,
				)
			})

			It("queries the API and returns the matching space quotas", func() {
				Expect(executeErr).ToNot(HaveOccurred())

				Expect(fakeCloudControllerClient.GetSpaceQuotasCallCount()).To(Equal(1))
				query := fakeCloudControllerClient.GetSpaceQuotasArgsForCall(0)
				Expect(query).To(ConsistOf(
					ccv3.Query{Key: ccv3.OrganizationGUIDFilter, Values: []string{orgGUID}},
				))

				Expect(warnings).To(ConsistOf("some-quota-warning"))
				Expect(quotas).To(ConsistOf(
					SpaceQuota{
						Quota: ccv3.Quota{
							GUID: "quota-guid",
							Name: "quota-beluga",
						},
						OrgGUID: orgGUID,
					},
					SpaceQuota{
						Quota: ccv3.Quota{
							GUID: "quota-2-guid",
							Name: "quota-manatee",
						},
						OrgGUID: orgGUID,
					},
				))
			})
		})
	})

	Describe("UpdateSpaceQuota", func() {
		var (
			oldQuotaName string
			orgGUID      string
			newQuotaName string
			quotaLimits  QuotaLimits
			warnings     Warnings
			executeErr   error
		)

		BeforeEach(func() {
			oldQuotaName = "old-quota-name"
			orgGUID = "some-org-guid"
			newQuotaName = "new-quota-name"

			quotaLimits = QuotaLimits{
				TotalMemoryInMB:       &types.NullInt{Value: 2048, IsSet: true},
				PerProcessMemoryInMB:  &types.NullInt{Value: 1024, IsSet: true},
				TotalInstances:        &types.NullInt{Value: 0, IsSet: false},
				TotalServiceInstances: &types.NullInt{Value: 0, IsSet: true},
				PaidServicesAllowed:   &trueValue,
				TotalRoutes:           &types.NullInt{Value: 6, IsSet: true},
				TotalReservedPorts:    &types.NullInt{Value: 5, IsSet: true},
			}

			fakeCloudControllerClient.GetSpaceQuotasReturns(
				[]ccv3.SpaceQuota{{Quota: ccv3.Quota{Name: oldQuotaName}}},
				ccv3.Warnings{"get-quotas-warning"},
				nil,
			)
		})

		JustBeforeEach(func() {
			warnings, executeErr = actor.UpdateSpaceQuota(oldQuotaName, orgGUID, newQuotaName, quotaLimits)
		})

		When("the update-quota endpoint returns an error", func() {
			BeforeEach(func() {
				fakeCloudControllerClient.UpdateSpaceQuotaReturns(
					ccv3.SpaceQuota{},
					ccv3.Warnings{"update-quota-warning"},
					errors.New("update-error"),
				)
			})

			It("returns the error and warnings", func() {
				Expect(fakeCloudControllerClient.UpdateSpaceQuotaCallCount()).To(Equal(1))

				Expect(warnings).To(ConsistOf("get-quotas-warning", "update-quota-warning"))
				Expect(executeErr).To(MatchError("update-error"))
			})
		})

		When("no quota limits are being updated", func() {
			var (
				ccv3Quota ccv3.SpaceQuota
			)

			BeforeEach(func() {
				quotaLimits = QuotaLimits{}

				ccv3Quota = ccv3.SpaceQuota{
					Quota: ccv3.Quota{
						Name: oldQuotaName,
						Apps: ccv3.AppLimit{
							TotalMemory:       nil,
							InstanceMemory:    nil,
							TotalAppInstances: nil,
						},
						Services: ccv3.ServiceLimit{
							TotalServiceInstances: nil,
							PaidServicePlans:      nil,
						},
						Routes: ccv3.RouteLimit{
							TotalRoutes:        nil,
							TotalReservedPorts: nil,
						},
					},
				}

				fakeCloudControllerClient.UpdateSpaceQuotaReturns(
					ccv3Quota,
					ccv3.Warnings{"update-quota-warning"},
					nil,
				)
			})

			It("calls the update endpoint with the respective values and returns warnings", func() {
				Expect(fakeCloudControllerClient.UpdateSpaceQuotaCallCount()).To(Equal(1))

				Expect(warnings).To(ConsistOf("get-quotas-warning", "update-quota-warning"))

				passedQuota := fakeCloudControllerClient.UpdateSpaceQuotaArgsForCall(0)

				updatedQuota := ccv3Quota
				updatedQuota.Name = newQuotaName

				Expect(passedQuota).To(Equal(updatedQuota))
			})
		})

		When("the update space quota has all values set to unlimited", func() {
			var (
				ccv3Quota ccv3.SpaceQuota
			)

			BeforeEach(func() {
				quotaLimits = QuotaLimits{
					TotalMemoryInMB:       &types.NullInt{Value: -1, IsSet: true},
					PerProcessMemoryInMB:  &types.NullInt{Value: -1, IsSet: true},
					TotalInstances:        &types.NullInt{Value: -1, IsSet: true},
					PaidServicesAllowed:   &falseValue,
					TotalServiceInstances: &types.NullInt{Value: -1, IsSet: true},
					TotalRoutes:           &types.NullInt{Value: -1, IsSet: true},
					TotalReservedPorts:    &types.NullInt{Value: -1, IsSet: true},
				}

				ccv3Quota = ccv3.SpaceQuota{
					Quota: ccv3.Quota{
						Name: oldQuotaName,
						Apps: ccv3.AppLimit{
							TotalMemory:       &types.NullInt{Value: 0, IsSet: false},
							InstanceMemory:    &types.NullInt{Value: 0, IsSet: false},
							TotalAppInstances: &types.NullInt{Value: 0, IsSet: false},
						},
						Services: ccv3.ServiceLimit{
							TotalServiceInstances: &types.NullInt{Value: 0, IsSet: false},
							PaidServicePlans:      &falseValue,
						},
						Routes: ccv3.RouteLimit{
							TotalRoutes:        &types.NullInt{Value: 0, IsSet: false},
							TotalReservedPorts: &types.NullInt{Value: 0, IsSet: false},
						},
					},
				}

				fakeCloudControllerClient.UpdateSpaceQuotaReturns(
					ccv3Quota,
					ccv3.Warnings{"update-quota-warning"},
					nil,
				)
			})

			It("calls the update endpoint with the respective values and returns warnings", func() {
				Expect(fakeCloudControllerClient.UpdateSpaceQuotaCallCount()).To(Equal(1))

				Expect(warnings).To(ConsistOf("get-quotas-warning", "update-quota-warning"))

				passedQuota := fakeCloudControllerClient.UpdateSpaceQuotaArgsForCall(0)

				updatedQuota := ccv3Quota
				updatedQuota.Name = newQuotaName

				Expect(passedQuota).To(Equal(updatedQuota))
			})
		})

		When("The update space quota endpoint succeeds", func() {
			var (
				ccv3Quota ccv3.SpaceQuota
			)

			BeforeEach(func() {
				ccv3Quota = ccv3.SpaceQuota{
					Quota: ccv3.Quota{
						Name: oldQuotaName,
						Apps: ccv3.AppLimit{
							TotalMemory:       &types.NullInt{Value: 2048, IsSet: true},
							InstanceMemory:    &types.NullInt{Value: 1024, IsSet: true},
							TotalAppInstances: &types.NullInt{Value: 0, IsSet: false},
						},
						Services: ccv3.ServiceLimit{
							TotalServiceInstances: &types.NullInt{Value: 0, IsSet: true},
							PaidServicePlans:      &trueValue,
						},
						Routes: ccv3.RouteLimit{
							TotalRoutes:        &types.NullInt{Value: 6, IsSet: true},
							TotalReservedPorts: &types.NullInt{Value: 5, IsSet: true},
						},
					},
				}

				fakeCloudControllerClient.UpdateSpaceQuotaReturns(
					ccv3Quota,
					ccv3.Warnings{"update-quota-warning"},
					nil,
				)
			})

			It("calls the update endpoint with the respective values and returns warnings", func() {
				Expect(fakeCloudControllerClient.UpdateSpaceQuotaCallCount()).To(Equal(1))

				Expect(warnings).To(ConsistOf("get-quotas-warning", "update-quota-warning"))

				passedQuota := fakeCloudControllerClient.UpdateSpaceQuotaArgsForCall(0)

				updatedQuota := ccv3Quota
				updatedQuota.Name = newQuotaName

				Expect(passedQuota).To(Equal(updatedQuota))
			})
		})

		When("the space quota name is not being updated", func() {
			var (
				ccv3Quota ccv3.SpaceQuota
			)

			BeforeEach(func() {
				newQuotaName = ""

				ccv3Quota = ccv3.SpaceQuota{
					Quota: ccv3.Quota{
						Name: oldQuotaName,
						Apps: ccv3.AppLimit{
							TotalMemory:       &types.NullInt{Value: 2048, IsSet: true},
							InstanceMemory:    &types.NullInt{Value: 1024, IsSet: true},
							TotalAppInstances: &types.NullInt{Value: 0, IsSet: false},
						},
						Services: ccv3.ServiceLimit{
							TotalServiceInstances: &types.NullInt{Value: 0, IsSet: true},
							PaidServicePlans:      &trueValue,
						},
						Routes: ccv3.RouteLimit{
							TotalRoutes:        &types.NullInt{Value: 6, IsSet: true},
							TotalReservedPorts: &types.NullInt{Value: 5, IsSet: true},
						},
					},
				}

				fakeCloudControllerClient.UpdateSpaceQuotaReturns(
					ccv3Quota,
					ccv3.Warnings{"update-quota-warning"},
					nil,
				)
			})
			It("uses the current space quota name in the API request", func() {
				Expect(executeErr).NotTo(HaveOccurred())
				inputQuota := fakeCloudControllerClient.UpdateSpaceQuotaArgsForCall(0)
				Expect(inputQuota.Name).To(Equal("old-quota-name"))
			})
		})
	})
})