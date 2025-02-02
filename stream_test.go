// Copyright (C) 2020 Michael J. Fromberger. All Rights Reserved.

package twitter_test

import (
	"context"
	"testing"

	"github.com/creachadair/jhttp"
	"github.com/nankys/twitter"
	"github.com/nankys/twitter/rules"
	"github.com/nankys/twitter/tweets"
	"github.com/nankys/twitter/types"
)

func TestStream(t *testing.T) {
	if *testMode == "record" && *maxBodyBytes == 0 {
		t.Fatal("Cannot record streaming methods without a -max-body-size")
	}
	ctx := context.Background()

	req := &jhttp.Request{
		Method: "2/tweets/sample/stream",
		Params: jhttp.Params{
			"tweet.fields": []string{"author_id", "entities"},
		},
	}

	const maxResults = 3

	nr := 0
	err := cli.Stream(ctx, req, func(rsp *twitter.Reply) error {
		nr++
		t.Logf("Msg %d: %s", nr, string(rsp.Data))
		if nr == maxResults {
			return jhttp.ErrStopStreaming
		}
		return nil
	})
	if err != nil {
		t.Errorf("Error from Stream: %v", err)
	}
}

func TestSearchStream(t *testing.T) {
	if *testMode == "record" && *maxBodyBytes == 0 {
		t.Fatal("Cannot record streaming methods without a -max-body-size")
	}
	ctx := context.Background()

	if *testMode == "record" {
		t.Log(`
WARNING: This test may take a while (minutes) to complete in record mode.
         Be patient, it is waiting for live data from a search query.
`)
	}
	r := rules.Adds{{Query: `cat has:images lang:en`}}
	rsp, err := rules.Update(r).Invoke(ctx, cli)
	if err != nil {
		t.Fatalf("Updating rules: %v", err)
	}
	id := rsp.Rules[0].ID

	t.Run("Search", func(t *testing.T) {
		err := tweets.SearchStream(func(rsp *tweets.Reply) error {
			for _, tw := range rsp.Tweets {
				t.Logf("Result: id=%s, author=%s, text=%s", tw.ID, tw.AuthorID, tw.Text)
			}
			return nil
		}, &tweets.StreamOpts{
			MaxResults: 3,
			Optional: []types.Fields{
				types.TweetFields{AuthorID: true},
			},
		}).Invoke(ctx, cli)
		if err != nil {
			t.Errorf("SearchStream failed: %v", err)
		}
	})

	del := rules.Deletes{id}
	if _, err := rules.Update(del).Invoke(ctx, cli); err != nil {
		t.Fatalf("Deleting rules: %v", err)
	}
}
