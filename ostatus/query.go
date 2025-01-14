// Copyright (C) 2020 Michael J. Fromberger. All Rights Reserved.

// Package ostatus implements queries that operate on statuses (tweets)
// using the Twitter API v1.1.
package ostatus

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/creachadair/jhttp"
	"github.com/nankys/twitter"
	"github.com/nankys/twitter/internal/otypes"
	"github.com/nankys/twitter/types"
)

// Create constructs a status update ("tweet") with the given text.
// This query requires user-context authorization.
//
// API: 1.1/statuses/update.json
func Create(text string, opts *CreateOpts) Query {
	q := Query{
		Request: &jhttp.Request{
			Method:     "1.1/statuses/update.json",
			HTTPMethod: "POST",
			Params: jhttp.Params{
				"status":    []string{text},
				"trim_user": []string{"true"},
			},
		},
	}
	opts.addQueryParams(&q)
	return q
}

func modQuery(path, id string, opts *Options) Query {
	q := Query{
		Request: &jhttp.Request{
			Method:     path + "/" + id + ".json", // N.B. parameter in path
			HTTPMethod: "POST",
			Params:     jhttp.Params{"trim_user": []string{"true"}},
		},
	}
	opts.addQueryParams(&q)
	return q
}

// Delete constructs a query to delete ("destroy") a tweet with the given ID.
// This query requires user-context authorization.
//
// API: 1.1/statuses/destroy/:id.json
func Delete(id string, opts *Options) Query {
	return modQuery("1.1/statuses/destroy", id, opts)
}

// Retweet constructs a query to retweet a tweet with the given ID.
// This query requires user-context authorization.
//
// API: 1.1/statuses/retweet/:id.json
func Retweet(id string, opts *Options) Query {
	return modQuery("1.1/statuses/retweet", id, opts)
}

// Unretweet constructs a query to un-retweet a tweet with the given ID.
// This query requires user-context authorization.
//
// API: 1.1/statuses/unretweet/:id.json
func Unretweet(id string, opts *Options) Query {
	return modQuery("1.1/statuses/unretweet", id, opts)
}

func likeQuery(path, id string, opts *Options) Query {
	q := Query{
		Request: &jhttp.Request{
			Method:     path + ".json",
			HTTPMethod: "POST",
			Params:     jhttp.Params{"id": []string{id}},
		},
	}
	opts.addQueryParams(&q)
	q.Request.Params.Set("include_entities", strconv.FormatBool(q.opts.Entities))
	return q
}

// Like constructs a query to like ("favorite") a tweet with the given ID.
// This query requires user-context authorization.
//
// API: 1.1/favorites/create.json
func Like(id string, opts *Options) Query {
	return likeQuery("1.1/favorites/create", id, opts)
}

// Unlike constructs a query to un-like ("unfavorite") a tweet with the given ID.
// This query requires user-context authorization.
//
// API: 1.1/favorites/destroy.json
func Unlike(id string, opts *Options) Query {
	return likeQuery("1.1/favorites/destroy", id, opts)
}

// Query is a query to post, delete, like, or retweet a status update.
type Query struct {
	*jhttp.Request
	opts types.TweetFields
}

// Invoke posts the update and reports the resulting tweet.
func (o Query) Invoke(ctx context.Context, cli *twitter.Client) (*Reply, error) {
	data, err := cli.CallRaw(ctx, o.Request)
	if err != nil {
		return nil, err
	}
	var rsp otypes.Tweet
	if err := json.Unmarshal(data, &rsp); err != nil {
		return nil, &jhttp.Error{Message: "decoding response body", Err: err}
	}
	return &Reply{
		Data:   data,
		Tweets: []*types.Tweet{rsp.ToTweetV2(o.opts)},
	}, nil
}

// CreateOpts provides parameters for tweet creation. A nil *CreateOpts
// provides zero values for all fields.
type CreateOpts struct {
	// Record the update as a reply to this tweet ID.  This will be ignored
	// unless the update text includes an @mention of the author of that tweet.
	InReplyTo string

	// Ask the server to automatically populate the reply target and mentions.
	AutoPopulateReply bool

	// User IDs to exclude when auto-populating mentions.
	AutoExcludeMentions []string

	// Optional tweet fields to report with a successful update.
	Optional types.TweetFields
}

