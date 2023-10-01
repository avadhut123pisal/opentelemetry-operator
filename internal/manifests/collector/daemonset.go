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

package collector

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/open-telemetry/opentelemetry-operator/internal/manifests"
	"github.com/open-telemetry/opentelemetry-operator/internal/naming"
)

// DaemonSet builds the deployment for the given instance.
func DaemonSet(params manifests.Params) *appsv1.DaemonSet {
	otelcol := params.OtelCol
	logger := params.Log

	name := naming.Collector(otelcol.Name)
	labels := Labels(otelcol, name, params.Config.LabelsFilter())

	annotations := Annotations(otelcol)
	podAnnotations := PodAnnotations(otelcol)
	return &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:        naming.Collector(otelcol.Name),
			Namespace:   otelcol.Namespace,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: SelectorLabels(otelcol),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: podAnnotations,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: ServiceAccountName(otelcol),
					InitContainers:     otelcol.Spec.InitContainers,
					Containers:         append(otelcol.Spec.AdditionalContainers, Container(params.Config, logger, otelcol, true)),
					Volumes:            Volumes(params.Config, otelcol),
					Tolerations:        otelcol.Spec.Tolerations,
					NodeSelector:       otelcol.Spec.NodeSelector,
					HostNetwork:        otelcol.Spec.HostNetwork,
					DNSPolicy:          getDNSPolicy(otelcol),
					SecurityContext:    otelcol.Spec.PodSecurityContext,
					PriorityClassName:  otelcol.Spec.PriorityClassName,
					Affinity:           otelcol.Spec.Affinity,
				},
			},
		},
	}
}
