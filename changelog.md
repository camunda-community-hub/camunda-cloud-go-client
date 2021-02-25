### Linux

```shell
curl -L https://github.com/camunda-community-hub/camunda-cloud-go-client/releases/download/v0.0.19/cc-linux-amd64.tar.gz | tar xzv 
sudo mv cdf /usr/local/bin
```

### macOS

```shell
curl -L  https://github.com/camunda-community-hub/camunda-cloud-go-client/releases/download/v0.0.19/cc-darwin-amd64.tar.gz | tar xzv
sudo mv cdf /usr/local/bin
```

## Changes

* adding hack headerfile (salaboy)
* adding changelog step (salaboy)
* triggers only on main (salaboy)
* adding release binaries and upload binaries to pipeline (salaboy)
