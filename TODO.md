- [x] set up config & env vars
- [x] set up connections with postgres & redis + migrations + models
- [x] abstractions (structs & interfaces; unimplemented methods)
- [x] repo layer (db & cache)
- [x] service layer (app logic)
- [x] jwt implementation & logs (via zap) - /pkg/
- [x] transport layer (grpc server conf & handlers)
- [x] configure app init - /internal/app/ && /cmd/sso/main.go
- [x] dockerize + Makefile
- [ ] hand testing grpc api via postman
- [ ] unit testing
- [ ] ...

Make sure to properly:
- log
- comment
- write code, which is:
    - independent (uses abstractions)
    - clean
    - simple