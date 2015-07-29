// Copyright (C) 2015 Space Monkey, Inc.

package sdl

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_mixer"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

var (
	execChan = make(chan func())
)

func Run(cont func()) {
	runtime.LockOSThread()
	defer panic("shouldn't exit")

	sdl.Init(sdl.INIT_EVERYTHING)
	ttf.Init()

	err := mix.OpenAudio(mix.DEFAULT_FREQUENCY, mix.DEFAULT_FORMAT, 2, 4096)
	if err != nil {
		fmt.Printf("Could not open audio: %v\n", err)
	}

	go cont()

	for {
		(<-execChan)()
	}
}

func exec(cb func() error) (err error) {
	var r interface{}
	var mtx sync.Mutex
	mtx.Lock()
	execChan <- func() {
		defer func() {
			r = recover()
			mtx.Unlock()
		}()
		err = cb()
	}
	mtx.Lock()
	mtx.Unlock()
	if r != nil {
		panic(r)
	}
	return err
}
