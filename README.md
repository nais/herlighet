# Herlighet

Let There Be Glory!


## Herlighet ops/database maintenance

VMs that the handlers run on today:

```
dev: b27apvl00485.preprod.local
prod: a01apvl00454.adeo.no
```

Getting credentials for the herlighet metadata/configuration database:

dev:
```
vault read postgresql/preprod-fss/static-creds/herlighet-static-admin
```

prod:
```
vault read postgresql/prod-fss/static-creds/herlighet-static-admin
```
