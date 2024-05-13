package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/storage"
	"github.com/mavryk-network/mavryk-snapshot/pkg/snapshot"
	"github.com/mavryk-network/mavryk-snapshot/pkg/store"
	"github.com/patrickmn/go-cache"
)

func getSnapshotResponseCached(ctx context.Context, goCache *cache.Cache, bucketName string) *SnapshotResponse {
	itemsFound, found := goCache.Get("response")
	if found {
		log.Println("Using response from cache...")
		response := (itemsFound.(SnapshotResponse))
		return &response
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	snapshotStorage := store.NewSnapshotStorage(client, bucketName)
	data := snapshotStorage.GetSnapshotItems(ctx)
	response := SnapshotResponse{
		DateGenerated: time.Now().UTC().Format("2006-01-02T15:04:05Z07:00"),
		Org:           "MavrykDynamics",
		Schema:        "https://raw.githubusercontent.com/mavryk-network/mavryk-snapshot-metadata-schema/main/mavryk-snapshot-metadata.schema.json",
		Data:          data,
	}

	goCache.Set("response", response, cache.DefaultExpiration)
	return &response
}

func getNewestSnapshot(
	ctx context.Context,
	goCache *cache.Cache,
	bucketName string,
	historyMode snapshot.HistoryModeType,
	chainName string,
) (*snapshot.SnapshotItem, error) {
	responseCached := getSnapshotResponseCached(ctx, goCache, bucketName)

	for _, item := range responseCached.Data {
		if item.HistoryMode == historyMode && item.ChainName == chainName {
			return &item, nil
		}
	}

	return nil, fmt.Errorf("Snapshot item not found")
}
