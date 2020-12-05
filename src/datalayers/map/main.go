package mapDatalayer

import (
  "context"
  "fmt"

  "github.com/google/trillian"
  "github.com/google/trillian/client"
  "github.com/google/trillian/types"
  "google.golang.org/grpc"

  _ "github.com/google/trillian/merkle/coniks"
)

const MAP_ID = int64(8499339143476310100)
const MAP_ADDRESS = "ec2-3-91-133-44.compute-1.amazonaws.com:8093"

func getAdminClient() (trillian.TrillianAdminClient, *grpc.ClientConn, error) {
  g, dialErr := grpc.Dial(MAP_ADDRESS, grpc.WithInsecure()) // TODO(z-tech): secure this
  if dialErr != nil {
    return nil, nil, dialErr
  }
  a := trillian.NewTrillianAdminClient(g)
  return a, g, nil
}

func getTrillianClient() (trillian.TrillianMapClient, *grpc.ClientConn, error) {
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

func AddLeaf(ctx context.Context, key string, data []byte) {
  // 1) initialize some stuff
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

  // 3) get tree root
  root, getRootErr := getRoot(ctx, tc)
  if getRootErr != nil {
    fmt.Printf("error: failed to get tree root %d: %v\n", MAP_ID, getRootErr)
  }
  revision := int64(root.Revision)

  // 4) get client (the kind that does QueueLeaves(), naming it "writer client")
  wc, getWriterErr := getWriterClient(tc, tree)
  if getWriterErr != nil {
    fmt.Printf("error: failed to get writer client %d: %v\n", MAP_ID, getWriterErr)
  }

  fmt.Printf("HELLLO %+v %v\n", wc, revision)
}
