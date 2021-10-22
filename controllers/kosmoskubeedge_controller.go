/*
Copyright 2021.

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

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	crdv1 "gitlab.inovex.de/proj-kosmos/infra/setting-up-kosmos-edge/api/v1"
)

// KosmosKubeEdgeReconciler reconciles a KosmosKubeEdge object
type KosmosKubeEdgeReconciler struct {
	//funtionality for interacting with kubernetes api servers
	client.Client
	//Scheme defines methods for serializing and deserializing API objects,
	//a type registry for converting group, version, and kind information to and from Go schemas,
	//and mappings between Go schemas of different versions.
	//A scheme is the foundation for a versioned API and versioned configuration over time.
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=crd.gitlab.inovex.de,resources=kosmoskubeedges,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=crd.gitlab.inovex.de,resources=kosmoskubeedges/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=crd.gitlab.inovex.de,resources=kosmoskubeedges/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the KosmosKubeEdge object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *KosmosKubeEdgeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	//
	log := log.FromContext(ctx)

	//log-Eintrag
	log.Info("I watched a kosmos kube edge")

	//packe type KosmosKubeEdge, gefunden in Go-Controller-Runtime in eine variabel,
	//sodass ich Zugriff auf alle types in der Struktur habe
	//obj
	//edge = struct pointer
	edge := &crdv1.KosmosKubeEdge{}

	// for reconcile interface, retrieves an obj + infos from cluster
	//(request-scoped values, uniquely identifies obj, aus dem struct der crd)
	//schreibt infos in edge
	r.Get(ctx, req.NamespacedName, edge)

	//log-Eintrag
	log.Info("found a signature", "signature", edge.Spec.Signature.Signature)

	for i, definition := range edge.Spec.Body.RequiredTechnicalContainers {
		deployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-%d", edge.Name, i),
				Namespace: edge.Namespace,
			},
		}

		res, err := ctrl.CreateOrUpdate(ctx, r.Client, deployment, func() error {
			deployment.Spec = appsv1.DeploymentSpec{
				Replicas: pointer.Int32Ptr(3),
				Template: v1.PodTemplateSpec{
					Spec: v1.PodSpec{
						Containers: []v1.Container{
							{
								Name:  fmt.Sprintf("requiredTechnicalContainer-%i", i+1),
								Image: fmt.Sprintf("requiredTechnicalContainer-image-%i", i),
							},
						},
					},
				},
			}
			return nil
		})
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to create deployment: %w", err)
		}
		log.Info("result from create or update", "res", res)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KosmosKubeEdgeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&crdv1.KosmosKubeEdge{}).
		Complete(r)
}
