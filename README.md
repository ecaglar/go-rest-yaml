# go-appmetadata-yaml
REST API to support read/write application metadata as yaml/json payloads with a integrated in-mem db 

- Extended Logger
- Custom validation function integration to handler
- GET and POST to create and get application metadata
- yaml and json payload support
- Acecept header support (application/json) default is yaml
- handler chaining for routers

Sample usage of the server is:


	logger  := logger.CreateLogger()
	db      := memstore.CreateInMemDB()
	server  := server.CreateServer(logger,db)

	http.Handle("/",server.Routers)
	logger.LogInfo("Listening localhost 8080...")
	http.ListenAndServe("localhost:8080", server.Routers)
