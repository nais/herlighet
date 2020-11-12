# Herlighet

Let There Be Glory!


## Herlighet ops/database maintenance

Getting credentials for the herlighet metadata/configuration database:

dev:
```
vault read postgresql/preprod-fss/static-creds/herlighet-static-admin
```

prod:
```
vault read postgresql/prod-fss/static-creds/herlighet-static-admin
```
