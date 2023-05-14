```mermaid
classDiagram

main --> config
main --> handler
main --> rotator

handler --> fileService
handler --> logService
handler --> queryService
rotator --> fileService

queryService --> logService
queryService --> fileService
queryService --> letheql
logService --> fileService
fileService --> storageDriver

storageDriver <|-- filesystemDriver

class main {
    config
    handler
    rotator
    main()
}

class rotator {
    fileService
    Start()
}

class handler {
    queryService
    logService
    fileService
    Metadata()
    Query()
    QueryRange()
    Targets()
}

class queryService {
    letheql
    logService
    ExecuteQuery()
}

class logService {
    fileService
    ListTargets()
}

class fileService {
    storageDriver
    ListDirs()
    DeleteByAge()
    DeleteBySize()
}


class letheql {
    ProcQuery()
}
```
