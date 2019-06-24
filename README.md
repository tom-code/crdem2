# crdem
minimal template for k8s crd generated client

```
cd $GOPATH/src
git clone <this>
cd crdem
sh hack/gen.sh

kubectl apply -f a.crd
kubectl apply -f a.obj

go run t1.go
```
