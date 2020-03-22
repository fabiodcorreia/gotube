package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"runtime/trace"

	"github.com/fabiodcorreia/gotube"
)

func main() {
	ft, err := os.Create("./trace.out")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	trace.Start(ft)
	defer trace.Stop()
	defer ft.Close()

	fc, err := os.Create("./cpu.pprof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	defer fc.Close()
	if err := pprof.StartCPUProfile(fc); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()

	target := "https://www.youtube.com/watch?v=urarTyKn9cg2"

	v, err := gotube.GetVideoDetails(target)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(v.Streams[0].ContentLength)

	file, err := os.Create("./" + v.Title + string(v.Streams[0].Extension))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	bc2, err := v.Download(file)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(bc2)

	fm, err := os.Create("./mem.pprof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer fm.Close()
	runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(fm); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}

}
