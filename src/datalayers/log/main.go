package logDatalayer

import (
  "context"
  "fmt"
  "os"
  "strconv"

  "google.golang.org/grpc/codes"
  "github.com/google/trillian"
  "github.com/google/trillian/client"
  "github.com/google/trillian/merkle"
  "github.com/google/trillian/merkle/rfc6962"
  "github.com/google/trillian/types"

  "github.com/z-tech/blue/src/datalayers/helpers"
)

func getLogRoot(ctx context.Context, logID int64, trillianClient trillian.TrillianLogClient) (*types.LogRootV1, error) {
  logRootRequest := &trillian.GetLatestSignedLogRootRequest{LogId: logID}
  logRootResponse, logRootRequestErr := trillianClient.GetLatestSignedLogRoot(ctx, logRootRequest)
  if logRootRequestErr != nil {
    return nil, logRootRequestErr
  }
  var logRoot types.LogRootV1
  unmarshalErr := logRoot.UnmarshalBinary(logRootResponse.SignedLogRoot.LogRoot)
  if (unmarshalErr != nil) {
    return nil, unmarshalErr
  }
  return &logRoot, nil
}

func AddLeaf(ctx context.Context, data []byte) (int64, int64, [][]byte, []byte, []byte, bool, error) {
  // 1) initialize some stuff
  LOG_ADDRESS := os.Getenv("LOG_ADDRESS")
  LOG_ID, strconvErr := strconv.ParseInt(os.Getenv("LOG_ID"), 10, 64)
  if strconvErr != nil {
    fmt.Printf("error: unable read log id from environment %+v\n", strconvErr)
  }

  // 2) dial grpc connection
  grpcClientConn, getGRPCClientConnErr := datalayerHelpers.GetGRPCConn(LOG_ADDRESS)
  if getGRPCClientConnErr != nil {
    fmt.Printf("error: failed to dial grpcClient in map datalayer %+v\n", getGRPCClientConnErr)
  }
  defer grpcClientConn.Close()

  // 3) get tree
  adminClient := trillian.NewTrillianAdminClient(grpcClientConn)
  tree, getTreeErr := adminClient.GetTree(ctx, &trillian.GetTreeRequest{TreeId: LOG_ID})
  if getTreeErr != nil {
    fmt.Printf("error: failed to get log tree %d: %v\n", LOG_ID, getTreeErr)
    return 0, 0, nil, nil, nil, false, getTreeErr
  }

  // 4) get log root
  trillianClient := trillian.NewTrillianLogClient(grpcClientConn)
  logRoot, getLogRootErr := getLogRoot(ctx, LOG_ID, trillianClient)
  if getLogRootErr != nil {
    fmt.Printf("error: failed to get tree root %d: %v\n", LOG_ID, getLogRootErr)
    return 0, 0, nil, nil, nil, false, getLogRootErr
  }

  // 5) queue the leaf
  logClient, getLogClientErr := client.NewFromTree(trillianClient, tree, *logRoot)
  if getLogClientErr != nil {
    fmt.Printf("error: failed to get log client %d: %v\n", LOG_ID, getLogClientErr)
    return 0, 0, nil, nil, nil, false, getLogClientErr
  }
  logLeaf := logClient.BuildLeaf(data)
  queueLeafResp, queueLeafErr := trillianClient.QueueLeaf(ctx, &trillian.QueueLeafRequest{LogId: LOG_ID, Leaf: logLeaf})
  if queueLeafErr != nil {
    fmt.Printf("error: failed to queue leaf %d: %v\n", LOG_ID, queueLeafErr)
    return 0, 0, nil, nil, nil, false, queueLeafErr
  }

  // 6) wait for inclusion
  inclusionErr := logClient.WaitForInclusion(ctx, data)
  if inclusionErr != nil {
    fmt.Printf("error: failed to wait for leaf inclusion %d: %v\n", LOG_ID, inclusionErr)
    return 0, 0, nil, nil, nil, false, inclusionErr
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
  newLogRoot, getNewLogRootErr := getLogRoot(ctx, LOG_ID, trillianClient)
  if getNewLogRootErr != nil {
    fmt.Printf("error: failed to get new tree root %d: %v\n", LOG_ID, getNewLogRootErr)
    return 0, 0, nil, nil, nil, false, getNewLogRootErr
  }

  // 9) Get the inclusion proof from hash
  getProofResp, getProofErr := trillianClient.GetInclusionProofByHash(ctx,
    &trillian.GetInclusionProofByHashRequest{
      LogId:    LOG_ID,
      LeafHash: logLeaf.MerkleLeafHash,
      TreeSize: int64(newLogRoot.TreeSize),
    })
  if getProofErr != nil {
    fmt.Printf("error: failed to get new tree root %d: %v\n", LOG_ID, getProofErr)
    return 0, 0, nil, nil, nil, false, getProofErr
  }
  leafIndex := getProofResp.Proof[0].LeafIndex
  proof := getProofResp.Proof[0].Hashes
  if leafIndex == 0 {
    proof = make([][]byte, 0)
  }
  treeSize := int64(newLogRoot.TreeSize)
  rootHash := newLogRoot.RootHash
  leafHash := logLeaf.MerkleLeafHash

  // 10) verify the inclusion proof
  verifier := merkle.NewLogVerifier(rfc6962.DefaultHasher)
  verifyErr := verifier.VerifyInclusionProof(leafIndex, treeSize, proof, rootHash, leafHash)
  if verifyErr != nil {
    fmt.Printf("error: failed to verify inclusion %d: %+v\n", LOG_ID, verifyErr)
    return 0, 0, nil, nil, nil, false, verifyErr
  }

  return leafIndex, treeSize, proof, rootHash, leafHash, isDup, nil
}
