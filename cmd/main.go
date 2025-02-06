package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/bluesky-social/jetstream/pkg/client"
	"github.com/bluesky-social/jetstream/pkg/client/schedulers/sequential"
	"github.com/bluesky-social/jetstream/pkg/models"
)

const (
	serverAddr = "wss://jetstream.atproto.tools/subscribe"
)

var lexicons map[string]int64
var eventTypes map[string]int64
var mu sync.Mutex

func main() {
	ctx := context.Background()
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	})))
	logger := slog.Default()

	config := client.DefaultClientConfig()
	config.WebsocketURL = serverAddr
	config.Compress = true

    lexicons = make(map[string]int64)
    eventTypes = make(map[string]int64)

	h := &handler{
		seenSeqs: make(map[int64]struct{}),
	}

	scheduler := sequential.NewScheduler("jetstream_localdev", logger, h.HandleEvent)

	c, err := client.NewClient(config, logger, scheduler)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	cursor := time.Now().Add(5 * -time.Minute).UnixMicro()

	// Every 5 seconds print the events read and bytes read and average event size
    var perflogInterval int64 = 15
	go func() {
		ticker := time.NewTicker(time.Duration(perflogInterval) * time.Second)
        var eventsReadPrev int64 = 0
        var kilobytesReadPrev int64 = 0
		for {
			select {
			case <-ticker.C:
				eventsRead := c.EventsRead.Load()
				kilobytesRead := c.BytesRead.Load()/1000
                eventsDiff := eventsRead - eventsReadPrev
                kilobytesDiff := kilobytesRead - kilobytesReadPrev

                var avgEventSize int64
                if eventsDiff > 0 {
				    avgEventSize = kilobytesDiff / eventsDiff
                } else {
                    avgEventSize = 0 
                }
                eps := eventsDiff / perflogInterval 
                KBps := kilobytesDiff / perflogInterval 
                
                fmt.Println(time.Now().Local(), "| total events in last", perflogInterval, "sec:", eventsDiff, "| total KB in last", perflogInterval, "sec:", kilobytesDiff, "|", eps, "events/s |", KBps, "KB/s |", avgEventSize, "avg event size")
                eventsReadPrev = eventsRead
                kilobytesReadPrev = kilobytesDiff
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(time.Minute)
		for {
			select {
			case <-ticker.C:
                mu.Lock()
                fmt.Println("\n* * * Events: * * *")
                printMapSorted(&eventTypes)
                fmt.Println()
                fmt.Println("* * * Lexicons: * * *")
                printMapSorted(&lexicons)
                fmt.Println()
                mu.Unlock()
			}
		}
	}()

	if err := c.ConnectAndRead(ctx, &cursor); err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	slog.Info("shutdown")
}

type handler struct {
	seenSeqs  map[int64]struct{}
	highwater int64
}

func (h *handler) HandleEvent(ctx context.Context, event *models.Event) error {
    if event.Account != nil {
        eventTypes["account"] = eventTypes["account"]+1
    } else if event.Identity != nil {
        eventTypes["identity"] = eventTypes["identity"]+1
    } else if event.Commit != nil {
        if event.Commit.Operation == models.CommitOperationCreate {
            eventTypes["commit.create"] = eventTypes["commit.create"]+1
        } else if event.Commit.Operation == models.CommitOperationUpdate {
            eventTypes["commit.update"] = eventTypes["commit.update"]+1
        } else if event.Commit.Operation == models.CommitOperationDelete {
            eventTypes["commit.delete"] = eventTypes["commit.delete"]+1
        }
    }
	// Unmarshal the record if there is one
	if event.Commit != nil && (event.Commit.Operation == models.CommitOperationCreate || event.Commit.Operation == models.CommitOperationUpdate) {
        lexicons[event.Commit.Collection] = lexicons[event.Commit.Collection]+1
        /*
		switch event.Commit.Collection {
		case "app.bsky.feed.post":
			var post apibsky.FeedPost
			if err := json.Unmarshal(event.Commit.Record, &post); err != nil {
				return fmt.Errorf("failed to unmarshal post: %w", err)
			}
			fmt.Printf("%v |(%s)| %s\n", time.UnixMicro(event.TimeUS).Local().Format("15:04:05"), event.Did, post.Text)
		}
        */
	}

	return nil
}

func printMapSorted(m *map[string]int64) {
        // Create a slice to hold the keys
        type kv struct {
                Key   string
                Value int64
        }

        var sortedPairs []kv
        for k, v := range *m {
                sortedPairs = append(sortedPairs, kv{k, v})
        }

        // Sort the slice based on the values in descending order
        sort.Slice(sortedPairs, func(i, j int) bool {
                return sortedPairs[i].Value > sortedPairs[j].Value
        })

        // Print the sorted key-value pairs
        for _, pair := range sortedPairs {
                fmt.Printf("%s: %d\n", pair.Key, pair.Value)
        }
}
