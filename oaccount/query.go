// Copyright (C) 2022 Michael J. Fromberger. All Rights Reserved.

// Package oaccount implements queries that operate on account (tweets)
// using the Twitter API v1.1.
package oaccount

import (
	"context"
	"github.com/creachadair/jhttp"
	"github.com/nankys/twitter"
	"github.com/nankys/twitter/internal/ocall"
	"github.com/nankys/twitter/types"
)

// Settings constructs a status update ("tweet") with the given text.
// This query requires user-context authorization.

// API: 1.1/account/settings.json
func Settings(text string) Query {
	q := Query{
		Request: &jhttp.Request{
			Method: "1.1/account/settings.json",
			Params: make(jhttp.Params),
		},
	}

	return q
}

// Query is a query for list memberships.
type Query struct {
	*jhttp.Request
	opts types.UserFields
}

// HasMorePages reports whether the query has more pages to fetch.  This is
// true for a freshly-constructed query, and for an invoked query where the
// server not reported a next-page token.
func (q Query) HasMorePages() bool { return ocall.HasMorePages(q.Request) }

// ResetPageToekn resets (clears) the query's current page token.  Subsequently
// invoking the query will then fetch the first page of results.
func (q Query) ResetPageToken() { ocall.ResetPageToken(q.Request) }

// Invoke executes the query and returns the matching users.
func (q Query) Invoke(ctx context.Context, cli *twitter.Client) (*Reply, error) {
	return ocall.GetUsers(ctx, q.Request, q.opts, cli)
}

// A Reply is the response from a Query.
type Reply = ocall.UsersReply

// verify_credentials use method to test
// This query requires user-context authorization.

// API: 1.1/account/verify_credentials.json
func VerifyCredentials() Query {
	q := Query{
		Request: &jhttp.Request{
			Method: "1.1/account/verify_credentials.json",
			Params: make(jhttp.Params),
		},
	}
	q.Request.Params.Set("include_email", "true")
	return q
}
