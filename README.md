# Herlighet

Let There Be Glory!


## Deployment

Deployment is managed using ansible.

```
ansible $ ansible-playbook -bKu RA_[RA user] -i [dev or prod] site.yml
```

You will be prompted for your RA user's password, used to sudo on the target machines.

Please note that the ansible playbook also invokes vault locally to obtain the passwords used to access the herlighet databases. 

Therefore, vault and a valid access token (obtained with `vault login -method=oidc`) is required.

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
