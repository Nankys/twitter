// Copyright (C) 2022 Michael J. Fromberger. All Rights Reserved.

// Package oaccount implements queries that operate on account (tweets)
// using the Twitter API v1.1.
package oaccount

import (
	"context"
	"encoding/json"
	"github.com/creachadair/jhttp"
	"github.com/nankys/twitter"
	"github.com/nankys/twitter/types"
)

// Query is a query for list memberships.
type Query struct {
	*jhttp.Request
	opts types.VerifyCredentials
}

// Invoke executes the query and returns the matching users.
func (q Query) Invoke(ctx context.Context, cli *twitter.Client) (*VerifyCredentialsResult, error) {
	data, err := cli.CallRaw(ctx, q.Request)
	if err != nil {
		return nil, err
	}
	var rsp *VerifyCredentialsResult
	if err := json.Unmarshal(data, &rsp); err != nil {
		return nil, &jhttp.Error{Message: "decoding response body", Err: err}
	}
	return rsp, nil
}

// API: 1.1/account/verify_credentials.json
func VerifyCredentials(opts *CredentialsOpt) Query {
	q := Query{
		Request: &jhttp.Request{
			Method: "1.1/account/verify_credentials.json",
			Params: make(jhttp.Params),
		},
	}
	opts.addQueryParams(&q)
	return q
}

type CredentialsOpt struct {
	Optional types.VerifyCredentials
}

func (o *CredentialsOpt) addQueryParams(q *Query) {
	q.opts = o.Optional
}

type VerifyCredentialsResult struct {
	ContributorsEnabled            bool        `json:"contributors_enabled"`
	CreatedAt                      string      `json:"created_at"`
	DefaultProfile                 bool        `json:"default_profile"`
	DefaultProfileImage            bool        `json:"default_profile_image"`
	Description                    string      `json:"description"`
	FavouritesCount                int         `json:"favourites_count"`
	FollowRequestSent              interface{} `json:"follow_request_sent"`
	FollowersCount                 int         `json:"followers_count"`
	Following                      interface{} `json:"following"`
	FriendsCount                   int         `json:"friends_count"`
	GeoEnabled                     bool        `json:"geo_enabled"`
	Id                             int         `json:"id"`
	IdStr                          string      `json:"id_str"`
	IsTranslator                   bool        `json:"is_translator"`
	Lang                           string      `json:"lang"`
	ListedCount                    int         `json:"listed_count"`
	Location                       string      `json:"location"`
	Name                           string      `json:"name"`
	Notifications                  interface{} `json:"notifications"`
	ProfileBackgroundColor         string      `json:"profile_background_color"`
	ProfileBackgroundImageUrl      string      `json:"profile_background_image_url"`
	ProfileBackgroundImageUrlHttps string      `json:"profile_background_image_url_https"`
	ProfileBackgroundTile          bool        `json:"profile_background_tile"`
	ProfileImageUrl                string      `json:"profile_image_url"`
	ProfileImageUrlHttps           string      `json:"profile_image_url_https"`
	ProfileLinkColor               string      `json:"profile_link_color"`
	ProfileSidebarBorderColor      string      `json:"profile_sidebar_border_color"`
	ProfileSidebarFillColor        string      `json:"profile_sidebar_fill_color"`
	ProfileTextColor               string      `json:"profile_text_color"`
	ProfileUseBackgroundImage      bool        `json:"profile_use_background_image"`
	Protected                      bool        `json:"protected"`
	ScreenName                     string      `json:"screen_name"`
	ShowAllInlineMedia             bool        `json:"show_all_inline_media"`
	Email                          string      `json:"email"`
	Status                         struct {
		Contributors interface{} `json:"contributors"`
		Coordinates  struct {
			Coordinates []float64 `json:"coordinates"`
			Type        string    `json:"type"`
		} `json:"coordinates"`
		CreatedAt string `json:"created_at"`
		Favorited bool   `json:"favorited"`
		Geo       struct {
			Coordinates []float64 `json:"coordinates"`
			Type        string    `json:"type"`
		} `json:"geo"`
		Id                   int64  `json:"id"`
		IdStr                string `json:"id_str"`
		InReplyToScreenName  string `json:"in_reply_to_screen_name"`
		InReplyToStatusId    int64  `json:"in_reply_to_status_id"`
		InReplyToStatusIdStr string `json:"in_reply_to_status_id_str"`
		InReplyToUserId      int    `json:"in_reply_to_user_id"`
		InReplyToUserIdStr   string `json:"in_reply_to_user_id_str"`
		Place                struct {
			Attributes struct {
			} `json:"attributes"`
			BoundingBox struct {
				Coordinates [][][]float64 `json:"coordinates"`
				Type        string        `json:"type"`
			} `json:"bounding_box"`
			Country     string `json:"country"`
			CountryCode string `json:"country_code"`
			FullName    string `json:"full_name"`
			Id          string `json:"id"`
			Name        string `json:"name"`
			PlaceType   string `json:"place_type"`
			Url         string `json:"url"`
		} `json:"place"`
		RetweetCount int    `json:"retweet_count"`
		Retweeted    bool   `json:"retweeted"`
		Source       string `json:"source"`
		Text         string `json:"text"`
		Truncated    bool   `json:"truncated"`
	} `json:"status"`
	StatusesCount int         `json:"statuses_count"`
	TimeZone      string      `json:"time_zone"`
	Url           interface{} `json:"url"`
	UtcOffset     int         `json:"utc_offset"`
	Verified      bool        `json:"verified"`
}
