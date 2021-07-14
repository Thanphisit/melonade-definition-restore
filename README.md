## Start import

## Dump
Dump all task/workflows to files
```sh
go run main.go -mode dump -url http://melonade.melonade.staging.tel.internal/api/process-manager
```

## Restore
Create new task/workflows from backup
```sh
go run main.go -mode restore -url http://melonade.melonade.uat1.tel.internal/api/process-manager
```

## Upgrade
Update existing task/workflows from backup
```sh
go run main.go -mode upgrade -url http://melonade.melonade.uat1.tel.internal/api/process-manager
```

## Clean
Remove all workflow
```sh
go run main.go -mode clean -url http://melonade.melonade.uat1.tel.internal/api/process-manager
```