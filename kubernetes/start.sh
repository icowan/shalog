


docker run --name=kube-apiserver -v /etc/ssl/certs:/etc/ssl/certs \
-v /etc/pki:/etc/pki \
-v /etc/kubernetes/pki:/etc/kubernetes/pki \
k8s.gcr.io/kube-apiserver:v1.13.4 \
kube-apiserver \
--authorization-mode=Node,RBAC \
--advertise-address=129.204.31.254 \
--allow-privileged=true \
--client-ca-file=/etc/kubernetes/pki/ca.crt \
--enable-admission-plugins=NodeRestriction \
--enable-bootstrap-token-auth=true \
--etcd-cafile=/etc/kubernetes/pki/etcd/ca.crt \
--etcd-certfile=/etc/kubernetes/pki/apiserver-etcd-client.crt \
--etcd-keyfile=/etc/kubernetes/pki/apiserver-etcd-client.key \
--etcd-servers=https://127.0.0.1:2379 \
--insecure-port=0 \
--kubelet-client-certificate=/etc/kubernetes/pki/apiserver-kubelet-client.crt \
--kubelet-client-key=/etc/kubernetes/pki/apiserver-kubelet-client.key \
--kubelet-preferred-address-types=InternalIP,ExternalIP,Hostname \
--proxy-client-cert-file=/etc/kubernetes/pki/front-proxy-client.crt \
--proxy-client-key-file=/etc/kubernetes/pki/front-proxy-client.key \
--requestheader-allowed-names=front-proxy-client \
--requestheader-client-ca-file=/etc/kubernetes/pki/front-proxy-ca.crt \
--requestheader-extra-headers-prefix=X-Remote-Extra- \
--requestheader-group-headers=X-Remote-Group \
--requestheader-username-headers=X-Remote-User \
--secure-port=6443 \
--service-account-key-file=/etc/kubernetes/pki/sa.pub \
--service-cluster-ip-range=10.96.0.0/12 \
--tls-cert-file=/etc/kubernetes/pki/apiserver.crt \
--tls-private-key-file=/etc/kubernetes/pki/apiserver.key

# etcd
nohup ./etcd \
--advertise-client-urls=https://172.16.0.4:2379 \
--cert-file=/etc/kubernetes/pki/etcd/server.crt \
--client-cert-auth=true \
--data-dir=/var/lib/etcd \
--initial-advertise-peer-urls=https://172.16.0.4:2380 \
--initial-cluster=master=https://172.16.0.4:2380 \
--key-file=/etc/kubernetes/pki/etcd/server.key \
--listen-client-urls=https://127.0.0.1:2379,https://172.16.0.4:2379 \
--listen-peer-urls=https://172.16.0.4:2380 \
--name=master \
--peer-cert-file=/etc/kubernetes/pki/etcd/peer.crt \
--peer-client-cert-auth=true \
--peer-key-file=/etc/kubernetes/pki/etcd/peer.key \
--peer-trusted-ca-file=/etc/kubernetes/pki/etcd/ca.crt \
--snapshot-count=10000 \
--trusted-ca-file=/etc/kubernetes/pki/etcd/ca.crt >> etcd.log 2>&1 &

test

ETCDCTL_API=3 ./etcdctl --endpoints=https://[127.0.0.1]:2379 \
--cacert=/etc/kubernetes/pki/etcd/ca.crt \
--cert=/etc/kubernetes/pki/etcd/healthcheck-client.crt \
--key=/etc/kubernetes/pki/etcd/healthcheck-client.key

## apiserver

```
nohup ./kube-apiserver \
--authorization-mode=Node,RBAC \
--advertise-address=129.204.31.254 \
--allow-privileged=true \
--client-ca-file=/etc/kubernetes/pki/ca.crt \
--enable-admission-plugins=NodeRestriction \
--enable-bootstrap-token-auth=true \
--etcd-cafile=/etc/kubernetes/pki/etcd/ca.crt \
--etcd-certfile=/etc/kubernetes/pki/apiserver-etcd-client.crt \
--etcd-keyfile=/etc/kubernetes/pki/apiserver-etcd-client.key \
--etcd-servers=https://127.0.0.1:2379 \
--insecure-port=0 \
--kubelet-client-certificate=/etc/kubernetes/pki/apiserver-kubelet-client.crt \
--kubelet-client-key=/etc/kubernetes/pki/apiserver-kubelet-client.key \
--kubelet-preferred-address-types=InternalIP,ExternalIP,Hostname \
--proxy-client-cert-file=/etc/kubernetes/pki/front-proxy-client.crt \
--proxy-client-key-file=/etc/kubernetes/pki/front-proxy-client.key \
--requestheader-allowed-names=front-proxy-client \
--requestheader-client-ca-file=/etc/kubernetes/pki/front-proxy-ca.crt \
--requestheader-extra-headers-prefix=X-Remote-Extra- \
--requestheader-group-headers=X-Remote-Group \
--requestheader-username-headers=X-Remote-User \
--secure-port=6443 \
--service-account-key-file=/etc/kubernetes/pki/sa.pub \
--service-cluster-ip-range=10.96.0.0/12 \
--tls-cert-file=/etc/kubernetes/pki/apiserver.crt \
--tls-private-key-file=/etc/kubernetes/pki/apiserver.key  >> kube-apiserver.log 2>&1 &
```

## kube-scheduler

