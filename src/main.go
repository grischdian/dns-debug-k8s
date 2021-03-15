package main

import (
	"context"
	"fmt"
	"time"
        "net"
        "os"

//	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

func main() {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
        var debug string
        debug = os.Getenv("DEBUG")

        for {
                currentTime := time.Now()
                ips, err := net.LookupIP(os.Getenv("SERVICENAME")) //lookup on specific kubernetes-service
                if err != nil { //if lookup failed
                   fmt.Printf("%s %s lookup was not succesfull from ClusterIP \n", currentTime.Format("2006.01.02 15:04:05"), os.Getenv("SERVICENAME"))
                   endpoints := clientset.CoreV1().Endpoints(os.Getenv("NAMESPACE")) //get endpoints from namespace
                   list, err := endpoints.List(context.TODO(), metav1.ListOptions{FieldSelector:"metadata.name=kubernetes"}) //get specific endpoint
                   if err != nil {
                     panic(err)
                   }
                   for _, d := range list.Items { //start iterating over specific endpoint to get endpoint IPs
                     for _, x := range d.Subsets {
                        for _, a := range x.Addresses { //iterate over all endpoints IPs
                          endpoints := a.IP
                          //fmt.Println(endpoints) //print endpoint IPs
                          endpoints += ":53" //sick hack to concat String
                          r := &net.Resolver{
                            PreferGo: true,
                            Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
                              d := net.Dialer{
                                Timeout: 1,
                              }
                              return d.DialContext(ctx, "udp", endpoints) //use endpoints as DNS-Server
                            },
                          }
                          ip, err := r.LookupHost(context.Background(), os.Getenv("SERVICENAME"))
                          if err !=nil { //DNS lookup failed for specific DNS-Server
                            fmt.Printf("%s %s lookup was not succesfull from %s \n", currentTime.Format("2006.01.02 15:04:05"), os.Getenv("SERVICENAME"), a.IP)
                          }
                          for _, b := range ip {
                            fmt.Printf("%s %s with ip %s lookup was succesfull from %s \n", currentTime.Format("2006.01.02 15:04:05"), os.Getenv("SERVICENAME"), b,  a.IP)
                          }
                        }
                     }
                   }
                }
                if debug == "enabled" {
                  for _, ip := range ips { //if lookup was good
                    fmt.Printf("%s %s with ip %s lookup was succesfull from ClusterIP \n", currentTime.Format("2006.01.02 15:04:05"), os.Getenv("SERVICENAME"), ip)
                  }
                }


                time.Sleep(2 * time.Second)
        }

}
