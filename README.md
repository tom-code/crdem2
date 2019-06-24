# crdem
minimal template for k8s crd generated client

```
cd $GOPATH/src
git clone <this>
cd crdem
go test
go get k8s.io/code-generator@kubernetes-1.14.0
sh hack/gen.sh

kubectl apply -f a.crd
kubectl apply -f a.obj

go run t1.go
```





during go.mod creation I had to experiment with versions:
```
go get k8s.io/client-go@v11.0.0
go get k8s.io/api@kubernetes-1.14.0
go get k8s.io/apimachinery@kubernetes-1.14.0 
```
