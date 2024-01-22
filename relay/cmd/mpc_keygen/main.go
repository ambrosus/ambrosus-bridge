package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/networking/client"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/networking/server"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/tss_wrap"
	"github.com/bnb-chain/tss-lib/ecdsa/keygen"
	zerolog "github.com/rs/zerolog/log"
)

// keugen examples:
// keygen server:
// go run main.go -server :8080 -meID A -partyIDs "A B C" -threshold 2 -shareDir /tmp/mpc
// keygen client:
// go run main.go -url http://localhost:8080 -meID B -partyIDs "A B C" -threshold 2 -shareDir /tmp/mpc

// reshare examples:
// run as server and client, member of old and new committees
// go run main.go -reshare -server :8080 -url http://localhost:8080 -meID A -partyIDs "A B C" -threshold 2 -meIDNew A2 -partyIDsNew "A2 B2 C2 D2" -thresholdNew 3 -shareDir /tmp/mpc
// run as client, member of old and new committees
// go run main.go -reshare -url http://localhost:8080 -meID B -partyIDs "A B C" -threshold 2 -meIDNew B2 -partyIDsNew "A2 B2 C2 D2" -thresholdNew 3 -shareDir /tmp/mpc
// run as client, member of new committee
// go run main.go -reshare -url http://localhost:8080 partyIDs "A B C" -threshold 2 -meIDNew C2 -partyIDsNew "A2 B2 C2 D2" -thresholdNew 3 -shareDir /tmp/mpc

var logger = zerolog.Logger

func main() {
	flagOperation := flag.Bool("reshare", false, "do reshare (default: keygen)")

	flagServerHost := flag.String("server", "", "if specified, run a server on this port (ex: ':8080')")
	flagUrl := flag.String("url", "", "url to which a client will connect")
	flagAccessToken := flag.String("accessToken", "", "url to which a client will connect")

	flagMeID := flag.String("meID", "", "my ID")
	flagPartyIDs := flag.String("partyIDs", "", "party IDs (space separated)")
	flagThreshold := flag.Int("threshold", -1, "threshold")

	flagShareDir := flag.String("shareDir", "", "path to directory with shares")

	// for reshare only
	flagMeIDNew := flag.String("meIDNew", "", "new my ID (for reshare)")
	flagPartyIDsNew := flag.String("partyIDsNew", "", "new party IDs (space separated) (for reshare)")
	flagThresholdNew := flag.Int("thresholdNew", -1, "new threshold (for reshare)")

	flag.Parse()

	partyIDs := strings.Split(*flagPartyIDs, " ")
	isServer := *flagServerHost != ""
	checkThreshold(*flagThreshold)
	checkPartyIDs(partyIDs)
	checkShareDir(*flagShareDir)

	if !*flagOperation {
		doKeygen(isServer, *flagServerHost, *flagUrl, *flagAccessToken, *flagMeID, partyIDs, *flagThreshold, *flagShareDir)
	} else {
		partyIDsNew := strings.Split(*flagPartyIDsNew, " ")
		checkThreshold(*flagThresholdNew)
		checkPartyIDs(partyIDsNew)

		reshareBothCommittee(isServer, *flagServerHost, *flagUrl, *flagAccessToken, *flagMeID, *flagMeIDNew, partyIDs, partyIDsNew, *flagThreshold, *flagThresholdNew, *flagShareDir)
	}
}

func doKeygen(isServer bool, hostUrl, serverURL string, accessToken string, id string, partyIDs []string, threshold int, shareDir string) {
	sharePath := getSharePath(shareDir, id)

	fmt.Println("=======================================================")
	fmt.Println("You are about to generate the MPC share")
	fmt.Println("IDS: ", partyIDs, "; threshold: ", threshold)
	fmt.Println("Your ID: ", id, "; share path: ", sharePath)
	fmt.Println("Is this server: ", isServer, "; server URL: ", serverURL)
	fmt.Println("=======================================================")

	if isShareExist(sharePath) {
		log.Fatal("share already exist")
	}

	log.Println("Generating some safe prime numbers...")
	preParams, err := keygen.GeneratePreParams(10 * time.Minute)
	if err != nil {
		log.Fatal("Failed to generate pre params", err)
	}
	log.Println("Done generating safe prime numbers!")
	log.Println("Starting MPC keygen protocol...")

	mpcc := tss_wrap.NewMpc(id, threshold, &logger)
	netOperation := createNetworking(isServer, hostUrl, serverURL, accessToken, mpcc)
	err = netOperation.Keygen(context.Background(), partyIDs, *preParams)
	if err != nil {
		log.Fatal(err)
	}
	saveShare(mpcc, sharePath)
}

