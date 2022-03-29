We use a custom buf Docker image in order to accommodate compiling our protobufs
to Typescript. 

### Rebuild our custom Buf image
You shouldn't have to do this often. Only if an internal NPM dependency needs
updating:
```
docker build \
  --no-cache \
  --platform linux/amd64 \
  --file=Dockerfile-buf \
  --tag=buftest:latest .
```

We need to specify the platform because grpc-tools aren't available for other 
architectures, [sadly](https://github.com/grpc/grpc-node/issues/1405#issuecomment-623677869).

### Generate Typescript protos
```
docker run \
  --platform linux/amd64 \
  -it \
  -v $(pwd):/workspace \
  buftest generate \
    --template /workspace/buf/buf.gen.ts.yaml \
    .
```