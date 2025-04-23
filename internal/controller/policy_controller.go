/*
Copyright 2025 The Ketches Authors.

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
	"slices"
	"time"

	"github.com/ketches/kube-recycle-bin/internal/api"
	"github.com/ketches/kube-recycle-bin/internal/consts"
	"github.com/ketches/kube-recycle-bin/internal/webhook"
	"github.com/ketches/kube-recycle-bin/pkg/tlog"
	"github.com/ketches/kube-recycle-bin/pkg/util"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// RecyclePolicyReconciler reconciles a api.RecyclePolicy object
type RecyclePolicyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *RecyclePolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	tlog.Infof("» reconciling RecyclePolicy [%s]...", req.Name)
	recyclePolicy := &api.RecyclePolicy{}
	if err := r.Get(ctx, req.NamespacedName, recyclePolicy); err != nil {
		if k8serrors.IsNotFound(err) {
			tlog.Infof("» watched RecyclePolicy [%s] deleted, reclaiming webhook...", req.Name)
			if err := r.tryReclaimWebhook(ctx, req.Name); err != nil {
				tlog.Errorf("✗ failed to reclaim webhook: %v", err)
				return ctrl.Result{RequeueAfter: time.Second * 10}, err
			}
			tlog.Infof("✓ webhook reclaimed for RecyclePolicy [%s] done.", req.Name)
			return ctrl.Result{}, nil
		}

		tlog.Errorf("✗ failed to get RecyclePolicy [%s]: %v", req.Name, err)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if err := r.tryBuildWebhook(ctx, recyclePolicy); err != nil {
		tlog.Errorf("✗ failed to build webhook for RecyclePolicy [%s]: %v", req.Name, err)
		return ctrl.Result{}, err
	}

	tlog.Infof("✓ webhook built for RecyclePolicy [%s] done.", req.Name)
	return ctrl.Result{}, nil
}

func (r *RecyclePolicyReconciler) tryReclaimWebhook(ctx context.Context, recyclePolicyName string) error {
	return r.Client.Delete(ctx, &admissionregistrationv1.ValidatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: webhookName(recyclePolicyName),
		},
	})
}

func (r *RecyclePolicyReconciler) tryBuildWebhook(ctx context.Context, recyclePolicy *api.RecyclePolicy) error {
	recyclePolicies := &api.RecyclePolicyList{}
	if err := r.List(ctx, recyclePolicies); err != nil {
		tlog.Errorf("✗ failed to list recycle policies: %v", err)
		return err
	}

	webhook := constructWebhookFromPolicy(recyclePolicy)
	currentWebhook := &admissionregistrationv1.ValidatingWebhookConfiguration{}
	if err := r.Get(ctx, types.NamespacedName{Name: webhook.Name}, currentWebhook); err != nil {
		if k8serrors.IsNotFound(err) {
			tlog.Infof("» creating webhook for RecyclePolicy [%s]...", recyclePolicy.Name)
			if err := r.Client.Create(ctx, webhook); err != nil {
				tlog.Errorf("✗ failed to create webhook for RecyclePolicy [%s]: %v", recyclePolicy.Name, err)
				return err
			}
			tlog.Infof("✓ webhook for RecyclePolicy [%s] created.", recyclePolicy.Name)
			return nil
		}
		return err
	}

	// update the webhook
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		webhook.SetResourceVersion(currentWebhook.ResourceVersion)

		if err := r.Client.Update(ctx, webhook); err != nil {
			// update failed, try to get the latest version, and retry or return error
			if err := r.Client.Get(ctx, types.NamespacedName{Name: webhook.Name}, currentWebhook); err != nil {
				return err
			}
		}
		tlog.Infof("✓ webhook for RecyclePolicy [%s] updated.", recyclePolicy.Name)
		return nil
	})
}

func constructWebhookFromPolicy(recyclePolicy *api.RecyclePolicy) *admissionregistrationv1.ValidatingWebhookConfiguration {
	result := &admissionregistrationv1.ValidatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: webhookName(recyclePolicy.Name),
			Labels: map[string]string{
				"krb.ketches.cn/recycle-policy": recyclePolicy.Name,
			},
		},
		Webhooks: []admissionregistrationv1.ValidatingWebhook{
			{
				AdmissionReviewVersions: []string{"v1"},
				ClientConfig: admissionregistrationv1.WebhookClientConfig{
					CABundle: getCertBytes(),
					Service: &admissionregistrationv1.ServiceReference{
						Name:      consts.WebhookName,
						Namespace: consts.WebhookNamespace,
						Path:      util.Ptr(consts.WebhookServicePath),
					},
				},
				FailurePolicy:  util.Ptr(admissionregistrationv1.Fail),
				MatchPolicy:    util.Ptr(admissionregistrationv1.Exact),
				Name:           consts.WebhookDNSName,
				SideEffects:    util.Ptr(admissionregistrationv1.SideEffectClassNone),
				TimeoutSeconds: util.Ptr(int32(5)),
			},
		},
	}

	result.Webhooks[0].Rules = append(result.Webhooks[0].Rules, admissionregistrationv1.RuleWithOperations{
		Operations: []admissionregistrationv1.OperationType{admissionregistrationv1.Delete},
		Rule: admissionregistrationv1.Rule{
			APIGroups:   []string{recyclePolicy.Target.Group},
			APIVersions: []string{"*"},
			Resources:   []string{recyclePolicy.Target.Resource},
		},
	})

	var namespaceSelector metav1.LabelSelectorRequirement
	if len(recyclePolicy.Target.Namespaces) == 0 || slices.Contains(recyclePolicy.Target.Namespaces, metav1.NamespaceAll) || slices.Contains(recyclePolicy.Target.Namespaces, "*") {
		namespaceSelector = metav1.LabelSelectorRequirement{
			Key:      "kubernetes.io/metadata.name",
			Operator: metav1.LabelSelectorOpExists, // match all namespaces
		}
	} else {
		namespaceSelector = metav1.LabelSelectorRequirement{
			Key:      "kubernetes.io/metadata.name",
			Operator: metav1.LabelSelectorOpIn,
			Values:   recyclePolicy.Target.Namespaces,
		}
	}

	result.Webhooks[0].NamespaceSelector = &metav1.LabelSelector{
		MatchExpressions: []metav1.LabelSelectorRequirement{namespaceSelector},
	}
	return result
}

func webhookName(policyName string) string {
	return consts.WebhookName + "-" + policyName
}

func getCertBytes() []byte {
	cert, _ := webhook.FetchWebhookCertAndKey()
	return cert
}

// SetupWithManager sets up the controller with the Manager.
func (r *RecyclePolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&api.RecyclePolicy{}).
		Complete(r)
}
