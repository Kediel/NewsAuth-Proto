package datalayers

import (
  "os"
  "fmt"
  "strconv"

  "google.golang.org/grpc"
)

func GetConfig() (string, int64, string, int64, error) {
  LOG_ADDRESS := os.Getenv("LOG_ADDRESS")
  LOG_ID, logIDConvErr := strconv.ParseInt(os.Getenv("LOG_ID"), 10, 64)
  if logIDConvErr != nil {
    fmt.Printf("error: unable read log id from environment %+v\n", logIDConvErr)
    return "", 0, "", 0, logIDConvErr
  }
  MAP_ADDRESS := os.Getenv("MAP_ADDRESS")
  MAP_ID, mapIDConvErr := strconv.ParseInt(os.Getenv("MAP_ID"), 10, 64)
  if mapIDConvErr != nil {
    fmt.Printf("error: unable read map id from environment %+v\n", mapIDConvErr)
    return "", 0, "", 0, mapIDConvErr
  }
  return LOG_ADDRESS, LOG_ID, MAP_ADDRESS, MAP_ID, nil
}

func GetGRPCConn(address string) (*grpc.ClientConn, error) {
  clientConn, dialErr := grpc.Dial(address, grpc.WithInsecure()) // TODO(z-tech): secure this
  if dialErr != nil {
    return nil, dialErr
  }
  return clientConn, nil
}
