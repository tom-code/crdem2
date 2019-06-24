package main

import (
  "k8s.io/apimachinery/pkg/fields"
  "k8s.io/client-go/kubernetes"
  "fmt"
  "k8s.io/client-go/tools/clientcmd"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  muj "crdem/pkg/apis/generated/clientset/versioned"
  "k8s.io/client-go/tools/cache"
  mj "crdem/pkg/apis/svc1/v1"
)

func main() {
  config, err := clientcmd.BuildConfigFromFlags("", "admin.conf")
  if err != nil {
    panic(err.Error())
  }

  //clientset for core stuff
  clientset, err := kubernetes.NewForConfig(config)
  if err != nil {
    panic(err.Error())
  }


  fmt.Println("----list pods")
  pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
  if err != nil {
    panic(err.Error())
  }
  for _, p := range pods.Items {
    fmt.Println(p.Name)
  }



  //client set for extension
  cs, err := muj.NewForConfig(config)
  if err != nil {
    panic(err.Error())
  }
  
  fmt.Println("----list extension objects")
  s := cs.Svc1V1()
  zz := s.TestTypes("default")
  lst, err := zz.List(metav1.ListOptions{})
  if err != nil {
    panic(err.Error())
  }

  for _, item := range lst.Items {
    fmt.Println(item.Name)
    fmt.Println(item.Options)
    fmt.Println(item.Spec.Blah)
    fmt.Println(item.Spec.Replicas)
  }

  //try to set watch for extensin list
  watch, err := zz.Watch(metav1.ListOptions{})
  if err != nil {
    panic(err)
  }
  watch.Stop()


  //lw := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "pods", "", fields.Everything())
  lw := cache.NewListWatchFromClient(cs.Svc1V1().RESTClient(), "testtypes", "", fields.Everything())
  //lw := cache.NewListWatchFromClient(clientset2.ApiextensionsV1beta1().RESTClient(), "testtypes", "", fields.Everything())


  fmt.Println("----list using watcher")
  ls, err := lw.List(metav1.ListOptions{})
  if err != nil {
    panic(err)
  }
  ls2 := ls.(*mj.TestTypeList)
  for _, item := range ls2.Items {
    fmt.Println(item.Name)
  }

  fmt.Println("----wait for notifications")
  //notifications through watcher
  _, controller := cache.NewInformer(
    lw,  &mj.TestType{}, 0,
      cache.ResourceEventHandlerFuncs{
        AddFunc: func(obj interface{}) {
          fmt.Printf("added: %s \n", obj)
        },
        DeleteFunc: func(obj interface{}) {
          fmt.Printf("deleted: %s \n", obj)
        },
        UpdateFunc: func(oldObj, newObj interface{}) {
          fmt.Printf("changed \n")
        },
      },
  )
  stop := make(chan struct{})
  controller.Run(stop)
}
