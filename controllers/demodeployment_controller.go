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

package controllers

import (
	"context"
	"fmt"
	"reflect"

	core "k8s.io/api/core/v1"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	v1alpha1 "github.com/kklysenko/k8s_test/api/v1alpha1"
)

const REPLICAS_MAXIMUM_NUMBER int32 = 2

// DemoDeploymentReconciler reconciles a DemoDeployment object
type DemoDeploymentReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=k8stest.com,resources=demodeployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=k8stest.com,resources=demodeployments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=k8stest.com,resources=demodeployments/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods/status,verbs=get
//+kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DemoDeployment object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *DemoDeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Enter Reconcile v.0.0.7", "req", req)

	instance := &v1alpha1.DemoDeployment{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		if !apierrors.IsNotFound(err) {
			log.Error(err, "unable to fetch DemoDeployment")
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if instance.Spec.Replicas > REPLICAS_MAXIMUM_NUMBER {
		r.Recorder.Event(instance, "Warning", "Rejected", fmt.Sprintf("Trying to scale for more than 2 pods %s/%s", instance.Namespace, instance.Name))
	}

	deploy := r.deploymentFromDemoDeployment(instance)
	log.Info(fmt.Sprintf("Instance %v\n", instance))

	found := &appsv1.Deployment{}
	err := r.Get(ctx, req.NamespacedName, found)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new Deployment", "Deployment.Namespace", deploy.Namespace, "Deployment.Name", deploy.Name)
		err = r.Create(ctx, deploy)
		if err != nil {
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", deploy.Namespace, "Deployment.Name", deploy.Name)
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return ctrl.Result{}, err
	}

	if err := ctrl.SetControllerReference(instance, deploy, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	if !reflect.DeepEqual(deploy.Spec, found.Spec) {
		if *deploy.Spec.Replicas > REPLICAS_MAXIMUM_NUMBER {
			r.Recorder.Event(instance, "Warning", "Rejected", fmt.Sprintf("Trying to scale for more than 2 pods %s/%s", deploy.Namespace, deploy.Name))
		}
		found.Spec = deploy.Spec
		log.Info("Updating Deployment %s/%s\n", found.Namespace, found.Name)
		err = r.Update(context.TODO(), found)
		if err != nil {
			return reconcile.Result{}, err
		}
	}
	return reconcile.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DemoDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.DemoDeployment{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}

func (r *DemoDeploymentReconciler) deploymentFromDemoDeployment(dd *v1alpha1.DemoDeployment) *appsv1.Deployment {
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dd.Name,
			Namespace: dd.Namespace,
			Labels:    dd.Labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &dd.Spec.Replicas,
			Selector: dd.Spec.Selector,
			Template: core.PodTemplateSpec{
				Spec: dd.Spec.Template.Spec,
				ObjectMeta: metav1.ObjectMeta{
					Labels: dd.Labels,
				},
			},
		},
	}

	return dep
}
