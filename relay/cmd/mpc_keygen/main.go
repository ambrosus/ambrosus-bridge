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
	zerolog "github.com/rs/zerolog/log"
)

// examples:
// keygen server:
// go run main.go -server -serverUrl :8080 -meID A -partyIDs "A B C" -threshold 2 -shareDir /tmp/mpc
// keygen client:
// go run main.go -serverUrl http://localhost:8080 -meID B -partyIDs "A B C" -threshold 2 -shareDir /tmp/mpc

// reshare server:
// go run main.go -reshare -server -serverUrl :8080 -meID A -partyIDs "A B C" -threshold 2 -meIdNew A2 -partyIDsNew "A2 B2 C2 D2" -thresholdNew 3 -shareDir /tmp/mpc
// reshare client:
// go run main.go -reshare -serverUrl http://localhost:8080 -meID B -partyIDs "A B C" -threshold 2 -meIdNew B2 -partyIDsNew "A2 B2 C2 D2" -thresholdNew 3 -shareDir /tmp/mpc

func main() {
	flagOperation := flag.Bool("reshare", false, "do reshare (default: keygen)")

	flagIsServer := flag.Bool("server", false, "is this server (default: client)")
	flagServerUrl := flag.String("serverUrl", "", "server url (use ':8080' for server)")

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
	checkThreshold(*flagThreshold)
	checkPartyIDs(partyIDs)
	checkShareDir(*flagShareDir)

	if !*flagOperation {
		keygen(*flagIsServer, *flagServerUrl, *flagMeID, partyIDs, *flagThreshold, *flagShareDir)
	} else {
		partyIDsNew := strings.Split(*flagPartyIDsNew, " ")
		checkThreshold(*flagThresholdNew)
		checkPartyIDs(partyIDsNew)

		var wg sync.WaitGroup
		if *flagMeIDNew != "" {
			// we are in new committee
			wg.Add(1)
			go func() {
				reshare(*flagIsServer, *flagServerUrl,
					*flagMeIDNew,
					partyIDs, partyIDsNew,
					*flagThresholdNew, *flagThreshold, *flagShareDir)
				wg.Done()
			}()
			// if we already runned as server set flagIsServer to false, coz can't run server twice
			*flagIsServer = false
		}
		if *flagMeID != "" {
			// we are in old committee
			reshare(*flagIsServer, *flagServerUrl,
				*flagMeID,
				partyIDs, partyIDsNew,
				*flagThreshold, *flagThresholdNew, *flagShareDir)
		}
		wg.Wait()
	}
}

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

func keygen(isServer bool, serverURL string, id string, partyIDs []string, threshold int, shareDir string) {
	sharePath := getSharePath(shareDir, id)
	if isShareExist(sharePath) {
		log.Fatal("share already exist")
	}

	fmt.Println("=======================================================")
	fmt.Println("You are about to generate the MPC share")
	fmt.Println("IDS: ", partyIDs, "; threshold: ", threshold)
	fmt.Println("Your ID: ", id, "; share path: ", sharePath)
	fmt.Println("Is this server: ", isServer, "; server URL: ", serverURL)
	fmt.Println("=======================================================")

	logger := zerolog.Logger
	mpcc := tss_wrap.NewMpc(id, threshold, &logger)

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Minute)
	if isServer {
		server_ := server.NewServer(mpcc, &logger)
		go http.ListenAndServe(serverURL, server_)

		err := server_.Keygen(ctx, partyIDs)
		if err != nil {
			panic(err)
		}
	} else {
		client_ := client.NewClient(mpcc, serverURL, nil, &logger)

		err := client_.Keygen(ctx, partyIDs)
		if err != nil {
			panic(err)
		}
	}
	saveShare(mpcc, sharePath)
}

func reshare(isServer bool, serverURL, id string, partyIDsOld, partyIDsNew []string, thresholdOld, thresholdNew int, shareDir string) {
	sharePath := getSharePath(shareDir, id)

	fmt.Println("=======================================================")
	fmt.Println("You are about to reshare the MPC share")
	fmt.Println("Old IDS: ", partyIDsOld, "; threshold: ", thresholdOld)
	fmt.Println("New IDS: ", partyIDsNew, "; threshold: ", thresholdNew)
	fmt.Println("Your ID: ", id, "; share path: ", sharePath)
	fmt.Println("Is this server: ", isServer, "; server URL: ", serverURL)
	fmt.Println("=======================================================")

	logger := zerolog.Logger
	mpcc := tss_wrap.NewMpc(id, thresholdOld, &logger)
	if isShareExist(sharePath) {
		share, err := os.ReadFile(sharePath)
		if err != nil {
			panic(fmt.Errorf("can't read share: %w", err))
		}
		mpcc.SetShare(share)
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Minute)
	if isServer {
		server_ := server.NewServer(mpcc, &logger)
		go http.ListenAndServe(serverURL, server_)

		err := server_.Reshare(ctx, partyIDsOld, partyIDsNew)
		if err != nil {
			panic(err)
		}
	} else {
		client_ := client.NewClient(mpcc, serverURL, nil, &logger)

		err := client_.Reshare(ctx, partyIDsOld, partyIDsNew)
		if err != nil {
			panic(err)
		}
	}
	saveShare(mpcc, sharePath)
}

func saveShare(tss *tss_wrap.Mpc, sharePath string) {
	share, err := tss.Share()
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(sharePath, share, 0644)
	if err != nil {
		panic(err)
	}
}

func isShareExist(sharePath string) bool {
	_, err := os.Stat(sharePath)
	return err == nil
}
func getSharePath(shareDir, id string) string {
	return filepath.Join(shareDir, fmt.Sprintf("share_%s", id))
}
