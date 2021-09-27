# Kubernetes learning path task solutions:

Install kind tool - https://kind.sigs.k8s.io/docs/user/quick-start/#installation


Cluster setup step: Create cluster: kind create cluster --config kind-cluster-config/kind-multi-node-cluster.yaml

Cluster clean-up step: Delete cluster: kind delete cluster


Use command 'kubectl get pods -w' to monitor live state changes in pods from a different Terminal session


``` {.sourceCode .bash}
1) Kubernetes Components

Task 1: kubectl describe pod kube-apiserver-kind-control-plane -n kube-system
        kubectl describe pod kube-controller-manager-kind-control-plane -n kube-system
        kubectl describe pod kube-scheduler-kind-control-plane -n kube-system

Task 2: https://kubernetes.io/docs/reference/command-line-tools-reference/kube-controller-manager/ --controller strings


2) Workloads

Task 1: kubectl run ping-google-pod --image=alpine ping 8.8.8.8
        kubectl delete pod ping-google-pod

Task 2: kubectl apply -f workloads/init-container.yaml
        kubectl logs static-file-app-pod
        kubectl exec static-file-app-pod -c static-file-app-container -- rm /tmp/intro.txt
        kubectl delete pod static-file-app-pod

The container (and not the pod) tries to restart after deleting /tmp/intro.txt file but fails with an Error status.
This happens because the /tmp/intro.txt file does not exist after container restart since file creation occurs in initContainer.
Container restarts do not trigger initContainer. It only runs on pod start/restart. 

Task 3: kubectl apply -f workloads/nginx-daemonset.yaml
        kubectl delete daemonset nginx-app
        
Task 4: kubectl apply -f workloads/nginx-daemonset.yaml
        Change the image tag to 'latest-not-available' in workloads/nginx-daemonset.yaml and run the Task 3 command again to apply the update
        kubectl rollout status daemonset.apps/nginx-app
        Ctrl + C
        kubectl rollout undo daemonset.apps/nginx-app
        kubectl delete daemonset nginx-app

Task 5: kubectl apply -f workloads/postgresql-statefulset.yaml
        kubectl delete statefulset postgresql-database

Task 6: docker cp workloads/kube-scheduler-modified.yaml kind-control-plane:/etc/kubernetes/manifests/
        docker exec -it kind-control-plane /bin/bash
        cd /etc/kubernetes/manifests/
        mv kube-scheduler-modified.yaml kube-scheduler.yaml
        Ctrl + D
        kubectl run ping-google-pod --image=alpine ping 8.8.8.8
        kubectl get pods -w
        Revert back the image tag change made to kube-scheduler.yaml manifest


3) Kubernetes Objects

Task 1: kubectl apply -f k8s-objects/nginx-deployment-replicas.yaml
        Manually change the label in the pod template of nginx-deployment-replicas.yaml file & run the above command again to update those changes
        kubectl delete deployment nginx-deployment


4) Networking

Task 1: kubectl apply -f networking/nginx-deployment.yaml
        kubectl run test-pod-ns-one -n ns-one --image=busybox:latest ping 8.8.8.8
        Note down the Endpoint in below command
        kubectl describe service -n ns-one nginx-service
        kubectl exec -it -n ns-one test-pod-ns-one -- /bin/sh
        wget http://(above-endpoint)
        cat index.html
        Ctrl + D
        Experiment - Exposed the nginx service out of the kind cluster as well (Try http://localhost:80 on browser)
        kubectl delete namespace ns-one

Task 2: kubectl apply -f networking/nginx-deployment.yaml
        kubectl create namespace ns-two
        kubectl run test-pod-ns-two -n ns-two --image=busybox:latest ping 8.8.8.8
        Note down the Endpoint in below command
        kubectl describe service -n ns-one nginx-service
        kubectl exec -it -n ns-two test-pod-ns-two -- /bin/sh
        wget http://(above-endpoint)
        cat index.html
        Ctrl + D
        kubectl delete namespace ns-one
        kubectl delete namespace ns-two

Task 3: kubectl create deployment nginx --image=nginx
        kubectl expose deployment nginx --port=80
        kubectl run busybox --rm -ti --image=busybox -- /bin/sh
        wget --spider --timeout=1 nginx
        kubectl apply -f networking/nginx-policy.yaml
        kubectl run busybox --rm -ti --image=busybox -- /bin/sh
        wget --spider --timeout=1 nginx
        kubectl run busybox --rm -ti --labels="access=true" --image=busybox -- /bin/sh
        wget --spider --timeout=1 nginx
        kubectl delete deployment nginx
        kubectl delete sevice nginx


5) Configuration

Task 1: kubectl apply -f configuration/configmap-pod.yaml
        kubectl exec -it test-pod -- /bin/sh
        env
        Look for environment variable 'NAME' with value as 'Shlok Chaudhari'
        Ctrl + D
        kubectl delete pod test-pod
        kubectl delete configmap configmap-example

Task 2: kubectl apply -f configuration/secret-pod.yaml
        kubectl exec -it test-pod -- /bin/bash
        env
        Look for environment variable 'PASSWORD' with value as 'wfkbkeo11'
        Ctrl + D
        kubectl delete pod test-pod
        kubectl delete secret secret-example

Task 3: echo -n '1f2d1e2e67df' > configuration/password.txt
        kubectl create secret generic pass-example --from-file=password=configuration/password.txt
        kubectl get secret pass-example -o yaml
        Observe that the value to password key is decoded
        kubectl create secret generic pass-example-literal --from-literal=password='1f2d1e2e67df'
        kubectl get secret pass-example-literal -o yaml
        Observe that the value to password key is decoded


6) Storage

Task 1: https://kubernetes.io/docs/concepts/storage/dynamic-provisioning/

Task 2: kubectl apply -f storage/pod-with-two-containers.yaml
        kubectl exec -c test-container-1 -it test-pod -- /bin/sh
        echo "My name is Shlok" >> /cache/intro.txt
        Ctrl + D
        kubectl exec test-pod -c test-container-2 -- cat /cache/intro.txt
        kubectl delete pod test-pod

Task 3: https://kubernetes.io/docs/concepts/storage/volumes/#hostpath
        kubectl apply -f storage/hostPath-example.yaml
        kubectl exec -it test-webserver -- ls /var/local/
        kubectl exec -it test-webserver -- ls /var/local/aaa
        kubectl describe pod test-webserver
        Note down the node pod is scheduled on
        docker exec -it (node-name) ls /var/local/
        docker exec -it (node-name) ls /var/local/aaa


7) Security

Task 1: https://kubernetes.io/docs/reference/access-authn-authz/rbac/
        https://stackoverflow.com/questions/47973570/kubernetes-log-user-systemserviceaccountdefaultdefault-cannot-get-services
        kubectl apply -f security/kanister-pod.yaml
        kubectl create rolebinding pod-reader-pod --role=pod-reader --serviceaccount=default:default
        kubectl exec -it kanister-pod -- kubectl get pods

Task 2: To enable RBAC, start the API server with the --authorization-mode flag set to a comma-separated list that includes RBAC
        For example: kube-apiserver --authorization-mode=Example,RBAC --other-options --more-options

Task 3: https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/
```
