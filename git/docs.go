// Package git implements interface for gathering git history diff data.
//
// Use constructors to initialize it (InitPR, InitLast)
package git

// GITLAB

// source branch (current), target branch (master, dev, etc)
// if: $CI_PIPELINE_SOURCE == 'merge_request_event'
// only:
//    - merge_requests
// CI --> CI_MERGE_REQUEST_SOURCE_BRANCH_NAME -- CI_MERGE_REQUEST_TARGET_BRANCH_NAME

// GITHUB

// merge branch/source branch (current), target branch (master, dev, etc)
// pull_request, pull_request_target
// CI --> GITHUB_REF_NAME/GITHUB_HEAD_REF (example, feature-branch-1) -- GITHUB_BASE_REF (example, main)