func (o *CreateOpts) addQueryParams(q *Query) {
	q.Request.Params.Set("tweet_mode", "extended")
	if o != nil {
		if o.InReplyTo != "" {
			q.Request.Params.Set("in_reply_to_status_id", o.InReplyTo)
		}
		if o.AutoPopulateReply {
			q.Request.Params.Set("auto_populate_reply_metadata", "true")
			if len(o.AutoExcludeMentions) != 0 {
				q.Request.Params.Add("exclude_reply_user_ids", o.AutoExcludeMentions...)
			}
		}
		q.opts = o.Optional
	}
	q.Request.SetBodyToParams()
}

// Options provides parameters for tweet modification. A nil *Options provides
// zero values for all fields.
type Options struct {
	Optional types.TweetFields
}

func (o *Options) addQueryParams(q *Query) {
	if o != nil {
		q.opts = o.Optional
	}
}

// A Reply is the response from a Query or TimelineQuery.
type Reply struct {
	Data   []byte
	Tweets []*types.Tweet
}

func makeTLQuery(id, method string, opts *TimelineOpts) TimelineQuery {
	q := TimelineQuery{
		Request: &jhttp.Request{
			Method: "1.1/statuses/" + method + "_timeline.json",
			Params: make(jhttp.Params),
		},
	}
	opts.addQueryParams(id, &q)
	return q
}

// UserTimeline constructs a query for the specified user's timeline.
//
// API: 1.1/statuses/user_timeline.json
func UserTimeline(id string, opts *TimelineOpts) TimelineQuery {
	return makeTLQuery(id, "user", opts)
}

// HomeTimeline constructs a query for the specified user's home timeline.
// This request requres user-context authorization.
//
// API: 1.1/statuses/home_timeline.json
func HomeTimeline(id string, opts *TimelineOpts) TimelineQuery {
	return makeTLQuery(id, "home", opts)
}

// MentionsTimeline constructs a query for the user's mentions timeline.
// This request requires user-context authorization.
//
// API: 1.1/statuses/mentions_timeline.json
func MentionsTimeline(id string, opts *TimelineOpts) TimelineQuery {
	return makeTLQuery(id, "mentions", opts)
}

// TimelineOpts provides parameters for timeline queries. A nil *TimelineOpts
// provides zero values for all fields.
type TimelineOpts struct {
	// Look up the user by ID instead of by username.
	ByID bool

	// The maximum number of results to return (limit 200).
	// If zero, use the server default.
	MaxResults int

	// Exclude replies from the reported timeline.
	ExcludeReplies bool

	// Include native retweets in the reported timeline.
	IncludeRetweets bool

	// Include entities in the response.
	IncludeEntities bool

	// If set, return results with IDs greater than this (exclusive).
	SinceID string

	// If set, return results with IDs smaller than this (exclusive).
	UntilID string

	// Optional tweet fields to report with the result.
	Optional types.TweetFields
}

func (o *TimelineOpts) keyField() string {
	if o != nil && o.ByID {
		return "user_id"
	}
	return "screen_name"
}

func (o *TimelineOpts) addQueryParams(key string, q *TimelineQuery) {
	q.Request.Params.Set(o.keyField(), key)
	q.Request.Params.Set("trim_user", "true")
	q.Request.Params.Set("tweet_mode", "extended")
	if o == nil {
		return
	}
	q.opts = o.Optional
	if o.MaxResults > 0 {
		q.Request.Params.Set("count", strconv.Itoa(o.MaxResults))
	}
	if o.ExcludeReplies {
		q.Request.Params.Set("exclude_replies", "true")
	}
	if o.IncludeRetweets {
		q.Request.Params.Set("include_rts", "true")
	}
	if o.IncludeEntities {
		q.Request.Params.Set("include_entities", "true")
		q.opts.Entities = true
	}
	if o.SinceID != "" {
		q.Request.Params.Set("since_id", o.SinceID)
	}
	if o.UntilID != "" {
		q.Request.Params.Set("max_id", o.UntilID)
	}
}

// TimelineQuery is a query to fetch a timeline of tweets.
type TimelineQuery struct {
	*jhttp.Request
	opts types.TweetFields
}

// Invoke posts the query and reports the matching tweets.
func (o TimelineQuery) Invoke(ctx context.Context, cli *twitter.Client) (*Reply, error) {
	data, err := cli.CallRaw(ctx, o.Request)
	if err != nil {
		return nil, err
	}
	var rsp []*otypes.Tweet
	if err := json.Unmarshal(data, &rsp); err != nil {
		return nil, &jhttp.Error{Message: "decoding response body", Err: err}
	}
	v2s := make([]*types.Tweet, len(rsp))
	for i, t := range rsp {
		v2s[i] = t.ToTweetV2(o.opts)
	}
	return &Reply{
		Data:   data,
		Tweets: v2s,
	}, nil
}
