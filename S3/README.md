# APPs

## simple-app

```
go run main.go -b app -f file-sample.json
awslocal s3 ls app/
awslocal s3 cp s3://app/file-sample.json myfile.json
```



# CLI
```
awslocal s3 mb s3://mybucket
awslocal s3 ls
```