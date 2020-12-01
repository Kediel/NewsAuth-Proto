package logDatalayer

import (
  "context"
  // "encoding/json"
  "fmt"
  // "strconv"

  "github.com/google/trillian"
  "github.com/google/trillian/client"
  "github.com/google/trillian/types"
  "google.golang.org/grpc"
  // "google.golang.org/grpc/codes"
)

const LOG_ID = int64(1067636684015883737)

func getTrillianClient() (trillian.TrillianLogClient, *grpc.ClientConn, error) {
  const LOG_ADDRESS = "ec2-3-91-133-44.compute-1.amazonaws.com:8090"
  g, dialErr := grpc.Dial(LOG_ADDRESS, grpc.WithInsecure()) // TODO(z-tech): secure this
  if dialErr != nil {
    return nil, nil, dialErr
  }
  tc := trillian.NewTrillianLogClient(g)
  return tc, g, nil
}

func getRoot(c context.Context, tc trillian.TrillianLogClient) (*types.LogRootV1, error) {
  rootRequest := &trillian.GetLatestSignedLogRootRequest{LogId: LOG_ID}
  rootResponse, reqErr := tc.GetLatestSignedLogRoot(c, rootRequest)
  if reqErr != nil {
    return nil, reqErr
  }

  var root types.LogRootV1 // TODO(z-tech): verify hash?
  unmarshalErr := root.UnmarshalBinary(rootResponse.SignedLogRoot.LogRoot)
  if (unmarshalErr != nil) {
    return nil, unmarshalErr
  }
  return &root, nil
}

func AddLeaf(data []byte) error {
  ctx := context.Background() // TODO(z-tech): what's the deal with this?

  tc, g, getLogClientErr := getTrillianClient()
  if getLogClientErr != nil {
    fmt.Printf("error: getRoot() %+v", getLogClientErr)
  }
  defer g.Close()

  root, getRootErr := getRoot(ctx, tc)
  if getRootErr != nil {
    fmt.Printf("error: getRoot() %+v", getRootErr)
  }

  admin := trillian.NewTrillianAdminClient(g)
  tree, getTreeErr := admin.GetTree(ctx, &trillian.GetTreeRequest{TreeId: LOG_ID})
  if getTreeErr != nil {
    fmt.Printf("failed to get tree %d: %v", LOG_ID, getTreeErr)
  }
  tree.HashStrategy = trillian.HashStrategy_RFC6962_SHA256
  fmt.Printf("Got tree %+v\n", tree)

  // logVerifier, newVerifierErr := client.NewLogVerifierFromTree(tree)
  // if newVerifierErr != nil {
  //   fmt.Printf("failed to get new verifier %+v\n", newVerifierErr)
  // }
  // fmt.Printf("debug: log verifier is %+v\n", logVerifier)

  logClient, newClientErr := client.NewFromTree(tc, tree, *root)
  if newClientErr != nil {
    fmt.Printf("failed to get new client %+v\n", newClientErr)
  }

  fmt.Printf("debug: logclient is %+v\n", logClient)

  return nil

  // rawLeaf, _ := json.Marshal(leaf)
  //
  // // Send to Trillian
  // tl := &trillian.LogLeaf{LeafValue: j}
  // q := &trillian.QueueLeafRequest{LogId: int64(LOG_ID), Leaf: tl}
  // r, _ := tc.QueueLeaf(ctx, q)
  //
  // // And check everything worked
  // c := codes.Code(r.QueuedLeaf.GetStatus().GetCode())
  // if c != codes.OK && c != codes.AlreadyExists {
  //   // return fmt.Errorf("bad return status: %v", r.QueuedLeaf.GetStatus())
  // }
  //
  // if c == codes.AlreadyExists {
  //   // return err
  // }
}
