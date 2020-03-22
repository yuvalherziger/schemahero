package migration

import (
	"context"

	"github.com/pkg/errors"
	"github.com/schemahero/schemahero/pkg/apis/databases/v1alpha3"
	schemasv1alpha3 "github.com/schemahero/schemahero/pkg/apis/schemas/v1alpha3"
	"github.com/schemahero/schemahero/pkg/logger"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileMigration) getInstance(request reconcile.Request) (*schemasv1alpha3.Migration, error) {
	v1alpha3instance := &schemasv1alpha3.Migration{}
	err := r.Get(context.Background(), request.NamespacedName, v1alpha3instance)
	if err != nil {
		return nil, err // don't wrap
	}

	return v1alpha3instance, nil
}

func (r *ReconcileMigration) reconcileInstance(instance *schemasv1alpha3.Migration) (reconcile.Result, error) {
	logger.Debug("reconciling migration",
		zap.String("kind", instance.Kind),
		zap.String("name", instance.Name),
		zap.String("tableName", instance.Spec.TableName))

	if instance.Status.ApprovedAt > 0 && instance.Status.ExecutedAt == 0 {
		// TODO incomplete code
		connectionURI, err := r.readConnectionURI(instance.Name, v1alpha3.ValueOrValueFrom{})
		if err != nil {
			return reconcile.Result{}, errors.Wrap(err, "failed to get connection uri")
		}

		pod, err := getApplyPod(connectionURI, instance.Namespace, nil, nil)
		if err != nil {
			return reconcile.Result{}, errors.Wrap(err, "failed to get apply pod")
		}

		if err := r.Create(context.Background(), pod); err != nil {
			return reconcile.Result{}, errors.Wrap(err, "failed to create apply pod")
		}

		return reconcile.Result{}, nil
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileMigration) reconcilePod(pod *corev1.Pod) (reconcile.Result, error) {
	// podLabels := pod.GetObjectMeta().GetLabels()
	// role, ok := podLabels["schemahero-role"]
	// if !ok {
	// 	return reconcile.Result{}, nil
	// }

	return reconcile.Result{}, nil
}
