package ui

// like : いいね
func (t *tweets) like() {
	c := t.getSelectTweet()

	err := shared.api.Like(c.Tweet.ID)
	if err != nil {
		shared.SetErrorStatus("Like", err.Error())
		return
	}

	shared.SetStatus("Liked", createTweetSummary(c))
}

// unLike : いいね解除
func (t *tweets) unLike() {
	c := t.getSelectTweet()

	err := shared.api.UnLike(c.Tweet.ID)
	if err != nil {
		shared.SetErrorStatus("UnLike", c.Tweet.Text)
		return
	}

	shared.SetStatus("UnLiked", createTweetSummary(c))
}
