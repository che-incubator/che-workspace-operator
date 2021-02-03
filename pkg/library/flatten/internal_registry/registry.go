//
// Copyright (c) 2019-2020 Red Hat, Inc.
// This program and the accompanying materials are made
// available under the terms of the Eclipse Public License 2.0
// which is available at https://www.eclipse.org/legal/epl-2.0/
//
// SPDX-License-Identifier: EPL-2.0
//
// Contributors:
//   Red Hat, Inc. - initial API and implementation
//

package registry

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	devworkspace "github.com/devfile/api/pkg/apis/workspaces/v1alpha2"
	"github.com/devfile/devworkspace-operator/internal/images"
	"sigs.k8s.io/yaml"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	// RegistryDirectory is the directory where the YAML files for the internal registry exist
	RegistryDirectory = "internal-registry/"
)

var log = logf.Log.WithName("registry")

func getPluginPath(pluginID string) string {
	return filepath.Join(RegistryDirectory, pluginID, "devworkspacetemplate.yaml")
}

// IsInInternalRegistry checks if pluginID is in the internal registry
func IsInInternalRegistry(pluginID string) bool {
	if _, err := os.Stat(getPluginPath(pluginID)); err != nil {
		if os.IsNotExist(err) {
			log.Info(fmt.Sprintf("Could not find %s in the internal registry", pluginID))
		}
		return false
	}
	return true
}

func ReadPluginFromInternalRegistry(pluginID string) (*devworkspace.DevWorkspaceTemplate, error) {
	yamlBytes, err := ioutil.ReadFile(getPluginPath(pluginID))
	if err != nil {
		return nil, err
	}
	plugin := &devworkspace.DevWorkspaceTemplate{}
	if err := yaml.Unmarshal(yamlBytes, plugin); err != nil {
		return nil, err
	}
	resolvedPlugin, err := images.FillPluginEnvVars(plugin)
	if err != nil {
		return nil, err
	}
	return resolvedPlugin, nil
}
