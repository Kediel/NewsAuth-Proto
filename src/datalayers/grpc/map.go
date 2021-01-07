package grpcDatalayer

import (
  "context"
  "fmt"

  "github.com/google/trillian"
  "github.com/google/trillian/client"
  "github.com/google/trillian/types"
  _ "github.com/google/trillian/merkle/coniks"
)

func getMapRoot(ctx context.Context, mapID int64, trillianClient trillian.TrillianMapClient) (*types.MapRootV1, error) {
  mapRootRequest := &trillian.GetSignedMapRootRequest{MapId: mapID}
  mapRootResponse, mapRootRequestErr := trillianClient.GetSignedMapRoot(ctx, mapRootRequest)
  if mapRootRequestErr != nil {
    return nil, mapRootRequestErr
  }
  var mapRoot types.MapRootV1
  unmarshalErr := mapRoot.UnmarshalBinary(mapRootResponse.MapRoot.MapRoot)
  if (unmarshalErr != nil) {
    return nil, unmarshalErr
  }
  return &mapRoot, nil
}

func AddMapLeaf(ctx context.Context, mapAddress string, mapID int64, key []byte, data []byte) error {
  // 1) dial grpc connection
  grpcClientConn, getGRPCClientConnErr := GetGRPCConn(mapAddress)
  if getGRPCClientConnErr != nil {
    fmt.Printf("error: failed to dial grpcClient in map datalayer %+v\n", getGRPCClientConnErr)
  }
  defer grpcClientConn.Close()

  // 2) get the tree root
  trillianClient := trillian.NewTrillianMapClient(grpcClientConn)
  mapRoot, getMapRootErr := getMapRoot(ctx, mapID, trillianClient)
  if getMapRootErr != nil {
    fmt.Printf("error: failed to get map root %d: %v\n", mapID, getMapRootErr)
  }
  revisionNum := int64(mapRoot.Revision)

  // 3) write revision to map
  mapLeaf := trillian.MapLeaf{Index: key, LeafValue: data}
  writeReq := trillian.SetMapLeavesRequest{MapId: mapID, Leaves: []*trillian.MapLeaf{&mapLeaf}, Revision: revisionNum + 1}
  _, writeErr := trillianClient.SetLeaves(ctx, &writeReq)
  if writeErr != nil {
    fmt.Printf("error: failed to write map revision %d: %v\n", mapID, writeErr)
    return writeErr
  }

  return nil
}

func GetMapLeaf(ctx context.Context, mapAddress string, mapID int64, key []byte) (bool, []byte, []byte, [][]byte, error) {
  // 2) dial grpc connection
  grpcClientConn, getGRPCClientConnErr := GetGRPCConn(mapAddress)
  if getGRPCClientConnErr != nil {
    fmt.Printf("error: failed to dial grpcClient in map datalayer %+v\n", getGRPCClientConnErr)
  }
  defer grpcClientConn.Close()

  // 3) get the map tree
  adminClient := trillian.NewTrillianAdminClient(grpcClientConn)
  tree, getTreeErr := adminClient.GetTree(ctx, &trillian.GetTreeRequest{TreeId: mapID})
  if getTreeErr != nil {
    fmt.Printf("error: failed to get tree in map datalayer %d: %v\n", mapID, getTreeErr)
  }

  // 4) get and verify the leaf
  trillianClient := trillian.NewTrillianMapClient(grpcClientConn)
  mapClient, getMapClientErr := client.NewMapClientFromTree(trillianClient, tree)
  if getMapClientErr != nil {
    fmt.Printf("error: failed to get map client %d: %v\n", mapID, getTreeErr)
  }
  indexes := [][]byte{key}
  mapLeaf, _, getAndVerifyErr := mapClient.GetAndVerifyMapLeaves(ctx, indexes)
  if getAndVerifyErr != nil {
    fmt.Printf("error: failed to verify map inclusion %d: %v\n", mapID, getAndVerifyErr)
  }

  // 5) get proof
  getLeavesResp, _ := mapClient.Conn.GetLeaves(ctx, &trillian.GetMapLeavesRequest{
    MapId: mapID,
    Index: indexes,
  })
  proof := getLeavesResp.MapLeafInclusion[0].Inclusion
  isExists := true
  if mapLeaf[0].LeafValue == nil {
    isExists = false
  }

  return isExists, mapLeaf[0].LeafHash, mapLeaf[0].LeafValue, proof, nil
}
