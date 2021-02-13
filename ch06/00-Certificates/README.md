# RSA 인증서 만들기

SSL/TLS 적용을 위해 다음과 같은 단계를 거쳐 PKI(Public Key Infrastructure) 인증서를 만듭니다. 

## 1. 비밀키(Private Key) 생성
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

## 2. 공개키(Public Key) 생성
이제 인증서(certficate)에 해당되는 공개키를 만듭니다. 
원래는 인증기관(certificate authority, CA)을 통해 발급 받아야 하지만, 
테스트를 위해 자체 서명 인증서(self-signed certificate)를 사용합니다.

인증서를 만들기 위해 아래와 같이 실행하고, 관련된 값을 입력합니다. 
다른 값들은 임의로 입력해도 되지만, "Common Name"은 인증서와 관련된 도메인, 호스트명 또는 IP를 입력해야 합니다.
해당 값과 접속하는 도메인, 호스트명 또는 IP가 다른 경우 경고를 내기 때문입니다.

본 예제들의 경우는 접속 시에 "localhost"를 사용하므로 반드시 "localhost"로 입력해 주셔야 합니다.

```shell
$ openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
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
```
- `-x509` : 인증서 형식을 지정하는 옵션으로 일반적으로 인터넷에서 많이 사용되는 X.509 지정
- `-sha256` : 보안 해쉬알고리즘 지정(SHA-256)
- `-key server.key` : 사용한 비밀키 파일 지정
- `-out server.crt` : 생성할 인증서/공개키 파일 지정
- `-days 3650` : 인증서 유효기관 (10년)

### ※ Common Name은 반드시 접속하는 호스트명으로 'localhost'를 입력해 주세요.

## 기타
Java에서는 키 저장소(`.pem`)를 사용해야 하는데, 다음과 같은 명령을 통해 변환이 가능합니다.

```shell
$ openssl pkcs8 -topk8 -inform pem -in server.key -outform pem -nocrypt -out server.pem
```