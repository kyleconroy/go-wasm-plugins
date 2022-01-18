# https://github.com/bytecodealliance/wasmtime-go/pull/102
 
```
cd wasmtime
cargo build -p wasmtime-c-api --release
```

cp the release into the correct spot 

https://awstip.com/containerize-go-sqlite-with-docker-6d7fbecd14f0

```
RUN CGO_ENABLED=1 GOOS=linux go build -o /app -a -ldflags '-linkmode external -extldflags "-static"' .
```

cd ../wasmtime-go
./ci/local.sh ../wasmtime

https://stackoverflow.com/questions/23879205/how-to-change-lib-path-for-go-build
-ldflag "-L/usr/lib -lncursesw
go get -v -ldflags "-L/usr/lib" code.google.com/p/goncurses

CGO_CFLAGS="-I/Users/kyle/projects/go-wasm-plugins/wasmtime/crates/c-api/wasm-c-api/include" CGO_LDFLAGS="-L/Users/kyle/projects/go-wasm-plugins/wasmtime/target/release/" go build ./...

I need to fork wasmtime-go, create the correct releases and upload them with the correct builds.

