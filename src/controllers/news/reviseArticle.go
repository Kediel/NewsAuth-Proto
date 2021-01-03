package newsController

import (
  b64 "encoding/base64"
  "encoding/json"
  "fmt"
  "net/http"

  "github.com/gin-gonic/gin"
  "github.com/gin-gonic/gin/binding"

  // "github.com/z-tech/blue/src/datalayers/log"
  "github.com/z-tech/blue/src/datalayers/map"
)

type ReviseArticleData struct {
  ArticleKey string `json:"ArticleKey"`
}

func ValidateReviseArticle(ctx *gin.Context) {
  ValidatePublishArticle(ctx)
  // reviseArticleData := ReviseArticleData{}
  // bindErr := ctx.ShouldBindBodyWith(&reviseArticleData, binding.JSON)
  // if bindErr != nil {
  //   ctx.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
  //   ctx.Abort()
  //   return
  // }
  // articleKey, decodeArticleKeyErr := b64.URLEncoding.DecodeString(reviseArticleData.ArticleKey)
  // if decodeArticleKeyErr != nil {
  //   ctx.JSON(http.StatusBadRequest, gin.H{"error": decodeArticleKeyErr.Error()})
  //   ctx.Abort()
  //   return
  // }
  // if len(articleKey) != 32 {
  //   ctx.JSON(http.StatusBadRequest, gin.H{"error": "width of ArticleKey must be 32"})
  //   ctx.Abort()
  //   return
  // }
  // ctx.Set("articleKey", articleKey)
}

// type ArticleRevisionMapData struct {
//   PreviousArticleRevisionKey string `json:"ArticleBody"`
//   ArticleBody string `json:"ArticleBody"`
//   Author string `json:"Author"`
//   Dateline string `json:"Dateline"`
// }

// func getPreviousArticleRevisionKey(mapLeafValue []byte) ([]byte, error) {
//
// }

func ReviseArticle(ctx *gin.Context) {
  // 1) validate the articleKey
  reviseArticleData := ReviseArticleData{}
  bindErr := ctx.ShouldBindBodyWith(&reviseArticleData, binding.JSON)
  if bindErr != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
    ctx.Abort()
    return
  }
  articleKey, decodeArticleKeyErr := b64.URLEncoding.DecodeString(reviseArticleData.ArticleKey)
  if decodeArticleKeyErr != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": decodeArticleKeyErr.Error()})
    ctx.Abort()
    return
  }
  if len(articleKey) != 32 {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": "width of ArticleKey must be 32"})
    ctx.Abort()
    return
  }

  // 2) get the latest revision of this article from the map
  isExists, mapLeafHash, mapLeafValue, getLeafErr := mapDatalayer.GetLeaf(ctx, articleKey)
  if isExists != true {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": "ArticleID does not exist"})
    ctx.Abort()
    return
  }

  //3) get previousArticleRevisionKey


  publishArticleData, _ := ctx.Get("publishArticleData")
  leafData, marshalErr := json.Marshal(publishArticleData)
  if marshalErr != nil {
    ctx.JSON(http.StatusInternalServerError, gin.H{})
    ctx.Abort()
    return
  }

  fmt.Printf("HELLO: %+v %+v %+v %+v\n", leafData, mapLeafValue, mapLeafHash, getLeafErr)
  //
  // _, _, proof, _, _, isDup, addLeafErr := logDatalayer.AddLeaf(ctx, leafData)
  // if addLeafErr != nil {
  //   ctx.JSON(http.StatusInternalServerError, gin.H{})
  //   ctx.Abort()
  //   return
  // }
  //
  // fmt.Printf("HLLOO %+v %v %v\n", proof, isDup, articleID)
  //
  // uDec, _ := b64.URLEncoding.DecodeString(articleID)
  // hash := []byte(uDec)
  // key := hash[:]
  // if (isDup == false) {
  //   addLeafErr := mapDatalayer.AddLeaf(ctx, key, leafData)
  //   if addLeafErr != nil {
  //     ctx.JSON(http.StatusInternalServerError, gin.H{})
  //     ctx.Abort()
  //     return
  //   }
  // }
  // mapLeaf, mapRoot, getLeafErr := mapDatalayer.GetLeaf(ctx, key)
  // if getLeafErr != nil {
  //   ctx.JSON(http.StatusInternalServerError, gin.H{})
  //   ctx.Abort()
  //   return
  // }
  //
  // if isDup == true {
  //   ctx.JSON(200, gin.H{"proof": proof, "mapLeaf": mapLeaf, "mapRoot": mapRoot})
  //   ctx.Abort()
  //   return
  // }
  //
  // ctx.JSON(201, gin.H{"proof": proof, "mapLeaf": mapLeaf, "mapRoot": mapRoot})
}
