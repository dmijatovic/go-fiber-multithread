# Fiber basic api

This repo tests basic fiber api setup while learning go at the same time ;-).

Fiber documentation is [here](https://docs.gofiber.io/)


## Initialize Go module

```bash
# initialize new app/module
go mod init fiber-api
```

## Install fiber package

```bash
# install fiber package
go get github.com/gofiber/fiber/v2
```

## Run fiber using multiple threads

```go

func main() {
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
	// SHUTDOWN
	_ = app.Shutdown()

	fmt.Println("Gracefully shutting down...")
}

func startApi(app *fiber.App) {
	// listen to 8080 port
	err := app.Listen(":8080")
	if err != nil {
		panic(err)
	}
}

```