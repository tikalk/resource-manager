### How to develop with Tilt

From the root folder, do the following (assuimng [Tilt](https://docs.tilt.dev) is installed)
```
cd tilt
tilt up
```


### how to generate this repo

```bash

mkdir resource-manager
cd resource-manager
operator-sdk init --domain tikalk.com --repo gitlab.com/tikalk.com/resource-manager
operator-sdk create api --group resource-managment --version v1alpha1 --kind ResourceManager --resource --controller
# to create the resources
make generate
# to create the CRDs
make manifests
```
