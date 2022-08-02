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
	"strconv"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	crdv1 "github.com/kosmos-industrie40/setting-up-kosmos-edge/api/v1"
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

//+kubebuilder:rbac:groups=crd.github.inovex.de,resources=kosmoskubeedges,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=crd.github.inovex.de,resources=kosmoskubeedges/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=crd.github.inovex.de,resources=kosmoskubeedges/finalizers,verbs=update

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
	edge := &crdv1.KosmosKubeEdge{}

	// retrieves an obj + infos from cluster
	err := r.Get(ctx, req.NamespacedName, edge)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create deployment: %w", err)
	}

	//build deployments for neccessary containers
	for i, machineConnection := range edge.Spec.Body.MachineConnection {
		// var j int
		res, err := r.deployContainers(ctx, edge, machineConnection.Connector, i, 0, "machineconnection")
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to create deployment: %w", err)
		}
		log.Info("result from create or update", "res", res)
	}

	for i, analysisSystems := range edge.Spec.Body.Analysis.Systems {
		// var j int
		res, err := r.deployContainers(ctx, edge, analysisSystems.Connection.Container, i, 0, "analysissystemconnection")
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to create deployment: %w", err)
		}
		log.Info("result from create or update", "res", res)
	}

	for i, blockchainConnection := range edge.Spec.Body.BlockchainConnection.ContainerList {
		// var j int
		res, err := r.deployContainers(ctx, edge, blockchainConnection, i, 0, "blockchainconnection")
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to create deployment: %w", err)
		}
		log.Info("result from create or update", "res", res)
	}

	for i, requTechCont := range edge.Spec.Body.RequiredTechnicalContainers {
		for j, containers := range requTechCont.Containers {
			res, err := r.deployContainers(ctx, edge, containers, i, j, "requtechcont")
			if err != nil {
				return ctrl.Result{}, fmt.Errorf("failed to create deployment: %w", err)
			}
			log.Info("result from create or update", "res", res)
		}
	}

	return ctrl.Result{}, nil
}

func (r *KosmosKubeEdgeReconciler) deployContainers(
	ctx context.Context, edge *crdv1.KosmosKubeEdge, container crdv1.ContainerDefinitions,
	i int, j int, containerName string) (controllerutil.OperationResult, error) {

	//Deployment enables declarative updates for Pods and ReplicaSets
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s-%d-%d", edge.Name, containerName, i+1, j+1),
			Namespace: edge.Namespace,
		},
	}

	// CreateOrUpdate creates or updates the given object obj in the Kubernetes cluster.
	// The object's desired state should be reconciled with the existing state using the passed in ReconcileFn.
	// obj must be a struct pointer so that obj can be updated with the content returned by the Server.
	res, err := ctrl.CreateOrUpdate(ctx, r.Client, deployment, func() error {

		// SetOwnerReference is a helper method to make sure the given object contains an object reference
		// to the object provided. This allows you to declare that owner has a dependency on the object without
		// specifying it as a controller. If a reference to the same object already exists,
		// it'll be overwritten with the newly provided version.
		if err := controllerutil.SetOwnerReference(edge, deployment, r.Scheme); err != nil {
			return err
		}

		// Specification of the desired behavior of the Deployment.
		deployment.Spec = appsv1.DeploymentSpec{
			Replicas: pointer.Int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: edge.Labels,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: edge.Labels,
				},
				Spec: v1.PodSpec{
					NodeSelector: edge.Labels,
					Containers: []v1.Container{
						{
							Name:  "container",
							Image: fmt.Sprintf("%s:%s", container.Url, container.Tag),
							//Args:  containers.Arguments,
							Env:   getEnvVarValues(container.Environment),
							Ports: getListOfPorts(container.Ports),
						},
					},
				},
			},
		}
		return nil
	})
	return res, err
}

func getEnvVarValues(contractEnvironment []string) []v1.EnvVar {
	var res = []v1.EnvVar{}
	for _, envValues := range contractEnvironment {
		values := strings.SplitN(envValues, "=", 2)
		name, value := values[0], values[1]
		res = append(res, v1.EnvVar{
			Name:  name,
			Value: value,
		})
	}
	return res
}

func getListOfPorts(contractPorts []crdv1.PortsContainer) []v1.ContainerPort {
	res := []v1.ContainerPort{}
	for i, port := range contractPorts {
		res = append(res, v1.ContainerPort{
			Name:          alignLabelsToString(port.Label, i+1),
			ContainerPort: int32(port.Src),
		})
	}
	return res
}

func alignLabelsToString(portLabels []string, portNum int) string {
	res := ""
	for _, label := range portLabels {
		if res != "" {
			res = res + "-" + label
			continue
		} else {
			res = res + label
		}
	}

	if res == "" {
		return res
	} else {
		return res + "-" + strconv.Itoa(portNum)
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *KosmosKubeEdgeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&crdv1.KosmosKubeEdge{}).
		Complete(r)
}
