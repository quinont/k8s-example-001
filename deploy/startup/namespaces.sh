kubectl create namespace api 
kubectl create namespace app2
kubectl create namespace database
kubectl create namespace web 
kubectl label namespace web istio-injection=enabled
kubectl label namespace app2 istio-injection=enabled
kubectl label namespace api istio-injection=enabled
kubectl label namespace database istio-injection=enabled
