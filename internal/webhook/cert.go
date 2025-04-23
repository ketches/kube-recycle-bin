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

package webhook

import (
	"context"

	"github.com/ketches/kube-recycle-bin/internal/consts"
	"github.com/ketches/kube-recycle-bin/pkg/kube"
	"github.com/ketches/kube-recycle-bin/pkg/tlog"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	certuitl "k8s.io/client-go/util/cert"
)

var (
	cert, key []byte
)

func FetchWebhookCertAndKey() ([]byte, []byte) {
	if cert == nil {
		if secret, err := kube.Client().CoreV1().Secrets(consts.WebhookNamespace).Get(context.Background(), consts.WebhookTLSCertSecretName, metav1.GetOptions{}); err != nil {
			if k8serrors.IsNotFound(err) {
				cert, key, err = certuitl.GenerateSelfSignedCertKey(consts.WebhookDNSName, nil, []string{consts.WebhookDNSName})
				if err != nil {
					tlog.Fatalf("✗ failed to generate self-signed cert and key: %v", err)
				}
				secret = &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      consts.WebhookTLSCertSecretName,
						Namespace: consts.WebhookNamespace,
					},
					Data: map[string][]byte{
						corev1.TLSCertKey:       cert,
						corev1.TLSPrivateKeyKey: key,
					},
					Type: corev1.SecretTypeTLS,
				}
				if _, err := kube.Client().CoreV1().Secrets(consts.WebhookNamespace).Create(context.Background(), secret, metav1.CreateOptions{}); err != nil {
					if !k8serrors.IsAlreadyExists(err) {
						tlog.Fatalf("✗ failed to create secret [%s]: %v", consts.WebhookTLSCertSecretName, err)
					}
				}
				tlog.Infof("✓ cert secret [%s] created.", consts.WebhookTLSCertSecretName)
			} else {
				tlog.Fatalf("✗ failed to get secret [%s]: %v", consts.WebhookTLSCertSecretName, err)
			}
		} else {
			tlog.Infof("✓ cert secret [%s] found.", consts.WebhookTLSCertSecretName)
			cert = secret.Data[corev1.TLSCertKey]
			key = secret.Data[corev1.TLSPrivateKeyKey]
		}
	}
	return cert, key
}
