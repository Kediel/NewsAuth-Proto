package grpcDatalayer

import (
  "context"
  "fmt"

  "github.com/google/trillian"
  "github.com/google/trillian/client"
  "github.com/google/trillian/merkle"
  "github.com/google/trillian/merkle/rfc6962"
  "github.com/google/trillian/types"
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

func AddLogLeaf(ctx context.Context, logAddress string, logID int64, data []byte) error {
  // 1) dial grpc connection
  grpcClientConn, getGRPCClientConnErr := GetGRPCConn(logAddress)
  if getGRPCClientConnErr != nil {
    fmt.Printf("error: failed to dial grpcClient in map datalayer %d: %v\n", logID, getGRPCClientConnErr)
    return getGRPCClientConnErr
  }
  defer grpcClientConn.Close()

  // 2) get tree
  adminClient := trillian.NewTrillianAdminClient(grpcClientConn)
  tree, getTreeErr := adminClient.GetTree(ctx, &trillian.GetTreeRequest{TreeId: logID})
  if getTreeErr != nil {
    fmt.Printf("error: failed to get log tree %d: %v\n", logID, getTreeErr)
    return getTreeErr
  }

  // 3) get log root
  trillianClient := trillian.NewTrillianLogClient(grpcClientConn)
  logRoot, getLogRootErr := getLogRoot(ctx, logID, trillianClient)
  if getLogRootErr != nil {
    fmt.Printf("error: failed to get tree root %d: %v\n", logID, getLogRootErr)
    return getLogRootErr
  }

  // 4) queue the leaf
  logClient, getLogClientErr := client.NewFromTree(trillianClient, tree, *logRoot)
  if getLogClientErr != nil {
    fmt.Printf("error: failed to get log client %d: %v\n", logID, getLogClientErr)
    return getLogClientErr
  }
  logLeaf := logClient.BuildLeaf(data)
  _, queueLeafErr := trillianClient.QueueLeaf(ctx, &trillian.QueueLeafRequest{LogId: logID, Leaf: logLeaf})
  if queueLeafErr != nil {
    fmt.Printf("error: failed to queue leaf %d: %v\n", logID, queueLeafErr)
    return queueLeafErr
  }

  // 6) wait for inclusion
  inclusionErr := logClient.WaitForInclusion(ctx, data)
  if inclusionErr != nil {
    fmt.Printf("error: failed to wait for leaf inclusion %d: %v\n", logID, inclusionErr)
    return inclusionErr
  }

  return nil
}

func GetLogLeaf(ctx context.Context, logAddress string, logID int64, data []byte) (int64, int64, [][]byte, []byte, []byte, error) {
  // 1) dial grpc connection
  grpcClientConn, getGRPCClientConnErr := GetGRPCConn(logAddress)
  if getGRPCClientConnErr != nil {
    fmt.Printf("error: failed to dial grpcClient in map datalayer %+v\n", getGRPCClientConnErr)
  }
  defer grpcClientConn.Close()

  // 2) get tree
  adminClient := trillian.NewTrillianAdminClient(grpcClientConn)
  tree, getTreeErr := adminClient.GetTree(ctx, &trillian.GetTreeRequest{TreeId: logID})
  if getTreeErr != nil {
    fmt.Printf("error: failed to get log tree %d: %v\n", logID, getTreeErr)
    return 0, 0, nil, nil, nil, getTreeErr
  }

  // 3) get log root
  trillianClient := trillian.NewTrillianLogClient(grpcClientConn)
  logRoot, getLogRootErr := getLogRoot(ctx, logID, trillianClient)
  if getLogRootErr != nil {
    fmt.Printf("error: failed to get tree root %d: %v\n", logID, getLogRootErr)
    return 0, 0, nil, nil, nil, getLogRootErr
  }
  treeSize := int64(logRoot.TreeSize)
  rootHash := logRoot.RootHash

  // 4) build leaf
  logClient, getLogClientErr := client.NewFromTree(trillianClient, tree, *logRoot)
  if getLogClientErr != nil {
    fmt.Printf("error: failed to get log client %d: %v\n", logID, getLogClientErr)
    return 0, 0, nil, nil, nil, getLogClientErr
  }
  logLeaf := logClient.BuildLeaf(data)
  leafHash := logLeaf.MerkleLeafHash

  // 5) get inclusion proof from the leaf, return if not exists
  getProofResp, getProofErr := trillianClient.GetInclusionProofByHash(ctx,
    &trillian.GetInclusionProofByHashRequest{
      LogId:    logID,
      LeafHash: logLeaf.MerkleLeafHash,
      TreeSize: treeSize,
    })
  if getProofErr != nil {
    // NOT SO MUCH AN ERROR: there is no proof for this hash, meaning it's not in the log
    return -1, treeSize, nil, rootHash, leafHash, nil
  }
  leafIndex := getProofResp.Proof[0].LeafIndex
  proof := getProofResp.Proof[0].Hashes
  if leafIndex == 0 {
    proof = make([][]byte, 0)
  }

  // 6) verify the inclusion proof
  verifier := merkle.NewLogVerifier(rfc6962.DefaultHasher)
  verifyErr := verifier.VerifyInclusionProof(leafIndex, treeSize, proof, rootHash, leafHash)
  if verifyErr != nil {
    fmt.Printf("error: failed to verify inclusion %d: %+v\n", logID, verifyErr)
    return 0, 0, nil, nil, nil, verifyErr
  }

  // 7) done, yay
  return leafIndex, treeSize, proof, rootHash, leafHash, nil
}
