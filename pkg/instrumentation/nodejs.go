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

package instrumentation

import (
	"fmt"
	"path"

	corev1 "k8s.io/api/core/v1"

	"github.com/open-telemetry/opentelemetry-operator/apis/v1alpha1"
)

const (
	envSprklPrefix = "SPRKL_PREFIX"
	envNodePath    = "NODE_PATH"
	envNodeOptions = "NODE_OPTIONS"

	nodeRequireArgument = "-r @sprkl/obs"
	sprklPrefix         = "/.sprkl"
)

var (
	sprklNodePath         = path.Join(sprklPrefix, "lib", "node_modules")
	nodeVolumeName        = fmt.Sprintf("%s-nodejs", volumeName)
	nodeInitContainerName = fmt.Sprintf("%s-nodejs", initContainerName)
)

func injectNodeJSSDK(nodeJSSpec v1alpha1.NodeJS, pod corev1.Pod, index int) (corev1.Pod, error) {
	// caller checks if there is at least one container.
	container := &pod.Spec.Containers[index]

	err := validateContainerEnv(container.Env, envNodeOptions)
	if err != nil {
		return pod, err
	}

	// inject NodeJS instrumentation spec env vars.
	for _, env := range nodeJSSpec.Env {
		idx := getIndexOfEnv(container.Env, env.Name)
		if idx == -1 {
			container.Env = append(container.Env, env)
		}
	}
	// @eliran TODO: export appending
	idx := getIndexOfEnv(container.Env, envNodeOptions)
	if idx == -1 {
		container.Env = append(container.Env, corev1.EnvVar{
			Name:  envNodeOptions,
			Value: nodeRequireArgument,
		})
	} else if idx > -1 {
		container.Env[idx].Value = fmt.Sprintf("%s %s", container.Env[idx].Value, nodeRequireArgument)
	}

	idx = getIndexOfEnv(container.Env, envNodePath)
	if idx == -1 {
		container.Env = append(container.Env, corev1.EnvVar{
			Name:  envNodePath,
			Value: sprklNodePath,
		})
	} else if idx > -1 {
		container.Env[idx].Value = fmt.Sprintf("%s:%s", container.Env[idx].Value, sprklNodePath)
	}

	container.Env = append(container.Env, corev1.EnvVar{
		Name:  envSprklPrefix,
		Value: sprklPrefix,
	})

	container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{
		Name:      nodeVolumeName,
		MountPath: sprklPrefix,
	})

	// We just inject Volumes and init containers for the first processed container
	if isInitContainerMissing(pod, nodeInitContainerName) {
		pod.Spec.Volumes = append(pod.Spec.Volumes, corev1.Volume{
			Name: nodeVolumeName,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			}})

		pod.Spec.InitContainers = append(pod.Spec.InitContainers, corev1.Container{
			Name:    nodeInitContainerName,
			Image:   nodeJSSpec.Image,
			Command: []string{"cp", "-a", "/root/.sprkl/.", sprklPrefix},
			VolumeMounts: []corev1.VolumeMount{{
				Name:      nodeVolumeName,
				MountPath: sprklPrefix,
			}},
		})
	}
	return pod, nil
}
