/*
Copyright 2023 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package e2e

import (
	"time"

	. "github.com/onsi/ginkgo/v2"

	"github.com/openshift/cluster-control-plane-machine-set-operator/test/e2e/framework"
	"github.com/openshift/cluster-control-plane-machine-set-operator/test/e2e/helpers"
)

var _ = Describe("ControlPlaneMachineSet Operator", framework.Periodic(), func() {
	BeforeEach(func() {
		helpers.EventuallyClusterOperatorsShouldStabilise(10*time.Minute, 10*time.Second)
	})

	Context("With an active ControlPlaneMachineSet", func() {
		BeforeEach(func() {
			helpers.EnsureActiveControlPlaneMachineSet(testFramework)
		})

		Context("with a cluster wide proxy", Ordered, func() {
			BeforeAll(func() {
				helpers.DeployProxy(testFramework)
				helpers.ConfigureClusterWideProxy(testFramework)
				helpers.EventuallyClusterOperatorsShouldStabilise(45*time.Minute, 10*time.Second)
			})

			AfterAll(func() {
				helpers.UnconfigureClusterWideProxy(testFramework)
				helpers.EventuallyClusterOperatorsShouldStabilise(45*time.Minute, 10*time.Second)
				helpers.DeleteProxy(testFramework)
			})

			Context("and the instance type of index 1 is not as expected", func() {
				BeforeEach(func() {
					helpers.IncreaseControlPlaneMachineInstanceSize(testFramework, 1)
				})
				helpers.ItShouldRollingUpdateReplaceTheOutdatedMachine(testFramework, 1)
			})
		})

		Context("and the instance type is changed", func() {
			BeforeEach(func() {
				helpers.IncreaseControlPlaneMachineSetInstanceSize(testFramework)
			})

			helpers.ItShouldPerformARollingUpdate(&helpers.RollingUpdatePeriodicTestOptions{
				TestFramework: testFramework,
			})
		})

	})
})
