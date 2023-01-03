@televisionbbs

None of this works.  I am moving crap around and wanted to get as much as possible then shift to a more comfy chair...


```
go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

```
migrate -path migrations -database sqlite3://bbs.db up
```
