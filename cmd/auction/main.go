package main

import (
	"auction-house/internal/auction"
	"auction-house/internal/storage"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/google/uuid"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return errors.New("expected commands")
	}

	repo := storage.NewMemoryRepository()
	service := auction.NewService(repo)

	switch args[0] {
	case "create":
		return handleCreate(service, args[1:])
	case "bid":
		return handleBid(service, args[1:])
	case "close":
		return handleClose(service, args[1:])
	case "inspect":
		return handleInspect(service, args[1:])
	case "list":
		return handleList(service)
	default:
		return errors.New("wrong command provided")
	}
}

func handleCreate(service *auction.Service, args []string) error {
	fs := flag.NewFlagSet("create", flag.ExitOnError)

	productName := fs.String("product-name", "", "name of the product being auctioned")
	startingPrice := fs.Float64("starting-price", 0.0, "initial price point")
	minIncrement := fs.Float64("min-increment", 0.0, "minimum increment over the current highest bid")
	buyoutPrice := fs.Float64("buyout-price", 0.0, "buyout price")

	parseArgs(fs, args)

	id, err := service.HandleCreate(*productName, *startingPrice, *minIncrement, *buyoutPrice)

	if err != nil {
		return fmt.Errorf("There was a problem creating an auction: %w", err)
	}

	fmt.Println(id)
	return nil
}

func handleBid(service *auction.Service, args []string) error {
	fs := flag.NewFlagSet("bid", flag.ExitOnError)

	//auctionId := fs.String("auction-id", "", "id of the auction to bid")
	//amount := fs.Float64("amount", 0.0, "amount of auction to bid")

	parseArgs(fs, args)

	return nil
}

func handleClose(service *auction.Service, args []string) error {
	fs := flag.NewFlagSet("close", flag.ExitOnError)

	//auctionId := fs.String("auction-id", "", "id of the auction to close")

	parseArgs(fs, args)

	return nil
}

func handleInspect(service *auction.Service, args []string) error {
	fs := flag.NewFlagSet("inspect", flag.ExitOnError)

	auctionId := fs.String("auction-id", "", "id of the auction to inspect")

	parseArgs(fs, args)
	id, err := uuid.Parse(*auctionId)
	if err != nil {
		return errors.New("there was an error parsing id")
	}

	auction, err := service.HandleGet(id)
	if err != nil {
		return fmt.Errorf("there was an error getting auction: %w", err)
	}

	fmt.Println(auction.String())

	return nil
}

func handleList(service *auction.Service) error {
	fmt.Println("List of all of the auctions")

	return nil
}

func parseArgs(fs *flag.FlagSet, args []string) {
	err := fs.Parse(args)
	if err != nil {
		fmt.Printf("Error parsing args: %v\n", err)
	}
}
