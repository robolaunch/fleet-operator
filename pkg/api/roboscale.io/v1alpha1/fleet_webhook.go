/*
Copyright 2022.

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

package v1alpha1

import (
	"errors"

	"github.com/robolaunch/fleet-operator/internal"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var fleetlog = logf.Log.WithName("fleet-resource")

func (r *Fleet) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-fleet-roboscale-io-v1alpha1-fleet,mutating=true,failurePolicy=fail,sideEffects=None,groups=fleet.roboscale.io,resources=fleets,verbs=create;update,versions=v1alpha1,name=mfleet.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Fleet{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Fleet) Default() {
	fleetlog.Info("default", "name", r.Name)
}

//+kubebuilder:webhook:path=/validate-fleet-roboscale-io-v1alpha1-fleet,mutating=false,failurePolicy=fail,sideEffects=None,groups=fleet.roboscale.io,resources=fleets,verbs=create;update,versions=v1alpha1,name=vfleet.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Fleet{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Fleet) ValidateCreate() error {
	fleetlog.Info("validate create", "name", r.Name)

	err := r.checkTenancyLabels()
	if err != nil {
		return err
	}

	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Fleet) ValidateUpdate(old runtime.Object) error {
	fleetlog.Info("validate update", "name", r.Name)

	err := r.checkTenancyLabels()
	if err != nil {
		return err
	}

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Fleet) ValidateDelete() error {
	fleetlog.Info("validate delete", "name", r.Name)
	return nil
}

func (r *Fleet) checkTenancyLabels() error {
	labels := r.GetLabels()

	if _, ok := labels[internal.ORGANIZATION_LABEL_KEY]; !ok {
		return errors.New("organization label should be added with key " + internal.ORGANIZATION_LABEL_KEY)
	}

	if _, ok := labels[internal.TEAM_LABEL_KEY]; !ok {
		return errors.New("team label should be added with key " + internal.TEAM_LABEL_KEY)
	}

	if _, ok := labels[internal.REGION_LABEL_KEY]; !ok {
		return errors.New("super cluster label should be added with key " + internal.REGION_LABEL_KEY)
	}

	if _, ok := labels[internal.CLOUD_INSTANCE_LABEL_KEY]; !ok {
		return errors.New("cloud instance label should be added with key " + internal.CLOUD_INSTANCE_LABEL_KEY)
	}

	if _, ok := labels[internal.CLOUD_INSTANCE_ALIAS_LABEL_KEY]; !ok {
		return errors.New("cloud instance alias label should be added with key " + internal.CLOUD_INSTANCE_ALIAS_LABEL_KEY)
	}

	return nil
}