func doReshare(isServer bool, hostUrl, serverURL, accessToken string, id string, meInNewCommittee bool, partyIDsOld, partyIDsNew []string, thresholdOld, thresholdNew int, shareDir string) {
	sharePath := getSharePath(shareDir, id)

	fmt.Println("=======================================================")
	fmt.Println("You are about to reshare the MPC share")
	fmt.Println("Old IDS: ", partyIDsOld, "; threshold: ", thresholdOld)
	fmt.Println("New IDS: ", partyIDsNew, "; threshold: ", thresholdNew)
	fmt.Println("Your ID: ", id, "; share path: ", sharePath)
	fmt.Println("Is your in new committee: ", meInNewCommittee)
	fmt.Println("Is this server: ", isServer, "; server URL: ", serverURL)
	fmt.Println("=======================================================")

	mpcc := tss_wrap.NewMpc(id, thresholdOld, &logger)
	if !meInNewCommittee {
		readShare(mpcc, sharePath)
	}

	netOperation := createNetworking(isServer, hostUrl, serverURL, accessToken, mpcc)
	err := netOperation.Reshare(context.Background(), partyIDsOld, partyIDsNew, thresholdNew)
	if err != nil {
		log.Fatal(err)
	}

	// save share only if we are in new committee
	// coz old share changes somehow after reshare (Xi component zeroed)
	if meInNewCommittee {
		saveShare(mpcc, sharePath)
	}
}

func reshareBothCommittee(isServer bool, serverHost, url, accessToken, meID, meIDNew string, partyIDs, partyIDsNew []string, threshold, thresholdNew int, shareDir string) {
	var wg sync.WaitGroup
	if meIDNew != "" { // we are in new committee
		wg.Add(1)
		go func() {
			doReshare(isServer, serverHost, url, accessToken,
				meIDNew, true, partyIDs, partyIDsNew,
				threshold, thresholdNew, shareDir)
			wg.Done()
		}()

		// if new committee just launched as server, old committee must be run as client, coz can't run server twice
		time.Sleep(2 * time.Second)
		isServer = false
	}
	if meID != "" { // we are in old committee
		doReshare(isServer, serverHost, url, accessToken,
			meID, false, partyIDs, partyIDsNew,
			threshold, thresholdNew, shareDir)
	}
	wg.Wait()
}

// utils

type networkingOperations interface {
	Reshare(ctx context.Context, partyIDsOld, partyIDsNew []string, thresholdNew int) error
	Keygen(ctx context.Context, party []string, optionalPreParams ...keygen.LocalPreParams) error
}

func createNetworking(isServer bool, hostUrl string, serverURL string, accessToken string, mpcc *tss_wrap.Mpc) networkingOperations {
	if isServer {
		server_ := server.NewServer(mpcc, accessToken, &logger)
		go http.ListenAndServe(hostUrl, server_)
		return server_
	} else {
		return client.NewClient(mpcc, serverURL, accessToken, &logger)
	}
}

// share file utils

func saveShare(tss *tss_wrap.Mpc, sharePath string) {
	if share, err := tss.Share(); err != nil {
		log.Fatal(fmt.Errorf("can't get share: %w", err))
	} else if err = os.WriteFile(sharePath, share, 0644); err != nil {
		log.Fatal(fmt.Errorf("can't save share: %w", err))
	}
}

func readShare(tss *tss_wrap.Mpc, sharePath string) {
	if share, err := os.ReadFile(sharePath); err != nil {
		log.Fatal(fmt.Errorf("can't read share: %w", err))
	} else if err = tss.SetShare(share); err != nil {
		log.Fatal(fmt.Errorf("can't set share: %w", err))
	}
}

func isShareExist(sharePath string) bool {
	_, err := os.Stat(sharePath)
	return err == nil
}
func getSharePath(shareDir, id string) string {
	return filepath.Join(shareDir, fmt.Sprintf("share_%s", id))
}

// flags checks

func checkShareDir(dirPath string) {
	fileInfo, err := os.Stat(dirPath)
	if err != nil {
		log.Fatalf("something wring with dirPath (%v): %v", dirPath, err)
	}
	if !fileInfo.IsDir() {
		log.Fatalf("dirPath (%v) is not a directory", dirPath)
	}
}

func checkThreshold(t int) {
	if t < 2 {
		log.Fatal("threshold must be >= 2")
	}
}
func checkPartyIDs(partyIDs []string) {
	if len(partyIDs) < 2 {
		log.Fatal("partyIDs must be >= 2")
	}
}
