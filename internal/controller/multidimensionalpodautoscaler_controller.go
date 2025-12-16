/*
Copyright 2025.

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

package controller

import (
	"context"
	"time"

	autoscalingv1alpha1 "github.com/nulldoot2k/mpa-k8s/api/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// MultidimensionalPodAutoscalerReconciler reconciles a MultidimensionalPodAutoscaler object
type MultidimensionalPodAutoscalerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// RBAC
//+kubebuilder:rbac:groups=autoscaling.hacker-mpa.io,resources=multidimensionalpodautoscalers,verbs=get;list;watch;update;patch
//+kubebuilder:rbac:groups=autoscaling.hacker-mpa.io,resources=multidimensionalpodautoscalers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=autoscaling.hacker-mpa.io,resources=multidimensionalpodautoscalers/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;update;patch

func (r *MultidimensionalPodAutoscalerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var mpa autoscalingv1alpha1.MultidimensionalPodAutoscaler
	if err := r.Get(ctx, req.NamespacedName, &mpa); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var deploy appsv1.Deployment
	if err := r.Get(ctx, types.NamespacedName{
		Name:      mpa.Spec.TargetRef.Name,
		Namespace: mpa.Namespace,
	}, &deploy); err != nil {
		log.Error(err, "unable to fetch Deployment")
		return ctrl.Result{}, err
	}

	currentReplicas := int32(1)
	if deploy.Spec.Replicas != nil {
		currentReplicas = *deploy.Spec.Replicas
	}

	desiredReplicas := currentReplicas

	// MVP logic: enforce min/max replicas
	if currentReplicas < mpa.Spec.Horizontal.MinReplicas {
		desiredReplicas = mpa.Spec.Horizontal.MinReplicas
	}

	if currentReplicas > mpa.Spec.Horizontal.MaxReplicas {
		desiredReplicas = mpa.Spec.Horizontal.MaxReplicas
	}

	if desiredReplicas != currentReplicas {
		log.Info("Scaling deployment",
			"deployment", deploy.Name,
			"from", currentReplicas,
			"to", desiredReplicas,
		)

		deploy.Spec.Replicas = pointer.Int32(desiredReplicas)
		if err := r.Update(ctx, &deploy); err != nil {
			return ctrl.Result{}, err
		}

		mpa.Status.CurrentReplicas = desiredReplicas
		mpa.Status.LastScaleTime = metav1.Now()
		_ = r.Status().Update(ctx, &mpa)
	}

	return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
}

func (r *MultidimensionalPodAutoscalerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&autoscalingv1alpha1.MultidimensionalPodAutoscaler{}).
		Complete(r)
}
