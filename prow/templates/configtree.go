/*
Copyright 2019 The Kubernetes Authors.

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

package templates

// T is the interface holding configuration fields.
type T interface {
	// Apply returns a config that merges the child into the parent
	Apply(T) T
}

// ConfigTree specifies the global generic config for a plugin.
type ConfigTree struct {
	T
	Orgs map[string]Org `json:"orgs,omitempty"`
}

// Org holds the default config for an entire org, as well as any repo overrides.
type Org struct {
	T
	Repos map[string]Repo `json:"repos,omitempty"`
}

// Repo holds the default config for all branches in a repo, as well as specific branch overrides.
type Repo struct {
	T
	Branches map[string]T `json:"branches,omitempty"`
}

// GetOrg returns the org config after merging in any global config.
func (t ConfigTree) GetOrg(name string) Org {
	o, ok := t.Orgs[name]
	if ok {
		o.T = t.Apply(o.T)
	} else {
		o.T = t.T
	}
	return o
}

// GetRepo returns the repo config after merging in any org config.
func (o Org) GetRepo(name string) Repo {
	r, ok := o.Repos[name]
	if ok {
		r.T = o.Apply(r.T)
	} else {
		r.T = o.T
	}
	return r
}

// GetBranch returns the branch config after merging in any repo config.
func (r Repo) GetBranch(name string) T {
	b, ok := r.Branches[name]
	if ok {
		b = r.Apply(b)
	} else {
		b = r.T
	}
	return b
}

// BranchOptions returns the plugin configuration for a given org/repo/branch.
func (t *ConfigTree) BranchOptions(org, repo, branch string) T {
	return t.GetOrg(org).GetRepo(repo).GetBranch(branch)
}

// RepoOptions returns the plugin configuration for a given org/repo.
func (t *ConfigTree) RepoOptions(org, repo string) T {
	return t.GetOrg(org).GetRepo(repo).T
}

// OrgOptions returns the plugin configuration for a given org.
func (t *ConfigTree) OrgOptions(org string) T {
	return t.GetOrg(org).T
}
