package controller

import (
	"context"
	"time"

	autoscalingv2 "k8s.io/api/autoscaling/v2"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	autoscalingv1alpha1 "github.com/nulldoot2k/mpa-k8s/api/v1alpha1"
)

type MultidimensionalPodAutoscalerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// RBAC
// +kubebuilder:rbac:groups=autoscaling.hacker-mpa.io,resources=multidimensionalpodautoscalers,verbs=get;list;watch;update;patch
// +kubebuilder:rbac:groups=autoscaling.hacker-mpa.io,resources=multidimensionalpodautoscalers/status,verbs=update;patch
// +kubebuilder:rbac:groups=autoscaling,resources=horizontalpodautoscalers,verbs=*

func (r *MultidimensionalPodAutoscalerReconciler) Reconcile(
	ctx context.Context,
	req ctrl.Request,
) (ctrl.Result, error) {

	mpa := &autoscalingv1alpha1.MultidimensionalPodAutoscaler{}
	if err := r.Get(ctx, req.NamespacedName, mpa); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// ✅ VALIDATION GUARD
	ref := mpa.Spec.ScaleTargetRef
	if ref.APIVersion == "" || ref.Kind == "" || ref.Name == "" {
		mpa.Status.LastAction = "WaitingForValidScaleTargetRef"
		_ = r.Status().Update(ctx, mpa)
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	}

	hpaName := mpa.Name + "-hpa"
	hpa := &autoscalingv2.HorizontalPodAutoscaler{}

	err := r.Get(ctx, types.NamespacedName{
		Name:      hpaName,
		Namespace: mpa.Namespace,
	}, hpa)

	if apierrors.IsNotFound(err) {
		hpa = r.buildHPA(mpa, hpaName)
		if err := ctrl.SetControllerReference(mpa, hpa, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		if err := r.Create(ctx, hpa); err != nil {
			return ctrl.Result{}, err
		}
	} else if err == nil {
		desired := r.buildHPA(mpa, hpaName)
		hpa.Spec = desired.Spec
		if err := r.Update(ctx, hpa); err != nil {
			return ctrl.Result{}, err
		}
	} else {
		return ctrl.Result{}, err
	}

	now := metav1.NewTime(time.Now())
	mpa.Status.LastScaleTime = &now
	mpa.Status.LastAction = "Horizontal"

	_ = r.Status().Update(ctx, mpa)

	return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
}

func (r *MultidimensionalPodAutoscalerReconciler) buildHPA(
	mpa *autoscalingv1alpha1.MultidimensionalPodAutoscaler,
	name string,
) *autoscalingv2.HorizontalPodAutoscaler {

	min := int32(1)
	max := int32(3)

	if mpa.Spec.MinReplicas != nil {
		min = *mpa.Spec.MinReplicas
	}
	if mpa.Spec.MaxReplicas != nil && *mpa.Spec.MaxReplicas >= min {
		max = *mpa.Spec.MaxReplicas
	}

	return &autoscalingv2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: mpa.Namespace,
		},
		Spec: autoscalingv2.HorizontalPodAutoscalerSpec{
			MinReplicas: &min,
			MaxReplicas: max,
			ScaleTargetRef: autoscalingv2.CrossVersionObjectReference{
				APIVersion: mpa.Spec.ScaleTargetRef.APIVersion,
				Kind:       mpa.Spec.ScaleTargetRef.Kind,
				Name:       mpa.Spec.ScaleTargetRef.Name,
			},
		},
	}
}

func (r *MultidimensionalPodAutoscalerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&autoscalingv1alpha1.MultidimensionalPodAutoscaler{}).
		Owns(&autoscalingv2.HorizontalPodAutoscaler{}).
		Complete(r)
}
