Add a tag
 git tag -a v1.34.0-alpha.dracelexpr

KUBE_BUILD_PLATFORMS=linux/amd64 KUBE_DOCKER_IMAGE_TAG=v1.34.0-dracelexpr make quick-release-images

kind load docker-image registry.k8s.io/kube-controller-manager-amd64:v1.34.0-dracelexpr --name dra-example-driver-cluster
kind load docker-image registry.k8s.io/kube-apiserver-amd64:v1.34.0-dracelexpr --name dra-example-driver-cluster
kind load docker-image registry.k8s.io/kube-scheduler-amd64:v1.34.0-dracelexpr --name dra-example-driver-cluster
kind load docker-image registry.k8s.io/kube-proxy-amd64:v1.34.0-dracelexpr --name dra-example-driver-cluster


then docker exec -it /bin/bash into control plane container and update manifests - change the image and change log verbosity to 6 -- pods will restart automatically. Verify with get pod that control plane components are running with the above images - change apiserver, kcm and scheduler ..not sure if kcm needs to be updated.

Logs from dra code are in scheduler log
kubectl logs kube-scheduler-dra-example-driver-cluster-control-plane -n kube-system > s.log


