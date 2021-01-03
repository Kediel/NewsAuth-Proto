package mapDatalayer

import (
  "context"
  "fmt"
  "os"
  "strconv"

  "github.com/google/trillian"
  "github.com/google/trillian/client"
  "github.com/google/trillian/types"
  _ "github.com/google/trillian/merkle/coniks"

  "github.com/z-tech/blue/src/datalayers/helpers"
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

func AddLeaf(ctx context.Context, key []byte, data []byte) error {
  // 1) initialize some stuff
  MAP_ADDRESS := os.Getenv("MAP_ADDRESS")
  MAP_ID, strconvErr := strconv.ParseInt(os.Getenv("MAP_ID"), 10, 64)
  if strconvErr != nil {
    fmt.Printf("error: unable read map id from environment %+v\n", strconvErr)
  }

  // 2) dial grpc connection
  grpcClientConn, getGRPCClientConnErr := datalayerHelpers.GetGRPCConn(MAP_ADDRESS)
  if getGRPCClientConnErr != nil {
    fmt.Printf("error: failed to dial grpcClient in map datalayer %+v\n", getGRPCClientConnErr)
  }
  defer grpcClientConn.Close()

  // 2) get the tree root
  trillianClient := trillian.NewTrillianMapClient(grpcClientConn)
  mapRoot, getMapRootErr := getMapRoot(ctx, MAP_ID, trillianClient)
  if getMapRootErr != nil {
    fmt.Printf("error: failed to get map root %d: %v\n", MAP_ID, getMapRootErr)
  }
  revisionNum := int64(mapRoot.Revision)

  // 3) write revision to map
  mapLeaf := trillian.MapLeaf{Index: key, LeafValue: data}
  writeReq := trillian.SetMapLeavesRequest{MapId: MAP_ID, Leaves: []*trillian.MapLeaf{&mapLeaf}, Revision: revisionNum + 1}
  writeResponse, writeErr := trillianClient.SetLeaves(ctx, &writeReq)
  if writeErr != nil {
    fmt.Printf("error: failed to write map revision %d: %v\n", MAP_ID, writeErr)
    return writeErr
  }

  GetLeaf(ctx, key)
  fmt.Printf("WriteResponse %+v\n", writeResponse)
  return nil
}

func GetLeaf(ctx context.Context, key []byte) (*trillian.MapLeaf, *types.MapRootV1, error) {
  // 1) initialize some stuff
  MAP_ADDRESS := os.Getenv("MAP_ADDRESS")
  MAP_ID, strconvErr := strconv.ParseInt(os.Getenv("MAP_ID"), 10, 64)
  if strconvErr != nil {
    fmt.Printf("error: unable read map id from environment %+v\n", strconvErr)
  }

  // 2) dial grpc connection
  grpcClientConn, getGRPCClientConnErr := datalayerHelpers.GetGRPCConn(MAP_ADDRESS)
  if getGRPCClientConnErr != nil {
    fmt.Printf("error: failed to dial grpcClient in map datalayer %+v\n", getGRPCClientConnErr)
  }
  defer grpcClientConn.Close()

  // 3) get the map tree
  adminClient := trillian.NewTrillianAdminClient(grpcClientConn)
  tree, getTreeErr := adminClient.GetTree(ctx, &trillian.GetTreeRequest{TreeId: MAP_ID})
  if getTreeErr != nil {
    fmt.Printf("error: failed to get tree in map datalayer %d: %v\n", MAP_ID, getTreeErr)
  }

  // 4) get and verify the leaf
  trillianClient := trillian.NewTrillianMapClient(grpcClientConn)
  mapClient, getMapClientErr := client.NewMapClientFromTree(trillianClient, tree)
  if getMapClientErr != nil {
    fmt.Printf("error: failed to get map client %d: %v\n", MAP_ID, getTreeErr)
  }
  indexes := [][]byte{key}
  mapLeaf, mapRoot, getAndVerifyErr := mapClient.GetAndVerifyMapLeaves(ctx, indexes)
  if getAndVerifyErr != nil {
    fmt.Printf("error: failed to verify map inclusion %d: %v\n", MAP_ID, getAndVerifyErr)
  }
  fmt.Printf("mapLeaf, mapRoot %+v %+v\n", mapLeaf, mapRoot)
  return mapLeaf[0], mapRoot, nil
}
