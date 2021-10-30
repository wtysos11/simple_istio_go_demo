## Simple Go App in Istio

Changed from [httpbin](https://raw.githubusercontent.com/istio/istio/release-1.11/samples/httpbin/httpbin.yaml), and add gateway and virtual services following [official document](https://istio.io/latest/docs/tasks/traffic-management/ingress/ingress-control/)

K8s and istio is necessary.

* build image: `docker build . -t wtyfft:0.0.1`
* tag and upload image 
* Deploy: `kubectl apply -f app.yaml`
* get ip following official document
  * `export INGRESS_HOST=$(kubectl get po -l istio=ingressgateway -n istio-system -o jsonpath='{.items[0].status.hostIP}')`
  * `export INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].nodePort}')`
* get result: `curl http://$INGRESS_HOST:$INGRESS_PORT/fft/real/10`
* prometheus PromQL: `istio_requests_total{destination_workload_namespace="wty-istio"}`

