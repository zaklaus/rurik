all:
	go build -o build/rurik.exe src/demo/*.go

wasm:
	GOOS=js GOARCH=wasm go build -o build/rurik.exe src/demo/*.go

perf:
	go tool pprof --pdf build/cpu.pprof > build/shit.pdf

clean:
	rm -rf build/*

play:
	./play.sh
