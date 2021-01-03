package datalayerHelpers

import (
  "google.golang.org/grpc"
)

func GetGRPCConn(address string) (*grpc.ClientConn, error) {
  clientConn, dialErr := grpc.Dial(address, grpc.WithInsecure()) // TODO(z-tech): secure this
  if dialErr != nil {
    return nil, dialErr
  }
  return clientConn, nil
}
