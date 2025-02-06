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

## Observed in 30s

Events:
commit.create: 146702
commit.delete: 5488
commit.update: 461
account: 311
identity: 262

Lexicons:
app.bsky.feed.like: 92032
app.bsky.graph.follow: 22612
app.bsky.feed.post: 17092
app.bsky.feed.repost: 13114
app.bsky.graph.block: 1221
app.bsky.actor.profile: 655
app.bsky.graph.listitem: 273
app.bsky.feed.threadgate: 85
app.bsky.feed.postgate: 46
app.bsky.graph.list: 9
app.bsky.graph.listblock: 6
app.bsky.feed.generator: 6
app.bsky.graph.starterpack: 5
chat.bsky.actor.declaration: 5
uk.skyblur.post: 1
jp.5leaf.sync.mastodon: 1
