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

// Apply returns a policy that merges the child into the parent
func (parent Approve) Apply(child Approve) Approve {
	new := Approve{
		IssueRequired:                 child.IssueRequired,
		DeprecatedImplicitSelfApprove: selectBool(parent.DeprecatedImplicitSelfApprove, child.DeprecatedImplicitSelfApprove),
		RequireSelfApproval:           selectBool(parent.RequireSelfApproval, child.RequireSelfApproval),
		LgtmActsAsApprove:             child.LgtmActsAsApprove,
		DeprecatedReviewActsAsApprove: selectBool(parent.DeprecatedReviewActsAsApprove, child.DeprecatedReviewActsAsApprove),
		IgnoreReviewState:             selectBool(parent.IgnoreReviewState, child.IgnoreReviewState),
	}
	return new
}

// selectBool returns the child argument if set, otherwise the parent
func selectBool(parent, child *bool) *bool {
	if child != nil {
		return child
	}
	return parent
}
