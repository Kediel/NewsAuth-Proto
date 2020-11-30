package log

import (
	"context"
	"encoding/json"
	"log"
  "os"

  "github.com/google/trillian"
  "google.golang.org/grpc"
  "google.golang.org/grpc/codes"
)

type Leaf struct {
	Entry map[string]interface{}
	Hash  string
	Item  map[string]interface{}
}

func AddLeaf(leaf Leaf) error {
  ctx := context.Background() // TODO(z-tech): what's the deal with this?
  g, err := grpc.Dial(*trillianLog, grpc.WithInsecure()) // TODO(z-tech): secure this. use singleton?
  if err != nil {
    // log.Fatalf("Failed to dial Trillian Log: %v", err)
  }
  defer g.Close()

  tc := trillian.NewTrillianLogClient(g)

  j, err := json.Marshal(leaf)
  if err != nil {
    // return err
  }

  // Send to Trillian
  tl := &trillian.LogLeaf{LeafValue: j}
  q := &trillian.QueueLeafRequest{LogId: os.Getenv("LOG_ID"), Leaf: tl}
  r, err := tc.QueueLeaf(ctx, q)
  if err != nil {
    // return err
  }

  // And check everything worked
  c := codes.Code(r.QueuedLeaf.GetStatus().GetCode())
  if c != codes.OK && c != codes.AlreadyExists {
    // return fmt.Errorf("bad return status: %v", r.QueuedLeaf.GetStatus())
  }

  if c == codes.AlreadyExists {
    // return err
  }
}
