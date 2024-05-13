package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"github.com/mavryk-network/mavryk-snapshot/pkg/snapshot"
	"github.com/mavryk-network/mavryk-snapshot/pkg/store"
	"github.com/mavryk-network/mavryk-snapshot/pkg/util"
	"github.com/samber/lo"
)

func main() {
	godotenv.Load()
	cron := util.GetEnvString("CRON_EXPRESSION", "0 0 * * *")

	task()

	log.Println("Waiting for the snapshot job...")
	s := gocron.NewScheduler(time.UTC)
	s.Cron(cron).Do(task)
	s.StartBlocking()
}

func task() {
	log.Println("Starting the snapshot job...")
	start := time.Now()
	ctx := context.Background()
	bucketName := os.Getenv("BUCKET_NAME")
	maxDays := util.GetEnvInt("MAX_DAYS", 7)
	maxMonths := util.GetEnvInt("MAX_MONTHS", 6)
	network := strings.ToLower(os.Getenv("NETWORK"))
	snapshotsPath := util.GetEnvString("SNAPSHOTS_PATH", "/var/run/tezos/snapshots")
	mavkitNodepath := util.GetEnvString("MAVKIT_NODE_PATH", "/usr/local/bin/mavkit-node")
	mavrykPath := util.GetEnvString("MAVRYK_PATH", "/var/run/tezos/node")

	snapshotExec := NewSnapshotExec(snapshotsPath, mavkitNodepath, mavrykPath)

	if bucketName == "" {
		log.Fatalln("The BUCKET_NAME environment variable is empty.")
	}

	if network == "" {
		log.Fatalln("The NETWORK environment variable is empty.")
	}

	// Usually for basenet, because it's link
	if strings.Contains(network, "https://testnets.mavryk.network/") {
		network = strings.Replace(network, "https://testnets.mavryk.network/", "", -1)
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	snapshotStorage := store.NewSnapshotStorage(client, bucketName)

	// Check if today the rolling snapshot already exists
	execute(ctx, snapshotStorage, snapshot.ROLLING, network, snapshotExec)

	// Check if today the full snapshot already exists
	execute(ctx, snapshotStorage, snapshot.FULL, network, snapshotExec)

	snapshotStorage.DeleteExpiredSnapshots(ctx, maxDays, maxMonths)

	// Delete local snapshots
	snapshotExec.DeleteLocalSnapshots()

	log.Printf("Snapshot job took %s", time.Since(start))
}

func execute(ctx context.Context, snapshotStorage *store.SnapshotStorage, historyMode snapshot.HistoryModeType, chain string, snapshotExec *SnapshotExec) {
	todayItems := snapshotStorage.GetTodaySnapshotsItems(ctx)

	alreadyExist := lo.SomeBy(todayItems, func(item snapshot.SnapshotItem) bool {
		return item.ChainName == chain && item.HistoryMode == historyMode
	})

	if alreadyExist {
		log.Printf("Already exist a today snapshot with chain: %s and history mode: %s. \n", chain, string(historyMode))
		return
	}

	snapshotExec.CreateSnapshot(historyMode)
	snapshotfilename, err := snapshotExec.GetSnapshotName(historyMode)
	if err != nil {
		log.Fatalf("%v \n", err)
	}
	snapshotHeaderOutput := snapshotExec.GetSnapshotHeaderOutput(snapshotfilename)

	snapshotStorage.EphemeralUpload(ctx, snapshotfilename, snapshotHeaderOutput)
}
