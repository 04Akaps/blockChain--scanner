# blockChain--scanner
체인에 대한 Scanner를 개발하는 간단한 Scanner 모듈

읽어야 하는 블록이 있다면, 노드에 Call을 전송하여, Receipt를 불러오고
해당 정보를 토대로 Tx를 생성 그 후, Tx를 DB에 저장

`upsert`하는 부분을 최대한 줄이기 위해서, `WriteModel`을 사용을 하였고

빠른 처리를 위해서, 읽어야 하는 블록이 있다면, 모두 go루틴으로 전송
- 서버 성능에 따라서 부하가 발생 할 수 있으니, 고도화가 되어야 한다면, 따로 모듈을 돌려서, DB에 블록이 저장이 되면,
- 해당 컬렉션을 `Stream`하면서 처리하는 방법도 가능