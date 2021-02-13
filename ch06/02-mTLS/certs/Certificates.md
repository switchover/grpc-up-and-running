# RSA 인증서 만들기

mTLS를 적용하려면 클라이언트와 서버 인증서를 사용해야 합니다. 
자체 서명된 인증서로 CA를 생성하고 클라이언트와 서버 모두에 대한 인증서 서명요청을 생성해야 하며, CA를 사용해 전자서명합니다.
이에 대한 절차는 다음과 같습니다.

## 1. 서버 비밀키(Private Key) 생성
다음과 같은 명령을 통해 쉽게 비밀키 파일을 생성할 수 있습니다.
```shell
$ openssl genrsa -out server.key 2048
Generating RSA private key, 2048 bit long modulus (2 primes)
.................................................+++++
...............+++++
e is 65537 (0x010001)
```
- `genrsa` : 사용할 PKI 알고리즘 지정, 일반적인 웹 서비스용으로 사용되는 RSA 지정
- `-out server.key` : 생성할 비밀키 파일명 지정(server.key)
- `2048` : 키 크기를 지정 (기본적으로 512-bit지만 일반적으로 안전한 보안을 위해 2048-bit 사용)

별도로 비밀키에는 암호문(passphrase)를 지정해야 하지만, 예제에서는 암호문을 별도로 지정하지 않고 사용합니다.

## 2. 인증기관(CA) 및 자체 서명된(self-signed) 인증서 생성
다음으로 CA에 대한 비밀키를 다음과 같이 생성합니다.
```shell
$ openssl genrsa -aes256 -out ca.key 4096
Generating RSA private key, 4096 bit long modulus (2 primes)
...................................................................................................++++
..................++++
e is 65537 (0x010001)
Enter pass phrase for ca.key:gRPC
Verifying - Enter pass phrase for ca.key:gRPC
```
비밀번호를 입력해야 하는데, 여기서는 테스트용으로 "gRPC"를 입력(2번 입력)하였습니다.

이제 CA에 대한 인증서를 유효기간 2년과 SHA-256 해쉬 알고리즘을 지정하여 다음과 같이 생성합니다. 
(ca.ke에 대한 비밀번호를 입력해야 하는데, CA 비밀키 생성 시 입력한 "gRPC"를 입력합니다.)
```shell
$ openssl req -new -x509 -sha256 -days 730 -key ca.key -out ca.crt
Enter pass phrase for ca.key:gRPC
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Country Name (2 letter code) [AU]:KR
State or Province Name (full name) [Some-State]:Seoul
Locality Name (eg, city) []:
Organization Name (eg, company) [Internet Widgits Pty Ltd]:gRPC Ltd
Organizational Unit Name (eg, section) []:
Common Name (e.g. server FQDN or YOUR name) []:Self Signed CA
Email Address []:
```
※ CA 인증서의 Common Name는 서버 인증서와는 달리 도메인이나 호스트명일 필요는 없습니다.

이제 다음과 같은 명령을 통해 CA 인증서를 확인해 봅니다.
```shell
$ openssl x509 -noout -text -in ca.crt
Certificate:
    Data:
        Version: 3 (0x2)
        Serial Number:
            52:1e:b4:9a:30:f0:9d:b9:f3:40:1a:4b:dc:0a:71:16:5e:f5:ba:15
        Signature Algorithm: sha256WithRSAEncryption
        Issuer: C = KR, ST = Seoul, O = gRPC Ltd, CN = Self Signed CA
        ...
```

## 3. 서버 인증서(Certificate) 생성
서버 인증서를 생성하는데, 이전 방식과 달리 CA를 통해 전자서명된 인증서를 생성해야 합니다. 
이를 위해 우선 CSR(Certificate Signing Request)를 생성하고, CA 비밀키를 사용해 전자서명을 한 서버 인증서를 생성합니다.

우선, 다음과 같이 CSR을 생성합니다.
```shell
$ openssl req -new -sha256 -key server.key -out server.csr
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Country Name (2 letter code) [AU]:KR
State or Province Name (full name) [Some-State]:Seoul
Locality Name (eg, city) []:
Organization Name (eg, company) [Internet Widgits Pty Ltd]:gRPC Ltd
Organizational Unit Name (eg, section) []:
Common Name (e.g. server FQDN or YOUR name) []:localhost
Email Address []:webmaster@localhost

Please enter the following 'extra' attributes
to be sent with your certificate request
A challenge password []:
An optional company name []:
```
### ※ Common Name은 반드시 접속하는 호스트명으로 'localhost'를 입력해야 하며, 추가적으로 입력하는 challenge는 빈값으로 지정합니다.

이제 인증기관(CA)에서 다음과 같이 요청된 CSR에 대하여 인증서를 생성합니다.
```shell
$ openssl x509 -req -days 3650 -sha256 -in server.csr -CA ca.crt -CAkey ca.key -set_serial 1 -out server.crt
Signature ok
subject=C = KR, ST = Seoul, O = gRPC Ltd, CN = localhost, emailAddress = webmaster@localhost
Getting CA Private Key
Enter pass phrase for ca.key:gRPC
```
이 때 CA 비밀키에 대한 비밀번호를 입력해야 합니다.

## 4. 클라이언트 인증서(Certificate) 생성
다음으로 클라이언트 인증서도 생성하는데 CSR을 사용해 서버 인증서 발급과 동일하며, 
다음과 같은 명령들을 사용합니다.

```shell
$ openssl genrsa -out client.key 2048
$ openssl req -new -key client.key -out client.csr
$ openssl x509 -req -days 3650 -sha256 -in client.csr -CA ca.crt -CAkey ca.key -set_serial 2 -out client.crt
```
※ 입력 값들은 서버 인증서와 같은 값들을 사용하며, Common Name도 동일하게 "localhost"를 입력합니다.


## 기타
Java에서는 키 저장소(`.pem`)를 사용해야 하는데, 다음과 같은 명령을 통해 변환이 가능합니다.

```shell
$ openssl pkcs8 -topk8 -inform pem -in server.key -outform pem -nocrypt -out server.pem
$ openssl pkcs8 -topk8 -inform pem -in client.key -outform pem -nocrypt -out client.pem
```