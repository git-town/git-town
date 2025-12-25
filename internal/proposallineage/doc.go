// Package proposallineage implements the logic to embed
// a tree of branch lineage and associated proposals
// into the body of proposals.
package proposallineage

// Ideas for polishing this package:
//
// 1. Separate the proposal cache from the Tree.
//    These two data structures don't seem to belong together.
//    Instead, provide a "loadProposal" function that uses this cache internally,
//    and if the cache doesn't have it, loads via a connector.
//    Given that we have cached connectors now,
//    and need those anyways to reuse already loaded proposal data,
//    maybe we can just use the findProposal method from the CachedConnector instance here.
//
// 2. Extract a "render" function from Builder.
//    It takes a Tree instance and returns the rendered Markdown.
//
// 3. See if we can get rid of all these "rebuild" methods.
//    They don't seem to match the domain and introduce a ton of mutability.
