package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/briancain/go-tetris/test/integration"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run manual_client.go <username>")
		os.Exit(1)
	}

	username := os.Args[1]
	serverURL := "http://localhost:8080"

	fmt.Printf("🎮 Starting test client for user: %s\n", username)
	fmt.Printf("📡 Server: %s\n", serverURL)

	// Create client
	client := integration.NewTestClient(username, serverURL)

	// Login
	fmt.Print("🔐 Logging in... ")
	err := client.Login()
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}
	fmt.Printf("✅ Logged in as %s (ID: %s)\n", client.Username, client.PlayerID)

	// Connect WebSocket
	fmt.Print("🔌 Connecting WebSocket... ")
	err = client.ConnectWebSocket()
	if err != nil {
		log.Fatalf("WebSocket connection failed: %v", err)
	}
	fmt.Println("✅ WebSocket connected")

	defer client.Close()

	// Start message monitor
	go monitorMessages(client)

	// Interactive commands
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("\n📋 Available commands:")
	fmt.Println("  join     - Join matchmaking queue")
	fmt.Println("  status   - Check queue status")
	fmt.Println("  move <type> - Send game move (left/right/rotate/drop)")
	fmt.Println("  ping     - Send ping")
	fmt.Println("  quit     - Exit")
	fmt.Println()

	for {
		fmt.Print("💬 Command: ")
		if !scanner.Scan() {
			break
		}

		command := strings.TrimSpace(scanner.Text())
		parts := strings.Split(command, " ")

		switch parts[0] {
		case "join":
			handleJoinQueue(client)
		case "status":
			handleQueueStatus(client)
		case "move":
			if len(parts) > 1 {
				handleGameMove(client, parts[1])
			} else {
				fmt.Println("❌ Usage: move <type>")
			}
		case "ping":
			handlePing(client)
		case "quit":
			fmt.Println("👋 Goodbye!")
			return
		default:
			fmt.Printf("❌ Unknown command: %s\n", command)
		}
	}
}

func handleJoinQueue(client *integration.TestClient) {
	fmt.Print("🎯 Joining queue... ")
	err := client.JoinQueue()
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
		return
	}
	fmt.Println("✅ Joined queue")
}

func handleQueueStatus(client *integration.TestClient) {
	fmt.Print("📊 Checking queue status... ")
	position, err := client.GetQueueStatus()
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
		return
	}

	if position >= 0 {
		fmt.Printf("✅ Position in queue: %d\n", position)
	} else {
		fmt.Println("✅ Not in queue")
	}
}

func handleGameMove(client *integration.TestClient, moveType string) {
	fmt.Printf("🎮 Sending move: %s... ", moveType)
	err := client.SendGameMove(moveType)
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
		return
	}
	fmt.Println("✅ Move sent")
}

func handlePing(client *integration.TestClient) {
	fmt.Print("🏓 Sending ping... ")
	err := client.SendPing()
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
		return
	}
	fmt.Println("✅ Ping sent")
}

func monitorMessages(client *integration.TestClient) {
	lastCount := 0
	for {
		messages := client.GetMessages()
		if len(messages) > lastCount {
			for i := lastCount; i < len(messages); i++ {
				msg := messages[i]
				msgType, _ := msg["type"].(string)

				switch msgType {
				case "match_found":
					gameID, _ := msg["gameId"].(string)
					opponent, _ := msg["opponent"].(string)
					seed, _ := msg["seed"].(float64)
					fmt.Printf("\n🎉 MATCH FOUND! Game: %s, Opponent: %s, Seed: %.0f\n",
						gameID, opponent, seed)
				case "game_move":
					moveType, _ := msg["moveType"].(string)
					playerID, _ := msg["playerId"].(string)
					fmt.Printf("\n🎮 Opponent move: %s (from %s)\n", moveType, playerID)
				case "game_state":
					score, _ := msg["score"].(float64)
					level, _ := msg["level"].(float64)
					fmt.Printf("\n📊 Opponent state: Score=%.0f, Level=%.0f\n", score, level)
				case "game_over":
					winnerID, _ := msg["winnerId"].(string)
					fmt.Printf("\n🏁 GAME OVER! Winner: %s\n", winnerID)
				case "pong":
					fmt.Printf("\n🏓 Pong received\n")
				default:
					fmt.Printf("\n📨 Message: %v\n", msg)
				}
				fmt.Print("💬 Command: ")
			}
			lastCount = len(messages)
		}
		time.Sleep(100 * time.Millisecond)
	}
}
