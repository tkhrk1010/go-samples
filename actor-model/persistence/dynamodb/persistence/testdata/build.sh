protoc --go_out=. --go_opt=paths=source_relative \
    -I../../ -I. test_event.proto
