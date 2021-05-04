
# Build the tool:

```
go build cmd/myhttp/myhttp.go
```

# Unit tests:

```
go test -v ./internal/httpMD5
```

#Usage:

* Run with default limit of 10 parallel requests

```
$ ./myhttp http://www.adjust.com http://google.com
http://google.com 40946232cf19252110d95db2edd79794
http://www.adjust.com fa0072f4d48191d35c8592da67309f53
```

* Set limit for parallel requests

```
$ ./myhttp -parallel 3 adjust.com google.com facebook.com yahoo.com yandex.com twitter.com reddit.com/r/funny reddit.com/r/notfunny baroquemusiclibrary.com

http://google.com 58b2b8894f1d7baa823122648ccbbdcd
http://facebook.com 9d59fc5ee71fe140242d821da9e38c7e
http://adjust.com fa0072f4d48191d35c8592da67309f53
http://yandex.com 87ac35971e679c4c519c4dd54f3d7863
http://yahoo.com 9c813a810e6194080c3c541364de8156
http://twitter.com 0de39557fa10d36d6e5048df2c055d8e
http://reddit.com/r/funny c970e51512697a62fbd7a346d6ad8400
http://baroquemusiclibrary.com 273e0a982137d81490566eb12935870e
http://reddit.com/r/notfunny c7ac2f383be40c3cbfbe2fcba8f348c3
```