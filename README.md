# Herlighet

## Deployment

Deployment is managed using ansible. Although we have nowhere to run it from these days.

```
ansible $ ansible-playbook -bKu RA_[RA user] -i [dev or prod] site.yml
```

You will be prompted for your RA user's password, used to sudo on the target machines.

Please note that the ansible playbook also invokes vault locally to obtain the passwords used to access the herlighet databases. 

Therefore, vault and a valid access token (obtained with `vault login -method=oidc`) is required.

## Herlighet maintenance

VMs that the handlers run on today:

```
dev:  a30apvl042.oera.no (main), a30apvl044.oera.no (backup)
prod: a30apvl043.oera.no (main), a30apvl045.oera.no (backup)
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

These credentials can also be obtained from [vault GUI](https://vault.adeo.no). 

Since there are no nodes with access to the nodes in question we are updating the database password via the VMWare console (accessible from operations image):

* Log on with RA-account and change the password in the file /etc/herlighet.env.

* Restart the service: 
```
systemctl restart herlighet
```
