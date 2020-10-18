/*


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
	"errors"
	"fmt"
	"log"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	webappv1 "example3/api/v1"
)

// DeploymentReconciler reconciles a Deployment object
type DeploymentReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=webapp.my.domain,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=webapp.my.domain,resources=deployments/status,verbs=get;update;patch

func (r *DeploymentReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {

	freezedNamspaces := []string{"example-system2"}
	_ = context.Background()
	_ = r.Log.WithValues("deployment", req.NamespacedName)

	// your logic here
	// if req.Namespace == "example-system2" {
	// 	r.Log.WithValues("deployment", req.Namespace)
	// 	return ctrl.Result{}, errors.New("ramama")
	// }.Deployment
	found := &webappv1.Deployment{}
	err := r.Get(context.TODO(), types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, found)
	for _, frozenNamespace := range freezedNamspaces {
		if frozenNamespace == req.Namespace {
			return reconcile.Result{}, errors.New("This is a frozen namepace")
		}
	}

	if err == nil {
		log.Printf("Creating Deployment %s/%s\n", req.Namespace, req.Name)
		log.Printf("This is th deployment %+v\n", found)
		err = r.Create(context.TODO(), found)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Write an event to the ContainerSet instance with the namespace and name of the
		// created deployment
		r.Log.Info("Normal", "Created", fmt.Sprintf("Created deployment %s/%s", req.Namespace, req.Name))

	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Preform update
	// if !reflect.DeepEqual(deploy.Spec, found.Spec) {
	// 	found.Spec = deploy.Spec
	// 	log.Printf("Updating Deployment %s/%s\n", deploy.Namespace, deploy.Name)
	// 	err = r.Update(context.TODO(), found)
	// 	if err != nil {
	// 		return reconcile.Result{}, err
	// 	}

	// 	// Write an event to the ContainerSet instance with the namespace and name of the
	// 	// updated deployment
	// 	r.recorder.Event(instance, "Normal", "Updated", fmt.Sprintf("Updated deployment %s/%s", deploy.Namespace, deploy.Name))

	// }
	return reconcile.Result{}, nil

	//return ctrl.Result{}, nil
}

func (r *DeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1.Deployment{}).
		Complete(r)
}
