gen-protoc-ts:
    ls proto/tollroad | xargs -I {} protoc \
        --plugin="./scripts/node_modules/.bin/protoc-gen-ts_proto" \
        --ts_proto_out="./client/src/types/generated" \
        --proto_path="./proto" \
        --ts_proto_opt="esModuleInterop=true,forceLong=long,useOptionals=messages" \
        tollroad/{}
