package main

import (
	"errors"
	"fmt"
	"os"
)

type Seat struct {
	row     int
	col     int
	reseved bool
}

type Cinema struct {
	room          [][]Seat
	currentIncome int
}

type App struct {
	cinema *Cinema
}

func main() {
	app := &App{}

	app.initCinema()

	for {
		app.showMenu()
		app.handleUserEvent()
	}
}

func (a *App) showMenu() {
	fmt.Println("1. Show the seats")
	fmt.Println("2. Buy a ticket")
	fmt.Println("3. Statistics")
	fmt.Println("0. Exit")
}

func (a *App) handleUserEvent() {
	var event int
	fmt.Scan(&event)

	switch event {
	case 1:
		a.showCinema()
	case 2:
		a.buyTicket()
	case 3:
		a.showStatistics()
	case 0:
		os.Exit(0)
	}
}

func (a *App) initCinema() {
	var rows, seats int
	fmt.Println("Enter the number of rows:")
	fmt.Scan(&rows)
	fmt.Println("Enter the number of seats in each row:")
	fmt.Scan(&seats)

	a.cinema = NewCinema(rows, seats)
}

func (a *App) showCinema() {
	a.cinema.showRoom()
}

func (a *App) buyTicket() {

	for {
		var row, seat int
		fmt.Println("Enter a row number:")
		fmt.Scan(&row)
		fmt.Println("Enter a seat number in that row:")
		fmt.Scan(&seat)

		if price, err := a.cinema.buy(row, seat); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Ticket price: $%d\n", price)
			break
		}
	}
}

func (a *App) showStatistics() {
	ticketsCount := a.cinema.getTicketsCount()
	fmt.Printf("Number of purchased tickets: %d\n", ticketsCount)

	percentage := float64(ticketsCount*100) / float64(a.cinema.getSeatsCount())
	fmt.Printf("Percentage: %.2f%%\n", percentage)
	fmt.Printf("Current income: $%d\n", a.cinema.currentIncome)
	fmt.Printf("Total income: $%d\n", a.cinema.totalIncome())
}

func NewSeat(row, col int) Seat {
	return Seat{row: row, col: col, reseved: false}
}

func (s *Seat) symbol() string {
	if s.reseved {
		return "B"
	}
	return "S"
}

func NewCinema(rows, seats int) *Cinema {
	cinema := &Cinema{}
	cinema.initRoom(rows, seats)
	return cinema
}

func (c *Cinema) initRoom(rows, seats int) {
	room := make([][]Seat, rows)
	for r := 0; r < rows; r++ {
		row := make([]Seat, seats)
		for s := 0; s < seats; s++ {
			row[s] = NewSeat(r+1, s+1)
		}
		room[r] = row
	}
	c.room = room
}

func (c *Cinema) showRoom() {
	fmt.Println("Cinema:")

	for r := 0; r < len(c.room)+1; r++ {
		for s := 0; s < len(c.room[0])+1; s++ {
			switch {
			case r == 0 && s == 0:
				fmt.Print("  ")
			case r == 0:
				fmt.Printf("%d ", s)
			case s == 0:
				fmt.Printf("%d ", r)
			default:
				seat, _ := c.getSeat(r, s)
				fmt.Printf("%s ", seat.symbol())
			}
		}
		fmt.Println()
	}
}

func (c *Cinema) buy(row, seat int) (int, error) {
	s, err := c.getSeat(row, seat)
	if err != nil {
		return 0, err
	} else if s.reseved {
		return 0, errors.New("That ticket has already been purchased!")
	} else {
		price := c.getPrice(*s)
		s.reseved = true
		c.currentIncome += price
		return price, nil
	}
}

func (c *Cinema) totalIncome() int {
	totalIncome := 0

	for _, row := range c.room {
		for _, seat := range row {
			totalIncome += c.getPrice(seat)
		}
	}

	return totalIncome
}

func (c *Cinema) getSeat(row, seat int) (*Seat, error) {
	if !c.isValid(row, seat) {
		return nil, errors.New("Wrong input!")
	}

	return &c.room[row-1][seat-1], nil
}

func (c *Cinema) isValid(row, seat int) bool {
	if row > len(c.room) || row < 1 || seat > len(c.room[0]) || seat < 1 {
		return false
	}
	return true
}

func (c *Cinema) getPrice(s Seat) int {

	const (
		ordinaryRoomNumOfSeats = 60
		price                  = 10
		priceWithDiscont       = 8
	)

	if c.getSeatsCount() <= ordinaryRoomNumOfSeats {
		return price
	}

	frontRows := len(c.room) / 2

	if s.row > frontRows {
		return priceWithDiscont
	}

	return price
}

func (c *Cinema) getSeatsCount() int {
	return len(c.room) * len(c.room[0])
}

func (c *Cinema) getTicketsCount() int {
	count := 0

	for _, row := range c.room {
		for _, seat := range row {
			if seat.reseved {
				count++
			}
		}
	}

	return count
}
