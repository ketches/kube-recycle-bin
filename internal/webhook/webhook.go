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
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/ketches/kube-recycle-bin/internal/api"
	krbclient "github.com/ketches/kube-recycle-bin/internal/client"
	"github.com/ketches/kube-recycle-bin/internal/consts"
	"github.com/ketches/kube-recycle-bin/pkg/tlog"
	admissionv1 "k8s.io/api/admission/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Run starts the webhook server.
func Run() {
	tlog.Info("» starting admission webhook server...")

	ensureTLSFiles()
	http.HandleFunc(consts.WebhookServicePath, recycleDeleteObjects)

	if err := http.ListenAndServeTLS(":443", consts.WebhookServiceTLSCertFile, consts.WebhookServiceTLSKeyFile, nil); err != nil {
		tlog.Fatalf("✗ failed to listen and serve admission webhook: %v", err)
	}
}

func ensureTLSFiles() {
	cert, key := FetchWebhookCertAndKey()

	err := os.WriteFile(consts.WebhookServiceTLSCertFile, cert, 0644)
	if err != nil {
		tlog.Fatalf("✗ failed to write cert file: %v", err)
	}
	err = os.WriteFile(consts.WebhookServiceTLSKeyFile, key, 0644)
	if err != nil {
		tlog.Fatalf("✗ failed to write key file: %v", err)
	}

	tlog.Infof("✓ cert and key files generated.")
}

// recycleDeleteObjects webhook handler for recycling deleted objects.
func recycleDeleteObjects(w http.ResponseWriter, r *http.Request) {
	tlog.Infof("» received request: %s", r.URL.Path)

	review, err := parseRequest(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("✗ failed to parse request: %v", err), http.StatusBadRequest)
		return
	}

	request := review.Request
	tlog.Infof("» prepare to recycle deleted object: %s - %s", request.Kind.String(), types.NamespacedName{
		Name:      request.Name,
		Namespace: request.Namespace,
	}.String())

	// Create RecycleItem to recycle the deleted object.
	recycleItem := api.NewRecycleItem(*request.RequestKind, request.Namespace, request.Name, request.OldObject.Raw)
	if err := retry.OnError(retry.DefaultRetry, k8serrors.IsAlreadyExists, func() error {
		if err := krbclient.RecycleItem().Create(context.Background(), recycleItem, client.CreateOptions{}); err != nil {
			return err
		}

		tlog.Infof("✓ recycle deleted object %s done", types.NamespacedName{
			Name:      request.Name,
			Namespace: request.Namespace,
		}.String())

		return nil
	}); err != nil {
		tlog.Errorf("✗ failed to recycle deleted object: %v", err)
	}

	response(w, review)
}

// parseRequest parses the request of the admission webhook.
func parseRequest(r *http.Request) (*admissionv1.AdmissionReview, error) {
	var (
		request admissionv1.AdmissionReview
	)

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		tlog.Errorf("✗ failed to decode request: %v", err)
		return nil, err
	}

	return &request, nil
}

// response sends the response to the admission webhook.
func response(w http.ResponseWriter, request *admissionv1.AdmissionReview) {
	response := &admissionv1.AdmissionReview{
		TypeMeta: request.TypeMeta,
		Response: &admissionv1.AdmissionResponse{
			UID:     request.Request.UID,
			Allowed: true,
			Result:  nil,
		},
	}

	encodeResponse(w, response)
}

// encodeResponse encodes the response to the admission webhook.
func encodeResponse(w http.ResponseWriter, response *admissionv1.AdmissionReview) {
	if err := json.NewEncoder(w).Encode(response); err != nil {
		tlog.Errorf("✗ failed to encode response: %v", err)
		http.Error(w, fmt.Sprintf("✗ failed to encode response: %v", err), http.StatusInternalServerError)
	}
}
