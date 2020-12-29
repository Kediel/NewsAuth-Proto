# blue
API wrapper for [Google Trillian](https://github.com/google/trillian) cryptographically verified logs

### Flows
![Flow 1, editor publishes news](/diagrams/flow-1-diagram.png =250x)
![Flow 2, editor publishes revision of existing news](/diagrams/flow-2-diagram.png =250x)



### Ignore Below This Line
participant Editor
participant CMS
participant Client
participant Log Server
participant Map Server
participant MySQL
Title: Flow 1. Editor publishes news
Editor->CMS: presses publish
CMS->Client: POST /v1/news
Client->Client: key = articleID = sha256(content)
Client->Log Server: addLogLeaf(key, content)
Log Server->MySQL: write
MySQL-->Log Server: ok
Log Server-->Client: proof, isDup
Note right of Client: if isDup=false
Client->Map Server: addMapLeaf(key, content)
Map Server->MySQL: write
MySQL-->Map Server: ok
Map Server-->Client: ok
Note right of Client: endif
Client-->CMS: articleID, proof
CMS-->Editor: success
