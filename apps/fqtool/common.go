package main

import "runtime"

func setThread(thread int) {
	if cpu := runtime.NumCPU(); thread < 1 || thread > cpu {
		thread = cpu
	}
	runtime.GOMAXPROCS(thread)
}
