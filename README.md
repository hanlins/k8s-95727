# Test repo for issue 95727

# Download k8s mods

Credits to @abursavich for shared a [scrip](https://github.com/kubernetes/kubernetes/issues/79384#issuecomment-521493597) for downloading the k8s mods.
In our case we're debugging `v1.19.3` run the script to download everything needed:
```bash
bash -x hack/download.sh 1.19.3
```

# Compile the binary

To test it on a Linux VM, do the following to cross compile:
```bash
GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" .
```

