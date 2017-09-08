# Cisco Spark Command Line Client

## How To Use
### Messages
```
$ spark message "hello world"             # send message to default space
$ spark message -s <space> "hello world"  # send message to specific space
$ spark message -p <person> "hello world" # send message to specific person
$ spark message -f <filename>             # send file attachment to default space
$ spark message -f <filename> -s <space>  # send file attachment to specific space
$ spark message -f <filename> -s <person> # send file attachment to specific person
```

## How To Build
```
$ go build -o spark
```
