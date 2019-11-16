package main

import (
  "k8s.io/apimachinery/pkg/fields"
  "k8s.io/client-go/kubernetes"
  "fmt"
  "k8s.io/client-go/tools/clientcmd"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  muj "a.com/crdem/pkg/apis/generated/clientset/versioned"
  "k8s.io/client-go/tools/cache"
  mj "a.com/crdem/pkg/apis/svc1/v1"
  appsv1 "k8s.io/api/apps/v1"
  corev1 "k8s.io/api/core/v1"

)

var (
  clientset *kubernetes.Clientset
)

func createDepl() {
  replicas := int32(1)
  dep := &appsv1.Deployment{
    ObjectMeta: metav1.ObjectMeta {
      Name: "t1",
      Namespace: "default",
      Annotations: map[string]string {"a": "1"},
    },
    Spec: appsv1.DeploymentSpec {
      Replicas: &replicas,
      Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string {"sel": "abc"},
      },
      Template: corev1.PodTemplateSpec{
        ObjectMeta: metav1.ObjectMeta{
					Labels:      map[string]string {"sel": "abc"},
					Annotations: map[string]string {"a": "1"},
        },
        Spec: corev1.PodSpec{
          Containers: []corev1.Container { {
              Name: "cont1",
              Image: "10.10.2.213:5000/amf-eefat:iB4.0.1",
              Command: []string{"/bin/sleep"},
              Args: []string{"1000"},
              Env: []corev1.EnvVar {
                corev1.EnvVar { Name: "E1", Value: "V1" },
              },
            },
          },
        },
      },
    },
  }
  dep2, err := clientset.AppsV1().Deployments("default").Create(dep)
  fmt.Println(dep2)
  fmt.Println(err)
}

func DeleteDepl() {
  clientset.AppsV1().Deployments("default").Delete("t1", &metav1.DeleteOptions{})
}

func main() {
  config, err := clientcmd.BuildConfigFromFlags("", "admin.conf")
  if err != nil {
    panic(err.Error())
  }

  //clientset for core stuff
  clientset, err = kubernetes.NewForConfig(config)
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
  zz := s.TestTypes("")
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
          createDepl()
        },
        DeleteFunc: func(obj interface{}) {
          fmt.Printf("deleted: %s \n", obj)
          DeleteDepl()
        },
        UpdateFunc: func(oldObj, newObj interface{}) {
          fmt.Printf("changed %s\n", newObj)
        },
      },
  )
  stop := make(chan struct{})
  controller.Run(stop)
}
