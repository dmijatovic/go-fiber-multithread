package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"

	"github.com/gofiber/fiber/v2"
)

func Home(ctx *fiber.Ctx) error {
	return ctx.SendString("Hello world")
}

type SomeData struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

var Data = []SomeData{}

func getData(ctx *fiber.Ctx) error {
	// return data from array
	return ctx.Status(200).JSON(Data)
}

func postData(ctx *fiber.Ctx) error {
	// get request body
	body := new(SomeData)
	// parse data into structure
	err := ctx.BodyParser(body)
	if err != nil {
		// return error
		ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		return err
	}

	newData := SomeData{
		Id:   body.Id,
		Name: body.Name,
	}
	// append new data to data array
	Data = append(Data, newData)
	// return data
	return ctx.Status(fiber.StatusOK).JSON(newData)
}

func startApi(app *fiber.App) {
	// listen to 8080 port
	err := app.Listen(":8080")
	if err != nil {
		panic(err)
	}
}

func main() {
	// define max number of processes
	// by default prefork will spawn process on each processor core
	runtime.GOMAXPROCS(4)

	// create fiber app
	app := fiber.New(fiber.Config{
		// enable multiple processes to run
		Prefork: true,
	})

	// ROUTES
	app.Get("/", Home)
	app.Get("/data", getData)
	app.Post("/data", postData)

	// start api in separate thread
	go startApi(app)

	// Create channel to signify a signal being sent
	c := make(chan os.Signal, 1)
	// When an interrupt is sent, notify the channel
	signal.Notify(c, os.Interrupt)
	// This blocks the main thread until an interrupt is received
	_ = <-c

	// Log closing of app processes
	if fiber.IsChild() {
		fmt.Println("Gracefully shutting down child process...")
	} else {
		fmt.Println("Gracefully shutting down parent process...")
	}
	// SHUTDOWN
	_ = app.Shutdown()
}
