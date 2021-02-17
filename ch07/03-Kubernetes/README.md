# 쿠버네티스 배포 (Chapter 7 : 서비스 수준 gRPC 실행)

## 예제 코드 리스트
- 코드 7-4, 7-5 (서버용 쿠버네티스 배포 및 서비스 기술자) : [grpc-prodinfo-server.yaml](productinfo/server/grpc-prodinfo-server.yaml)
- 코드 7-6 (클라이언트용 쿠버네티스 잡 기술자) : [grpc-prodinfo-client-job.yaml](productinfo/client/grpc-prodinfo-client-job.yaml)
- 코드 7-7 (서비스용 쿠버네티스 인그레스 기술자) : [grpc-prodinfo-ingress.yaml](productinfo/ingress/grpc-prodinfo-ingress.yaml)

## 1. 서버 쿠버네티스 배포
다음과 같이 쿠버네티스 배포 기술자 파일을 사용해 쿠버네티스에 Deployment를 만듭니다. [grpc-prodinfo-server.yaml](productinfo/server/grpc-prodinfo-server.yaml) (코드 7-4)
```yml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: grpc-productinfo-server
spec:
  replicas: 1
  selector:
    matchLabels:
      app: grpc-productinfo-server
  template:
    metadata:
      labels:
        app: grpc-productinfo-server
    spec:
      containers:
      - name: grpc-productinfo-server
        image: kasunindrasiri/grpc-productinfo-server
        resources:
          limits:
            memory: "128Mi"
            cpu: "500m"
        ports:
        - containerPort: 50051
          name: grpc
```

다음으로 Service 정보를 배포 기술자 파일에 추가합니다. (별도로 만들 수 있지만, 같이 작성합니다.) [grpc-prodinfo-server.yaml](productinfo/server/grpc-prodinfo-server.yaml) (코드 7-5)

```yml
apiVersion: v1
kind: Service
metadata:
  name: productinfo
spec:
  selector:
    app: grpc-productinfo-server
  ports:
  - port: 50051
    targetPort: 50051
    name: grpc
  type: NodePort
```

이제 다음과 같이 쿠버네티스 명령을 사용해 배포 합니다.

```shell
$ cd productinfo
$ kubectl apply -f server/grpc-prodinfo-server.yaml
deployment.apps/grpc-productinfo-server created
service/productinfo created
```

## 2. 쿨라이언트 쿠버네티스 배포
클라이언트도 쿠버네티스 기술자 파일을 사용해 쿠버네티스 Job를 만듭니다. [grpc-prodinfo-client-job.yaml](productinfo/client/grpc-prodinfo-client-job.yaml) (쿄드 7-6)

```yml
apiVersion: batch/v1
kind: Job
metadata:
  name: grpc-productinfo-client
spec:
  completions: 1
  parallelism: 1
  template:
    spec:
      containers:
      - name: grpc-productinfo-client
        image: kasunindrasiri/grpc-productinfo-client
      restartPolicy: Never
  backoffLimit: 4
```

그런 다음 다음과 같이 쿠버네티스 명령을 사용해 Job을 실행합니다.

```shell
$ cd productinfo
$ kubectl apply -f client/grpc-prodinfo-client-job.yaml
job.batch/grpc-productinfo-client created
```

아울러, 다음과 같이 배포된 pod 정보를 확인할 수 있고, 서버가 실행된 로그를 확인할 수 있습니다.
```shell
$ kubectl get pods
NAME                                       READY   STATUS      RESTARTS   AGE
grpc-productinfo-client-988fv              0/1     Completed   0          7m34s
grpc-productinfo-server-68454d97c6-v257h   1/1     Running     0          16m

$ kubectl logs grpc-productinfo-server-68454d97c6-v257h
2021/02/17 01:50:28 New product added - ID : 82dda263-70c2-11eb-90b5-529a3dc63b21, Name : Sumsung S10
2021/02/17 01:50:28 New product retrieved - ID : value:"82dda263-70c2-11eb-90b5-529a3dc63b21
```

## 3. 서버 쿠버네티스 서비스 제공
이제 서비스를 외부로 공개하기 위해 쿠버네티스 인그레스를 다음고 같이 만듭니다. [grpc-prodinfo-ingress.yaml](productinfo/ingress/grpc-prodinfo-ingress.yaml) (쿄드 7-7)

```yml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.class: "nginx"
    nginx.ingress.kubernetes.io/ssl-redirect: "false"
    nginx.ingress.kubernetes.io/backend-protocol: "GRPC"
  name: grpc-prodinfo-ingress
spec:
  rules:
  - host: productinfo
    http:
      paths:
      - backend:
          serviceName: productinfo
          servicePort: grpc
```

우선, 다음과 같이 nginx ingress controller를 설치해야 합니다. (참조 : https://kubernetes.github.io/ingress-nginx/)

```shell
$ kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=120s
$ kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v0.44.0/deploy/static/provider/cloud/deploy.yaml
```

그런 다음 다음과 같이 쿠버네티스 명령을 사용해 인그래스를 생성합니다.

```shell
$ cd productinfo
$ kubectl apply -f ingress/grpc-prodinfo-ingress.yaml
ingress.extensions/grpc-prodinfo-ingress created
```

그리고 다음과 같이 인그래스를 확인할 수 있습니다.
```shell
$ kubectl get ingress
NAME                    CLASS    HOSTS         ADDRESS   PORTS   AGE
grpc-prodinfo-ingress   <none>   productinfo             80      2m21s
```
