// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package targetallocator

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/open-telemetry/opentelemetry-operator/internal/manifests"
	"github.com/open-telemetry/opentelemetry-operator/internal/naming"
)

// Deployment builds the deployment for the given instance.
func Deployment(params manifests.Params) *appsv1.Deployment {
	otelcol := params.OtelCol
	logger := params.Log

	name := naming.TargetAllocator(otelcol.Name)
	labels := Labels(otelcol, name)

	configMap, err := ConfigMap(params)
	if err != nil {
		logger.Info("failed to construct target allocator config map for annotations")
		configMap = nil
	}
	annotations := Annotations(otelcol, configMap)

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: otelcol.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: otelcol.Spec.TargetAllocator.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: annotations,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName:        ServiceAccountName(otelcol),
					Containers:                []corev1.Container{Container(params.Config, logger, otelcol)},
					Volumes:                   Volumes(params.Config, otelcol),
					NodeSelector:              otelcol.Spec.TargetAllocator.NodeSelector,
					TopologySpreadConstraints: otelcol.Spec.TargetAllocator.TopologySpreadConstraints,
				},
			},
		},
	}
}
