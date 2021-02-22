package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

func WaitForEnd(scanner *bufio.Scanner) {
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "<END") {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Scanner failed:", err.Error())
		os.Exit(1)
	}
}

func GetDirection(conn net.Conn, scanner *bufio.Scanner, id int) int {
	fmt.Fprintf(conn, "get(%d, dir)\n", id)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "<END") {
			break
		} else if !strings.HasPrefix(line, "<") {
			var _id, dir int
			_, err := fmt.Sscanf(line, "%d dir[%d]", &_id, &dir)
			if err != nil {
				fmt.Println("Scanf failed:", err.Error())
				os.Exit(1)
			}
			if _id == id {
				return dir
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Scanner failed:", err.Error())
		os.Exit(1)
	}

	fmt.Println("Didn't receive direction")
	os.Exit(1)
	return 0
}

func SetDirection(conn net.Conn, scanner *bufio.Scanner, id int, dir int) {
	fmt.Fprintf(conn, "set(%d, dir[%d])\n", id, dir)
	WaitForEnd(scanner)
}

func SetSpeed(conn net.Conn, scanner *bufio.Scanner, id int, speed int) {
	fmt.Fprintf(conn, "set(%d, speed[%d])\n", id, speed)
	WaitForEnd(scanner)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	conn, err := net.Dial("tcp", "ecos.local:15471")
	if err != nil {
		fmt.Println("Dial failed:", err.Error())
		os.Exit(1)
	}

	scanner := bufio.NewScanner(conn)

	var dir int
	dir = GetDirection(conn, scanner, 1001)

	// Assume train is pointing in the direction to move away from the buffers
	for {
		var move, slow time.Duration
		if dir == 0 {
			move = 72500
			slow = 10000
		} else {
			move = 68500
			slow = 10000
		}

		fmt.Println("Moving..")
		SetSpeed(conn, scanner, 1001, 8*127/28)
		time.Sleep(move * time.Millisecond)

		fmt.Println("Slowing..")
		SetSpeed(conn, scanner, 1001, 4*127/28)
		time.Sleep(slow * time.Millisecond)

		fmt.Println("Stopping..")
		SetSpeed(conn, scanner, 1001, 0)

		if dir == 0 {
			dir = 1
		} else {
			dir = 0
		}
		SetDirection(conn, scanner, 1001, dir)

		n := 15 + rand.Intn(10)*30
		fmt.Printf("Waiting %d.\n\n", n)
		time.Sleep(time.Duration(n) * time.Second)
	}
}
