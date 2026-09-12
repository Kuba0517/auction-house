package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Expected command")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "create":
		handleCreate(os.Args[2:])
	case "bid":
		handleBid(os.Args[2:])
	case "close":
		handleClose(os.Args[2:])
	case "inspect":
		handleInspect(os.Args[2:])
	case "list":
		handleList()
	default:
		fmt.Println("Wrong command provided")
	}

}

func handleCreate(args []string) {
	fs := flag.NewFlagSet("create", flag.ExitOnError)

	name := fs.String("product-name", "", "name of the product being auctioned")
	startingPrice := fs.Float64("starting-price", 0.0, "initial price point")
	minIncrement := fs.Float64("min-increment", 0.0, "minimum increment over the current highest bid")

	parseArgs(fs, args)

	fmt.Println(*name)
	fmt.Println(*startingPrice)
	fmt.Println(*minIncrement)
}

func handleBid(args []string) {
	fs := flag.NewFlagSet("bid", flag.ExitOnError)

	//auctionId := fs.String("auction-id", "", "id of the auction to bid")
	//amount := fs.Float64("amount", 0.0, "amount of auction to bid")

	parseArgs(fs, args)
}

func handleClose(args []string) {
	fs := flag.NewFlagSet("close", flag.ExitOnError)

	//auctionId := fs.String("auction-id", "", "id of the auction to close")

	parseArgs(fs, args)
}

func handleInspect(args []string) {
	fs := flag.NewFlagSet("inspect", flag.ExitOnError)

	//auctionId := fs.String("auction-id", "", "id of the auction to inspect")

	parseArgs(fs, args)
}

func handleList() {
	fmt.Println("List of all of the auctions")
}

func parseArgs(fs *flag.FlagSet, args []string) {
	err := fs.Parse(args)
	if err != nil {
		fmt.Printf("Error parsing args: %v\n", err)
	}
}
