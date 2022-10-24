// Copyright (C) 2020 Michael J. Fromberger. All Rights Reserved.

package types

// Fields defines a set of optional response fields to request.  This interface
// is satisfied by the generated enumeration types.
type Fields interface {
	// Return the parameter label for this field type.
	Label() string

	// Return the values selected for this field type.
	Values() []string
}

// Expansions represents a set of object field expansions.
type Expansions struct {
	// Return a user object representing the Tweet’s author.
	AuthorID bool `json:"author_id"`

	// Return a Tweet object that this Tweet is referencing (either as a
	// Retweet, Quoted Tweet, or reply).
	ReferencedTweetID bool `json:"referenced_tweets.id"`

	// Return a user object representing the Tweet author this requested Tweet
	// is a reply of.
	InReplyTo bool `json:"in_reply_to_user_id"`

	// Return a media object representing the images, videos, GIFs included in
	// the Tweet.
	MediaKeys bool `json:"attachments.media_keys"`

	// Return a poll object containing metadata for the poll included in the Tweet.
	PollID bool `json:"attachments.poll_ids"`

	// Return a place object containing metadata for the location tagged in the Tweet.
	PlaceID bool `json:"geo.place_id"`

	// Return a user object for the user mentioned in the Tweet.
	MentionUsername bool `json:"entities.mentions.username"`

	// Return a user object for the author of the referenced Tweet.
	ReferencedAuthorID bool `json:"referenced_tweets.id.author_id"`

	// Return a Tweet object representing the Tweet pinned to the top of the
	// user’s profile.
	PinnedTweetID bool `json:"pinned_tweet_id"`

	// Return a user object representing a list's owner.
	OwnerID bool `json:"owner_id"`
}

// Constants for the names of various metrics reported in a Metrics map.  The
// comment beside each constant describes its visibility.
//
// See https://developer.twitter.com/en/docs/twitter-api/metrics
const (
	Metric_FollowersCount    = "followers_count"     // public
	Metric_FollowingCount    = "following_count"     // public
	Metric_ImpressionCount   = "impression_count"    // non-public, organic, promoted
	Metric_LikeCount         = "like_count"          // public, organic, promoted
	Metric_ListedCount       = "listed_count"        // public
	Metric_QuoteCount        = "quote_count"         // public
	Metric_ReplyCount        = "reply_count"         // public, organic, promoted
	Metric_RetweetCount      = "retweet_count"       // public, organic, promoted
	Metric_TweetCount        = "tweet_count"         // public
	Metric_URLLinkClicks     = "url_link_clicks"     // non-public, organic, promoted
	Metric_UserProfileClicks = "user_profile_clicks" // non-public, organic, promoted
	Metric_ViewCount         = "view_count"          // public, organic, promoted

	// Video view quartile metrics. Non-public, organic, promoted.
	Metric_Playback0Count   = "playback_0_count"
	Metric_Playback25Count  = "playback_25_count"
	Metric_Playback50Count  = "playback_50_count"
	Metric_Playback75Count  = "playback_75_count"
	Metric_Playback100Count = "playback_100_count"
)

type VerifyCredentials struct {
	IncludeEntities bool `json:"include_entities"`
	SkipStatus      bool `json:"skip_status"`
	IncludeEmail    bool `json:"include_email"`
}
