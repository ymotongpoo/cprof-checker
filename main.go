// Copyright 2020 Yoshi Yamaguchi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"cloud.google.com/go/profiler"
)

const service = "cprof-checker"

const (
	height = 2048
	width  = 2048
)

var projectID = flag.String("p", "", "Google Cloud Platform project ID")

func main() {
	flag.Parse()

	user := os.Getenv("USER")

	cfg := profiler.Config{
		Service:        service,
		ServiceVersion: user + "-sample",
	}

	if *projectID != "" {
		cfg.ProjectID = *projectID
	}

	if err := profiler.Start(cfg); err != nil {
		panic(err)
	}

	for i := 0; ; i++ {
		log.Printf("cycle %v\n", i)
		if err := dummyHeavyLoadProcess(); err != nil {
			log.Println(err)
		}
	}
}

// dummyHeavyLoadProcess creates a random png file.
func dummyHeavyLoadProcess() error {
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	rand.Seed(time.Now().UnixNano())
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.NRGBA{
				R: uint8(rand.Intn(math.MaxUint8)),
				G: uint8(rand.Intn(math.MaxUint8)),
				B: uint8(rand.Intn(math.MaxUint8)),
				A: 255,
			})
		}
	}
	f, err := os.Create("cprof-sample.png")
	if err != nil {
		return err
	}
	if err := png.Encode(f, img); err != nil {
		f.Close()
		return err
	}
	stat, err := f.Stat()
	if err != nil {
		return err
	}
	log.Printf("file size: %v\n", stat.Size())
	if err := f.Close(); err != nil {
		return err
	}
	os.Remove(f.Name())
	return nil
}
