package logDatalayer

import (
  "context"
  "fmt"
  "os"
  "strconv"

  "google.golang.org/grpc"
  "google.golang.org/grpc/codes"
  "github.com/google/trillian"
  "github.com/google/trillian/client"
  "github.com/google/trillian/merkle"
  "github.com/google/trillian/merkle/rfc6962"
  "github.com/google/trillian/types"
)

func getAdminClient() (trillian.TrillianAdminClient, *grpc.ClientConn, error) {
  LOG_ADDRESS := os.Getenv("LOG_ADDRESS")
  g, dialErr := grpc.Dial(LOG_ADDRESS, grpc.WithInsecure()) // TODO(z-tech): secure this
  if dialErr != nil {
    return nil, nil, dialErr
  }
  a := trillian.NewTrillianAdminClient(g)
  return a, g, nil
}

func getTrillianClient() (trillian.TrillianLogClient, *grpc.ClientConn, error) {
  LOG_ADDRESS := os.Getenv("LOG_ADDRESS")
  g, dialErr := grpc.Dial(LOG_ADDRESS, grpc.WithInsecure()) // TODO(z-tech): secure this
  if dialErr != nil {
    return nil, nil, dialErr
  }
  tc := trillian.NewTrillianLogClient(g)
  return tc, g, nil
}

func getWriterClient(tc trillian.TrillianLogClient, tree *trillian.Tree, root types.LogRootV1) (*client.LogClient, error) {
  wc, newClientFromTreeErr := client.NewFromTree(tc, tree, root)
  if newClientFromTreeErr != nil {
    return nil, newClientFromTreeErr
  }
  return wc, nil
}

func getRoot(c context.Context, tc trillian.TrillianLogClient) (*types.LogRootV1, error) {
  LOG_ID, _ := strconv.ParseInt(os.Getenv("LOG_ID"), 10, 64)
  rootRequest := &trillian.GetLatestSignedLogRootRequest{LogId: LOG_ID}
  rootResponse, reqErr := tc.GetLatestSignedLogRoot(c, rootRequest)
  if reqErr != nil {
    return nil, reqErr
  }

  var root types.LogRootV1
  unmarshalErr := root.UnmarshalBinary(rootResponse.SignedLogRoot.LogRoot)
  if (unmarshalErr != nil) {
    return nil, unmarshalErr
  }
  return &root, nil
}

func AddLeaf(ctx context.Context, data []byte) (int64, int64, [][]byte, []byte, []byte, bool, error) {
  // 1) initialize some stuff
  LOG_ID, _ := strconv.ParseInt(os.Getenv("LOG_ID"), 10, 64)
  tc, g1, getLogClientErr := getTrillianClient()
  if getLogClientErr != nil {
    fmt.Printf("error: getTrillianClient() %+v\n", getLogClientErr)
    return 0, 0, nil, nil, nil, false, getLogClientErr
  }
  defer g1.Close()
  ac, g2, getAdminClientErr := getAdminClient()
  if getAdminClientErr != nil {
    fmt.Printf("error: getAdminClient() %+v\n", getAdminClientErr)
    return 0, 0, nil, nil, nil, false, getAdminClientErr
  }
  defer g2.Close()

  // 2) get tree
  tree, getTreeErr := ac.GetTree(ctx, &trillian.GetTreeRequest{TreeId: LOG_ID})
  if getTreeErr != nil {
    fmt.Printf("error: failed to get tree %d: %v\n", LOG_ID, getTreeErr)
    return 0, 0, nil, nil, nil, false, getTreeErr
  }

  // 3) get tree root
  root, getRootErr := getRoot(ctx, tc)
  if getRootErr != nil {
    fmt.Printf("error: failed to get tree root %d: %v\n", LOG_ID, getRootErr)
    return 0, 0, nil, nil, nil, false, getRootErr
  }

  // 4) get client (the kind that does QueueLeaves(), naming it "writer client")
  wc, getWriterErr := getWriterClient(tc, tree, *root)
  if getWriterErr != nil {
    fmt.Printf("error: failed to get writer client %d: %v\n", LOG_ID, getWriterErr)
    return 0, 0, nil, nil, nil, false, getWriterErr
  }

  // 5) Queue the leaf
  leaf := wc.BuildLeaf(data)
  queueLeafResp, queueLeafErr := tc.QueueLeaf(ctx, &trillian.QueueLeafRequest{
    LogId: wc.LogID,
    Leaf:  leaf,
  })
  if queueLeafErr != nil {
    fmt.Printf("error: failed to queue leaf %d: %v\n", LOG_ID, queueLeafErr)
    return 0, 0, nil, nil, nil, false, queueLeafErr
  }

  // 6) wait for inclusion
  waitErr := wc.WaitForInclusion(ctx, data)
  if waitErr != nil {
    fmt.Printf("error: failed to wait for leaf inclusion %d: %v\n", LOG_ID, waitErr)
    return 0, 0, nil, nil, nil, false, waitErr
  }

  // 7) Check if dup
  isDup := false
  if queueLeafResp.QueuedLeaf.Status != nil { // not sure why status missing for new leaves
    respCode := codes.Code(queueLeafResp.QueuedLeaf.Status.Code)
    if respCode != codes.OK && respCode != codes.AlreadyExists {
      fmt.Printf("error: queue leaf status is unsuccessful %d %v\n", LOG_ID, respCode)
      // TODO: return nil, nil, nil, nil, nil, nil, getNewRootErr
    } else if (respCode != codes.OK && respCode == codes.AlreadyExists) {
      isDup = true
    }
  }

  // 8) get the new tree root
  newRoot, getNewRootErr := getRoot(ctx, tc)
  if getNewRootErr != nil {
    fmt.Printf("error: failed to get new tree root %d: %v\n", LOG_ID, getNewRootErr)
    return 0, 0, nil, nil, nil, false, getNewRootErr
  }

  // 9) Get the inclusion proof from hash
  getProofResp, getProofErr := tc.GetInclusionProofByHash(ctx,
    &trillian.GetInclusionProofByHashRequest{
      LogId:    LOG_ID,
      LeafHash: leaf.MerkleLeafHash,
      TreeSize: int64(newRoot.TreeSize),
    })
  if getProofErr != nil {
    fmt.Printf("error: failed to get new tree root %d: %v\n", LOG_ID, getProofErr)
    return 0, 0, nil, nil, nil, false, getProofErr
  }
  leafIndex := getProofResp.Proof[0].LeafIndex
  proof := getProofResp.Proof[0].Hashes
  treeSize := int64(newRoot.TreeSize)
  rootHash := newRoot.RootHash
  leafHash := leaf.MerkleLeafHash

  // 10) verify the inclusion proof
  verifier := merkle.NewLogVerifier(rfc6962.DefaultHasher)
  verifyErr := verifier.VerifyInclusionProof(leafIndex, treeSize, proof, rootHash, leafHash)
  if verifyErr != nil {
    fmt.Printf("error: failed to verify inclusion %d: %+v\n", LOG_ID, verifyErr)
    return 0, 0, nil, nil, nil, false, verifyErr
  }

  return leafIndex, treeSize, proof, rootHash, leafHash, isDup, nil
}