nohup  ./kube-scheduler --address=127.0.0.1 --kubeconfig=/etc/kubernetes/scheduler.conf --leader-elect=true   >> kube-scheduler .log 2>&1 &

## kube-controller-manager

./kube-controller-manager \
--address=127.0.0.1 \
--allocate-node-cidrs=true \
--authentication-kubeconfig=/etc/kubernetes/controller-manager.conf \
--authorization-kubeconfig=/etc/kubernetes/controller-manager.conf \
--client-ca-file=/etc/kubernetes/pki/ca.crt \
--cluster-cidr=10.244.0.0/16 \
--cluster-signing-cert-file=/etc/kubernetes/pki/ca.crt \
--cluster-signing-key-file=/etc/kubernetes/pki/ca.key \
--controllers=*,bootstrapsigner,tokencleaner \
--kubeconfig=/etc/kubernetes/controller-manager.conf \
--leader-elect=true \
--node-cidr-mask-size=24 \
--requestheader-client-ca-file=/etc/kubernetes/pki/front-proxy-ca.crt \
--root-ca-file=/etc/kubernetes/pki/ca.crt \
--service-account-private-key-file=/etc/kubernetes/pki/sa.key \
--use-service-account-credentials=true

## kubelet

## kube-proxy




{
  "kind": "Service",
  "apiVersion": "v1",
  "metadata": {
    "name": "mysql",
    "namespace": "kpaas",
    "labels": {
      "app": "mysql"
    }
  },
  "spec": {
    "ports": [
      {
        "name": "tcp-3306",
        "protocol": "TCP",
        "port": 3306,
        "targetPort": 3306,
        "nodePort": 32306
      }
    ],
    "selector": {
      "app": "mysql"
    },
    "type": "NodePort",
    "externalTrafficPolicy": "Cluster"
  },
  "status": {
    "loadBalancer": {}
  }
}

kind: Service
apiVersion: v1
metadata:
  name: harbor-hub
  labels:
    app: harbor-hub
spec:
  ports:
    - name: http-80
      protocol: TCP
      port: 80
      targetPort: 80
  type: ClusterIP
  sessionAffinity: ClientIP
---
apiVersion: v1
kind: Endpoints
metadata:
  name: harbor-hub
subsets:
- addresses:
  - ip: 39.106.40.14
  ports:
  - port: 80
    protocol: TCP
---
kind: Ingress
apiVersion: extensions/v1beta1
metadata:
  name: kplcloud
spec:
  rules:
    - host: hub.kpaas.nsini.com
      http:
        paths:
          - backend:
              serviceName: harbor-hub
              servicePort: 80

---
kind: Deployment
apiVersion: apps/v1
metadata:
  generation: 1
  labels:
    app: cardbill
    language: Golang
  name: cardbill
  namespace: kpaas
spec:
  minReadySeconds: 10
  progressDeadlineSeconds: 600
  replicas: 1
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: cardbill
  strategy:
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
    type: RollingUpdate
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: cardbill
        language: Golang
    spec:
      containers:
      - env:
        - name: ENV
          value: prod
        - name: POD_NAME
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: metadata.name
        - name: POD_NAMESPACE
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: metadata.namespace
        - name: INSTANCE_IP
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: status.podIP
        image: hub.kpaas.nsini.com/golang
        imagePullPolicy: IfNotPresent
        name: cardbill
        ports:
        - containerPort: 8080
          protocol: TCP
        resources:
          limits:
            memory: 64Mi
          requests:
            memory: 64Mi
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
        volumeMounts:
        - mountPath: /etc/localtime
          name: tz-config
      dnsPolicy: ClusterFirst
      imagePullSecrets:
      - name: regcred
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      terminationGracePeriodSeconds: 30
      volumes:
      - hostPath:
          path: /usr/share/zoneinfo/Asia/Shanghai
          type: ""
        name: tz-config
status: {}

---
kind: Service
apiVersion: v1
metadata:
  labels:
    app: cardbill
  name: cardbill
  namespace: kpaas
spec:
  ports:
  - name: http-8080
    port: 8080
    protocol: TCP
    targetPort: 8080
  selector:
    app: cardbill
  sessionAffinity: None
  type: ClusterIP
status:
  loadBalancer: {}


---
kind: Ingress
apiVersion: extensions/v1beta1
metadata:
  annotations:
    kubernetes.io/ingress.class: traefik
  name: cardbill
  namespace: kpaas
spec:
  rules:
  - host: cardbill.nsini.com
    http:
      paths:
      - backend:
          serviceName: cardbill
          servicePort: 8080
status:
  loadBalancer: {}

apiVersion: v1
data:
  app.cfg: |

[server]
app_name = cardbill
debug = true
session_timeout = 7200
app_key = ab23&f9a812bd!@3r-=1203
http_static = ./dist/

[mysql]
mysql_host = mysql
mysql_port = 3306
mysql_user = root
mysql_password = gbA^zLR$FQsg
mysql_database = cardbill


[github]
client_id = 84a8596e2d0efc53e9d0
client_secret = 3fe726d8dfc0bb22a4f1cf797698edc72178052a

[cors]
allow = true
origin = http://localhost:8000
methods = GET,POST,OPTIONS,PUT,DELETE
headers = Origin,Content-Type,Authorization,mode,cors,x-requested-with,Access-Control-Allow-Origin,Access-Control-Allow-Credentials

kind: ConfigMap
metadata:
  name: cardbill
  namespace: kpaas
