package glab_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/forge/glab"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestParseJsonOutput(t *testing.T) {
	t.Parallel()

	t.Run("invalid JSON", func(t *testing.T) {
		t.Parallel()
		give := `[zonk`
		_, err := glab.ParseJSONOutput(give)
		must.Error(t, err)
	})

	t.Run("multiple results", func(t *testing.T) {
		t.Parallel()
		give := `[
		  {"id":394611593,"iid":5,"target_branch":"main","source_branch":"kg-test","project_id":61831152,"title":"foo","state":"opened","imported":false,"imported_from":"none","created_at":"2025-06-27T20:47:24.4Z","updated_at":"2025-06-27T20:47:25.483Z","upvotes":0,"downvotes":0,"author":{"id":22930529,"username":"kev.lar","name":"Kevin Goslar","state":"active","locked":false,"created_at":null,"avatar_url":"https://secure.gravatar.com/avatar/4905ffb1f10f3d653f25087d324657059267357bfd075294e1fda5c6c63c9f5f?s=80\u0026d=identicon","web_url":"https://gitlab.com/kev.lar"},"assignee":null,"assignees":[],"reviewers":[],"source_project_id":61831152,"target_project_id":61831152,"labels":[],"label_details":null,"description":"","draft":false,"milestone":null,"merge_when_pipeline_succeeds":false,"detailed_merge_status":"mergeable","merge_user":null,"merged_at":null,"merge_after":null,"prepared_at":"2025-06-27T20:47:25.477Z","closed_by":null,"closed_at":null,"sha":"b2b4af2c973ed3165f8561e6bc2d6eed43aea0b2","merge_commit_sha":"","squash_commit_sha":"","user_notes_count":0,"should_remove_source_branch":false,"force_remove_source_branch":true,"allow_collaboration":false,"allow_maintainer_to_push":false,"web_url":"https://gitlab.com/git-town-qa/test-repo/-/merge_requests/5","references":{"short":"!5","relative":"!5","full":"git-town-qa/test-repo!5"},"discussion_locked":false,"time_stats":{"human_time_estimate":"","human_total_time_spent":"","time_estimate":0,"total_time_spent":0},"squash":false,"squash_on_merge":false,"task_completion_status":{"count":0,"completed_count":0},"has_conflicts":false,"blocking_discussions_resolved":true,"merged_by":null},
			{"id":328691588,"iid":1,"target_branch":"main","source_branch":"branch-1","project_id":61831152,"title":"Changes in branch 1","state":"opened","imported":false,"imported_from":"none","created_at":"2024-09-19T15:27:48.828Z","updated_at":"2025-05-13T15:37:47.839Z","upvotes":0,"downvotes":0,"author":{"id":22930529,"username":"kev.lar","name":"Kevin Goslar","state":"active","locked":false,"created_at":null,"avatar_url":"https://secure.gravatar.com/avatar/4905ffb1f10f3d653f25087d324657059267357bfd075294e1fda5c6c63c9f5f?s=80\u0026d=identicon","web_url":"https://gitlab.com/kev.lar"},"assignee":null,"assignees":[],"reviewers":[],"source_project_id":61831152,"target_project_id":61831152,"labels":[],"label_details":null,"description":"This is the proposal body text.","draft":false,"milestone":null,"merge_when_pipeline_succeeds":false,"detailed_merge_status":"mergeable","merge_user":null,"merged_at":null,"merge_after":null,"prepared_at":"2024-09-19T15:27:50.592Z","closed_by":null,"closed_at":null,"sha":"5d3b08006a05b1c75879af8f074775686e49447c","merge_commit_sha":"","squash_commit_sha":"","user_notes_count":0,"should_remove_source_branch":false,"force_remove_source_branch":true,"allow_collaboration":false,"allow_maintainer_to_push":false,"web_url":"https://gitlab.com/git-town-qa/test-repo/-/merge_requests/1","references":{"short":"!1","relative":"!1","full":"git-town-qa/test-repo!1"},"discussion_locked":false,"time_stats":{"human_time_estimate":"","human_total_time_spent":"","time_estimate":0,"total_time_spent":0},"squash":false,"squash_on_merge":false,"task_completion_status":{"count":0,"completed_count":0},"has_conflicts":false,"blocking_discussions_resolved":true,"merged_by":null}
		]`
		have, err := glab.ParseJSONOutput(give)
		must.NoError(t, err)
		must.Len(t, 2, have)
		must.EqOp(t, have[0].Data.Data().Number, 5)
		must.EqOp(t, have[1].Data.Data().Number, 1)
	})

	t.Run("no results", func(t *testing.T) {
		t.Parallel()
		give := `[]`
		have, err := glab.ParseJSONOutput(give)
		must.NoError(t, err)
		must.Eq(t, []forgedomain.Proposal{}, have)
	})

	t.Run("single result", func(t *testing.T) {
		t.Parallel()
		give := `[{"id":394611593,"iid":5,"target_branch":"main","source_branch":"kg-test","project_id":61831152,"title":"foo","state":"opened","imported":false,"imported_from":"none","created_at":"2025-06-27T20:47:24.4Z","updated_at":"2025-06-27T20:47:25.483Z","upvotes":0,"downvotes":0,"author":{"id":22930529,"username":"kev.lar","name":"Kevin Goslar","state":"active","locked":false,"created_at":null,"avatar_url":"https://secure.gravatar.com/avatar/4905ffb1f10f3d653f25087d324657059267357bfd075294e1fda5c6c63c9f5f?s=80\u0026d=identicon","web_url":"https://gitlab.com/kev.lar"},"assignee":null,"assignees":[],"reviewers":[],"source_project_id":61831152,"target_project_id":61831152,"labels":[],"label_details":null,"description":"","draft":false,"milestone":null,"merge_when_pipeline_succeeds":false,"detailed_merge_status":"mergeable","merge_user":null,"merged_at":null,"merge_after":null,"prepared_at":"2025-06-27T20:47:25.477Z","closed_by":null,"closed_at":null,"sha":"b2b4af2c973ed3165f8561e6bc2d6eed43aea0b2","merge_commit_sha":"","squash_commit_sha":"","user_notes_count":0,"should_remove_source_branch":false,"force_remove_source_branch":true,"allow_collaboration":false,"allow_maintainer_to_push":false,"web_url":"https://gitlab.com/git-town-qa/test-repo/-/merge_requests/5","references":{"short":"!5","relative":"!5","full":"git-town-qa/test-repo!5"},"discussion_locked":false,"time_stats":{"human_time_estimate":"","human_total_time_spent":"","time_estimate":0,"total_time_spent":0},"squash":false,"squash_on_merge":false,"task_completion_status":{"count":0,"completed_count":0},"has_conflicts":false,"blocking_discussions_resolved":true,"merged_by":null}]`
		have, err := glab.ParseJSONOutput(give)
		must.NoError(t, err)
		want := []forgedomain.Proposal{
			{
				Data: forgedomain.ProposalData{
					Body:         None[gitdomain.ProposalBody](),
					MergeWithAPI: true,
					Number:       5,
					Source:       "kg-test",
					Target:       "main",
					Title:        "foo",
					URL:          "https://gitlab.com/git-town-qa/test-repo/-/merge_requests/5",
				},
				ForgeType: forgedomain.ForgeTypeGitlab,
			},
		}
		must.Eq(t, want, have)
	})
}
