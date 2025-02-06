# bsky-firehose-stats

API that computes and provides various metrics and stats related to the Bluesky firehose (from jetstream)

Keep storage requirements low. Only store aggregates

what stats do we want

- timewindow:
  - current:
  - last hour:
  - last 24h:
  - last month:
  - last year:
  - all time

lexicons/events:
app.bsky.feed:
app.bsky.feed.post
-> reply
-> quote?
app.bsky.feed.postgate
app.bsky.feed.repost
app.bsky.feed.like
app.bsky.feed.generator
app.bsky.feed.threadgate

app.bsky.graph:
app.bsky.graph.block
app.bsky.graph.follow
app.bsky.graph.listblock
app.bsky.graph.listitem ?

app.bsky.profile:
app.bsky.actor.profile

languages:

first step:
over the course of x minutes, run, collect and count all different lexicons occurring
