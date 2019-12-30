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

package plugins

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"sigs.k8s.io/yaml"
)

var (
	y   = true
	n   = false
	yes = &y
	no  = &n
)

func TestApproveConfigTreeApply(t *testing.T) {
	var cases = []struct {
		name     string
		child    Approve
		expected Approve
		parent   Approve
	}{
		{
			name:     "all empty",
			child:    Approve{},
			expected: Approve{},
			parent:   Approve{},
		},
		{
			name:  "empty child",
			child: Approve{},
			expected: Approve{
				RequireSelfApproval: yes,
				IgnoreReviewState:   yes,
			},
			parent: Approve{
				IssueRequired:       true,
				RequireSelfApproval: yes,
				LgtmActsAsApprove:   true,
				IgnoreReviewState:   yes,
			},
		},
		{
			name: "empty parent",
			child: Approve{
				IssueRequired:       true,
				RequireSelfApproval: yes,
				LgtmActsAsApprove:   true,
				IgnoreReviewState:   yes,
			},
			expected: Approve{
				IssueRequired:       true,
				RequireSelfApproval: yes,
				LgtmActsAsApprove:   true,
				IgnoreReviewState:   yes,
			},
			parent: Approve{},
		},
		{
			name: "all true",
			child: Approve{
				IssueRequired:       true,
				RequireSelfApproval: yes,
				LgtmActsAsApprove:   true,
				IgnoreReviewState:   yes,
			},
			expected: Approve{
				IssueRequired:       true,
				RequireSelfApproval: yes,
				LgtmActsAsApprove:   true,
				IgnoreReviewState:   yes,
			},
			parent: Approve{
				IssueRequired:       true,
				RequireSelfApproval: yes,
				LgtmActsAsApprove:   true,
				IgnoreReviewState:   yes,
			},
		},
		{
			name: "all false",
			child: Approve{
				IssueRequired:       false,
				RequireSelfApproval: no,
				LgtmActsAsApprove:   false,
				IgnoreReviewState:   no,
			},
			expected: Approve{
				RequireSelfApproval: no,
				IgnoreReviewState:   no,
			},
			parent: Approve{
				IssueRequired:       false,
				RequireSelfApproval: no,
				LgtmActsAsApprove:   false,
				IgnoreReviewState:   no,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if diff := cmp.Diff(c.expected, c.parent.Apply(c.child)); diff != "" {
				t.Error("returned config does not match expected for kubernetes\n", diff)
			}
		})
	}
}

func TestApproveConfigTree(t *testing.T) {
	var cases = []struct {
		name                  string
		config                []byte
		expectedApproveBranch Approve
		expectedApproveOrg    Approve
		expectedApproveRepo   Approve
	}{
		{
			name: "approve no orgs",
			config: []byte(`
issue_required: false
require_self_approval: true
lgtm_acts_as_approve: false
ignore_review_state: true
`),
			expectedApproveBranch: Approve{
				IssueRequired:       false,
				RequireSelfApproval: yes,
				LgtmActsAsApprove:   false,
				IgnoreReviewState:   yes,
			},
			expectedApproveOrg: Approve{
				IssueRequired:       false,
				RequireSelfApproval: yes,
				LgtmActsAsApprove:   false,
				IgnoreReviewState:   yes,
			},
			expectedApproveRepo: Approve{
				IssueRequired:       false,
				RequireSelfApproval: yes,
				LgtmActsAsApprove:   false,
				IgnoreReviewState:   yes,
			},
		},
		{
			name: "approve no default",
			config: []byte(`
orgs:
  bazelbuild:
    ignore_review_state: false
  kubernetes:
    lgtm_acts_as_approve: true
    repos:
      kops:
        lgtm_acts_as_approve: false
      kubernetes:
        require_self_approval: true
`),
			expectedApproveBranch: Approve{
				RequireSelfApproval: yes,
			},
			expectedApproveOrg: Approve{
				LgtmActsAsApprove: true,
			},
			expectedApproveRepo: Approve{
				RequireSelfApproval: yes,
			},
		},
		{
			name: "approve full",
			config: []byte(`
issue_required: false
require_self_approval: false
lgtm_acts_as_approve: false
ignore_review_state: true
orgs:
  bazelbuild:
    ignore_review_state: false
  kubernetes:
    lgtm_acts_as_approve: true
    repos:
      kops:
        lgtm_acts_as_approve: false
      kubernetes:
        require_self_approval: true
        branches:
          master:
            require_self_approval: false
`),
			expectedApproveBranch: Approve{
				RequireSelfApproval: no,
				IgnoreReviewState:   yes,
			},
			expectedApproveOrg: Approve{
				IssueRequired:       false,
				RequireSelfApproval: no,
				LgtmActsAsApprove:   true,
				IgnoreReviewState:   yes,
			},
			expectedApproveRepo: Approve{
				RequireSelfApproval: yes,
				IgnoreReviewState:   yes,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var tree ApproveConfigTree
			if err := yaml.Unmarshal(c.config, &tree); err != nil {
				t.Errorf("error unmarshaling config: %v", err)
			}
			if diff := cmp.Diff(c.expectedApproveOrg, tree.OrgOptions("kubernetes")); diff != "" {
				t.Error("returned config does not match expected for kubernetes\n", diff)
			}
			if diff := cmp.Diff(c.expectedApproveRepo, tree.RepoOptions("kubernetes", "kubernetes")); diff != "" {
				t.Error("returned config does not match expected for kubernetes/kubernetes\n", diff)
			}
			if diff := cmp.Diff(c.expectedApproveBranch, tree.BranchOptions("kubernetes", "kubernetes", "master")); diff != "" {
				t.Error("returned config does not match expected for kubernetes/kubernetes:master\n", diff)
			}
		})
	}
}
