package mapDatalayer

import (
  "context"
  "fmt"
  "os"
  "strconv"

  "github.com/google/trillian"
  "github.com/google/trillian/client"
  "github.com/google/trillian/types"
  "google.golang.org/grpc"

  _ "github.com/google/trillian/merkle/coniks"
)

func getAdminClient() (trillian.TrillianAdminClient, *grpc.ClientConn, error) {
  MAP_ADDRESS := os.Getenv("MAP_ADDRESS")
  g, dialErr := grpc.Dial(MAP_ADDRESS, grpc.WithInsecure()) // TODO(z-tech): secure this
  if dialErr != nil {
    return nil, nil, dialErr
  }
  a := trillian.NewTrillianAdminClient(g)
  return a, g, nil
}

func getTrillianClient() (trillian.TrillianMapClient, *grpc.ClientConn, error) {
  MAP_ADDRESS := os.Getenv("MAP_ADDRESS")
  g, dialErr := grpc.Dial(MAP_ADDRESS, grpc.WithInsecure()) // TODO(z-tech): secure this
  if dialErr != nil {
    return nil, nil, dialErr
  }
  tc := trillian.NewTrillianMapClient(g)
  return tc, g, nil
}

func getWriterClient(tc trillian.TrillianMapClient, tree *trillian.Tree) (*client.MapClient, error) {
  wc, newClientFromTreeErr := client.NewMapClientFromTree(tc, tree)
  if newClientFromTreeErr != nil {
    return nil, newClientFromTreeErr
  }
  return wc, nil
}

func getRoot(c context.Context, tc trillian.TrillianMapClient) (*types.MapRootV1, error) {
  MAP_ID, _ := strconv.ParseInt(os.Getenv("MAP_ID"), 10, 64)
  rootRequest := &trillian.GetSignedMapRootRequest{MapId: MAP_ID}
  rootResponse, reqErr := tc.GetSignedMapRoot(c, rootRequest)
  if reqErr != nil {
    return nil, reqErr
  }
  var root types.MapRootV1
  unmarshalErr := root.UnmarshalBinary(rootResponse.MapRoot.MapRoot)
  if (unmarshalErr != nil) {
    return nil, unmarshalErr
  }
  return &root, nil
}

func AddLeaf(ctx context.Context, key []byte, data []byte) error {
  // 1) initialize some stuff
  MAP_ID, _ := strconv.ParseInt(os.Getenv("MAP_ID"), 10, 64)
  tc, g, getMapClientErr := getTrillianClient()
  if getMapClientErr != nil {
    fmt.Printf("error: getTrillianClient() %+v\n", getMapClientErr)
  }
  defer g.Close()

  // 2) get tree root
  root, getRootErr := getRoot(ctx, tc)
  if getRootErr != nil {
    fmt.Printf("error: failed to get tree root %d: %v\n", MAP_ID, getRootErr)
  }
  revision := int64(root.Revision)

  // 3) write revision to map
  l := trillian.MapLeaf{
    Index:     key,
    LeafValue: data,
  }
  writeReq := trillian.SetMapLeavesRequest{
    MapId:  MAP_ID,
    Leaves: []*trillian.MapLeaf{&l},
    Revision: int64(revision + 1),
  }
  writeResponse, writeErr := tc.SetLeaves(ctx, &writeReq)
  if writeErr != nil {
    fmt.Printf("error: failed to write revision to tree %d: %v\n", MAP_ID, writeErr)
    return writeErr
  }

  fmt.Printf("WriteResponse %+v\n", writeResponse)
  return nil
}

func GetLeaf(ctx context.Context, key []byte) (*trillian.MapLeaf, *types.MapRootV1, error) {
  // 1) initialize some stuff
  MAP_ID, _ := strconv.ParseInt(os.Getenv("MAP_ID"), 10, 64)
  tc, g1, getMapClientErr := getTrillianClient()
  if getMapClientErr != nil {
    fmt.Printf("error: getTrillianClient() %+v\n", getMapClientErr)
  }
  defer g1.Close()
  ac, g2, getAdminClientErr := getAdminClient()
  if getAdminClientErr != nil {
    fmt.Printf("error: getAdminClient() %+v\n", getAdminClientErr)
  }
  defer g2.Close()

  // 2) get tree
  tree, getTreeErr := ac.GetTree(ctx, &trillian.GetTreeRequest{TreeId: MAP_ID})
  if getTreeErr != nil {
    fmt.Printf("error: failed to get tree %d: %v\n", MAP_ID, getTreeErr)
  }

  // 3) get client
  wc, getWriterErr := getWriterClient(tc, tree)
  if getWriterErr != nil {
    fmt.Printf("error: failed to get writer client %d: %v\n", MAP_ID, getWriterErr)
  }

  // 4)
  indexes := [][]byte{key}
  mapLeaf, mapRoot, verifyErr := wc.GetAndVerifyMapLeaves(ctx, indexes)
  if verifyErr != nil {
    fmt.Printf("error: failed to verify map inclusion %d: %v\n", MAP_ID, verifyErr)
  }

  return mapLeaf[0], mapRoot, verifyErr
}
