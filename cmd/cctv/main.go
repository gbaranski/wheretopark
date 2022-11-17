package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	wheretopark "wheretopark/go"
	"wheretopark/providers/cctv"

	"github.com/caarlos0/env/v6"
	"gocv.io/x/gocv"
)

type environment struct {
	DatabaseURL      string  `env:"DATABASE_URL" envDefault:"ws://localhost:8000"`
	DatabaseName     string  `env:"DATABASE_NAME" envDefault:"development"`
	DatabaseUser     string  `env:"DATABASE_USER" envDefault:"root"`
	DatabasePassword string  `env:"DATABASE_PASSWORD" envDefault:"root"`
	Configuration    *string `env:"CONFIGURATION"`
}

func RunParkingLot(parkingLot cctv.ParkingLot, window *gocv.Window) error {
	fmt.Printf("running parking lot %s\n", parkingLot.Name)
	video, err := gocv.OpenVideoCapture(parkingLot.CameraURL)
	if err != nil {
		panic(err)
	}
	defer video.Close()
	img := gocv.NewMat()
	for {
		fmt.Printf("reading frame\n")
		if ok := video.Read(&img); !ok {
			fmt.Printf("cannot read video\n")
			return err
		}
		fmt.Printf("got frame\n")

		for _, spot := range parkingLot.Spots {
			croppedImage := spot.CropOn(img)
			window.IMShow(croppedImage)
			for {
				if window.WaitKey(1) >= 0 {
					break
				}
			}
		}
		// time.Sleep(10 * time.Second)
	}
}

func main() {
	window := gocv.NewWindow("WhereToPark")
	defer window.Close()

	environment := environment{}
	if err := env.Parse(&environment); err != nil {
		log.Fatalf("%+v\n", err)
	}

	url, err := url.Parse(environment.DatabaseURL)
	if err != nil {
		log.Fatalf("invalid database url: %s", err)
	}
	client, err := wheretopark.NewClient(url, "wheretopark", environment.DatabaseName)
	if err != nil {
		log.Fatalf("failed to create database client: %v", err)
	}
	defer client.Close()
	err = client.SignInWithPassword(environment.DatabaseUser, environment.DatabasePassword)
	if err != nil {
		log.Fatalf("failed to sign in: %v", err)
	}

	var configuration cctv.Configuration
	if environment.Configuration == nil {
		configuration = cctv.DefaultConfiguration
	} else {
		newConfiguration, err := cctv.LoadConfiguration(*environment.Configuration)
		if err != nil {
			panic(err)
		}
		configuration = *newConfiguration
	}

	RunParkingLot(configuration.ParkingLots[0], window)
	// for _, parkingLot := range configuration.ParkingLots {
	// go RunParkingLot(parkingLot)
	// }

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	// video, err := gocv.OpenVideoCapture("http://91.238.55.4:5080/LiveApp/streams/435465478973256862461988.m3u8?token=null")
	// if err != nil {
	// 	panic(err)
	// }
	// defer video.Close()
	// window := gocv.NewWindow("Face Detect")
	// defer window.Close()
	// img := gocv.NewMat()
	// defer img.Close()
	// if ok := video.Read(&img); !ok {
	// 	fmt.Printf("cannot read video\n")
	// 	return
	// }
	// window.IMShow(img)
	// for {
	// 	if window.WaitKey(1) >= 0 {
	// 		break
	// 	}
	// }
}
