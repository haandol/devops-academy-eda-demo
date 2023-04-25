# Day-4 Hand-on Runbook

## 목표

- 개발환경을 설정합니다
- ECS 클러스터에 서비스들을 배포합니다

## EventEngine 설정

- 아래 링크의 내용대로 이벤트 엔진 접속을 진행합니다
  - https://catalog.us-east-1.prod.workshops.aws/workshops/9c0aa9ab-90a9-44a6-abe1-8dff360ae428/ko-KR/20-preq/200-event-engine

## AWS Cloud9 설정

#### AWS Cloud9으로 IDE 구성

1. [AWS Cloud9 콘솔창](https://console.aws.amazon.com/cloud9) 에 접속한 후, **Create environment** 버튼을 클릭합니다.

2. IDE 이름을 적은 후, Next step을 클릭합니다. 본 실습에서는 **devops-day4** 로 입력합니다.

3. 인스턴스 타입(instance type)을 **m5.large** 로 선택합니다.

4. 네트워크 설정 메뉴(networking settings)에서 **VPC 설정**을 클릭합니다. 기본 VPC 대신 **MskStack/Vpc** 으로 표시된 VPC 를 선택합니다. Subnet 은 **PublicSubnet1** 을 선택합니다

   > VPC 설정을 하지 않아도 핸즈온 진행에는 문제가 없습니다. 다만, Cloud9 에서 MSK 에 접근할 수 없으므로 카프카에 쌓인 메시지를 Kafka-UI 툴로 확인할 수 없습니다

![C9 network settings](/img/c9-network.png)

5. 하단의 **Create**를 클릭하여 생성합니다.

6. 생성이 완료되면 아래와 같은 화면이 나타납니다.

<img src="https://static.us-east-1.prod.workshops.aws/public/e7ab9b91-502d-4ada-84e2-5c8b92d8f791/static/images/30-setting/aws_cloud9_01.png" />

### IDE(AWS Cloud9 인스턴스)에 IAM Role 부여

AWS Cloud9 환경은 EC2 인스턴스로 구동됩니다. 따라서 EC2 콘솔에서 AWS Cloud9 인스턴스에 방금 생성한 IAM Role을 부여합니다.

1. [여기](https://console.aws.amazon.com/ec2/v2/home?#Instances:sort=desc:launchTime) 를 클릭하여 EC2 인스턴스 페이지에 접속합니다.

2. 해당 인스턴스를 선택 후, **Actions > Security > Modify IAM Role**을 클릭합니다.

3. IAM Role 에서 `Day4DemoAdminInstanceProfile`을 선택한 후, **Update IAM role** 버튼을 클릭합니다.

### IDE에서 IAM 설정 업데이트

AWS Cloud9의 경우, IAM credentials를 동적으로 관리합니다. 해당 credentials는 **실습상 리소스 프로비저닝시 제약이 있기 때문에** 이를 비활성화하고 IAM Role을 붙입니다.

1. AWS Cloud9 콘솔창에서 생성한 IDE로 다시 접속한 후, 우측 상단에 기어 아이콘을 클릭한 후, 사이드 바에서 **AWS SETTINGS**를 클릭합니다.

2. **Credentials** 항목에서 **AWS managed temporary credentials** 설정을 비활성화합니다.

3. Preference tab을 종료합니다.

<img src="https://static.us-east-1.prod.workshops.aws/public/e7ab9b91-502d-4ada-84e2-5c8b92d8f791/static/images/30-setting/aws_cloud9_05.png" />

4. **Temporary credentials** 이 없는지 확실히 하기 위해 기존의 자격 증명 파일도 제거합니다.

```bash
rm -vf ${HOME}/.aws/credentials
```

5. **GetCallerIdentity CLI** 명령어를 통해, Cloud9 IDE가 올바른 IAM Role을 사용하고 있는지 확인하세요. **결과 값이 나오면** 올바르게 설정된 것입니다.

```bash
aws sts get-caller-identity --query Arn | grep day4DemoAdminRole
```

### 추가설정

- 기본 리전을 한국(`ap-northeast-2`)으로 설정합니다

```bash
aws configure set default.region ap-northeast-2
```

- 원활한 도커 빌드를 위해 Cloud9 의 볼륨 크기를 조정합니다

```bash
wget https://gist.githubusercontent.com/haandol/45c1edfd1e3bf6f88655e655f161463d/raw/a95d52544e7e1a74ee720068b5281b23c4932df3/resize.sh
sh resize.sh
df -h
```

## 서비스 컨테이너 이미지 준비

### 코드 다운로드

```bash
git clone https://github.com/haandol/devops-academy-eda-demo.git ~/environment/devops-academy-eda-demo
cd ~/environment/devops-academy-eda-demo/app
```

### TaskCLI 설치

- 본 핸즈온에서는 GNU Make 를 사용하기 쉽게 만든 [Task](https://taskfile.dev/) 빌드툴을 사용합니다

```bash
npm install -g @go-task/cli
```

### ECR 레포지토리 생성

- 각 서비스별로 ECR 레포지토리(trip, car, hotel, flight) 를 AWSCLI 로 생성합니다

```bash
task create-repo
```

### 컨테이너 이미지 빌드

- [Dockerfile](/app/Dockerfile) 을 통해 각 서비스 컨테이너 이미지를 빌드합니다

```bash
task build-all
```

- 컨테이너 이미지가 잘 생성되었는지 확인합니다

```bash
docker images
```

### 컨테이너 이미지를 레포지토리로 푸시

- 빌드한 이미지들을 ECR 레포지토리로 푸시합니다

```bash
task push-all
```

- [ECR 콘솔](https://ap-northeast-2.console.aws.amazon.com/ecr/repositories?region=ap-northeast-2) 로 이동해 푸시한 이미지를 확인합니다

## CDK 로 MSK 클러스터 확인

- 실습시간 단축을 위해 MSK 클러스터는 각 계정에 미리 배포되어 있습니다
- 추후 개인계정에 MSK 를 배포해보실 분들은, 아래 레포지토리를 이용해 배포해볼 수 있습니다
  - https://github.com/haandol/cdk-msk-example

## CDK 로 ECS 클러스터 설정

### Infra 폴더로 이동

```bash
cd ~/environment/devops-academy-eda-demo/infra
```

### 의존성 설치

```bash
npm i -g aws-cdk@2.76.0 --force
npm i
```

### 설정파일 수정

- ECS 클러스터가 배포될 VPC, 각 서비스가 연결될 MSK 에 대한 설정을 수정합니다
- [config/dev.toml](/infra/config/dev.toml) 파일을 열고 아래 명령어들을 통해 비어있는 필드들의 값을 추가합니다

#### aws.account

```bash
aws sts get-caller-identity --query 'Account' --output text
```

#### vpc.id

```bash
aws ec2 describe-vpcs --filters 'Name=tag:namespace,Values=day4demo' --query 'Vpcs[0].VpcId' --output text
```

#### msk.seeds

```bash
export CLUSTER_ARN=$(aws cloudformation list-exports --query "Exports[?Name=='mskClusterArn'].Value" --output text)
echo $CLUSTER_ARN

aws kafka get-bootstrap-brokers --cluster-arn $CLUSTER_ARN --query 'BootstrapBrokerStringTls' --output text
```

#### msk.securityGroupId

```bash
aws cloudformation list-exports --query "Exports[?Name=='mskSecurityGroup'].Value" --output text
```

#### service.common.tag

```bash
git rev-parse --short=10 HEAD
```

### 설정파일 복사

- 내용을 수정한 뒤, [config/dev.toml](/infra/config/dev.toml) 파일을 인프라 루트에 `.toml` 파일로 복사합니다

```bash
cp ./config/dev.toml ./.toml
```

### 서비스 설정 확인

- ECS 태스크에서 사용할 이미지 정보를 수정합니다
- ECR 이미지 tag 는 _git commit hash_ 값을 사용하며, 이를 통해 시스템을 항상 코드로 부터 재현가능(reproducible)하도록 합니다

## 인프라 배포

### CDK 배포

```bash
cd ~/environment/devops-academy-eda-demo/infra
cdk bootstrap
cdk deploy "*" --concurrency 4 --require-approval never
```

### ECS 클러스터 확인

- [ECS 웹 콘솔](https://ap-northeast-2.console.aws.amazon.com/ecs/v2/clusters/DevOpsDemoDev/services?region=ap-northeast-2)에 방문하여 4개의 서비스가 모두 실행중인지 확인합니다

![ECS Cluster](/img/ecs-cluster.png)

## 서비스 테스트

### HTTPie 설치

- [Httpie](https://httpie.io/) 는 cURL 을 사용하기 쉽게 만든 터미널 기반 도구입니다

```bash
pip3 install httpie
```

### Application LoadBalancer DNS 주소 가져오기

```bash
export ALB=$(aws cloudformation list-exports --query "Exports[?Name=='DevOpsDemoDevLoadBalancerDNS'].Value" --output text)
echo $ALB
```

### 여행예약 요청하기

- 모든 요청에는 인증용 헤더 `x-auth-token:aws-devops` 가 필요합니다

```bash
http post $ALB/v1/trips/ x-auth-token:aws-devops tripId=myTrip1
```

### 생성된 데이터 확인하기

- [DynamoDB 웹 콘솔](https://ap-northeast-2.console.aws.amazon.com/dynamodbv2/home?region=ap-northeast-2#table?name=trip) 로 이동합니다
- Explore Table Item 버튼 클릭하여 생성된 레코드 확인합니다

![Dynamodb Table](/img/ddb-table.png)

![Dynamodb Items](/img/ddb-items.png)

### 생성된 여행 목록 가져오기

- _myTrip1_ 여행이 **Reserved** status 로 반환되는 것을 확인합니다

```bash
http get $ALB/v1/trips/ x-auth-token:aws-devops
```

## X-Ray 로 트레이스 확인

- [Cloudwatch Traces](https://ap-northeast-2.console.aws.amazon.com/cloudwatch/home?region=ap-northeast-2#xray:traces/query) 메뉴로 이동합니다
- 쿼리창에 `http.method = "POST"` 입력후, **Run Query** 버튼을 클릭합니다
- 하단 검색결과 Traces 에서 검색된 첫번째 트레이스 아이디를 클릭합니다
- 서비스 맵을 확인합니다

![X-Ray Traces](/img/xray-trace.png)

![X-Ray Service Map](/img/xray-service-map.png)

---

# Hands-on Lab 2

## 목표

- MSK 에서 처리중인 메시지들을 확인합니다
- 이벤트기반 서비스 장애상황시 대응방법을 간단히 알아봅니다
- 이벤트기반 서비스 구현시 고려할 내용을 알아봅니다

## 카프카 메시지 확인

> 카프카 메시지 확인 실습은 Cloud9 생성시 VPC 설정을 하신 분들만 진행할 수 있습니다

### 내 컴퓨터의 IP 주소 확인하기

- Cloud9 이 아니라 현재 워크샵을 진행중인 개인 컴퓨터의 IP 주소를 확인합니다
- 브라우저에서 [ifconfig.me](https://ifconfig.me/) 로 접근해서 확인하거나, 아래의 명령어로 확인합니다

```bash
curl -s ifconfig.me
```

### 내 컴퓨터에서 Cloud9 인스턴스의 8080 포트로 접근을 허용해주기

> 아이피 주소 앞에 `--` 입력하기

```bash
cd ~/environment/devops-academy-eda-demo/infra
task allow-ingress -- [위에서 확인 한 아이피주소]
# e.g. task allow-ingress -- 39.115.52.111
```

### docker-compose 설치

```bash
pip3 install docker-compose
```

### Kakfa-UI 실행

```bash
cd ~/environment/devops-academy-eda-demo/infra

export CLUSTER_ARN=$(aws cloudformation list-exports --query "Exports[?Name=='mskClusterArn'].Value" --output text)
echo $CLUSTER_ARN

echo KAFKA_CLUSTERS_0_BOOTSTRAPSERVERS=$(aws kafka get-bootstrap-brokers --cluster-arn $CLUSTER_ARN --query 'BootstrapBrokerStringTls' --output text) > .env
docker-compose up -d
```

### Kafka-UI 접속

- 아래 명령어로 Cloud9 이 구동중인 EC2 주소를 확인합니다

```bash
http get http://169.254.169.254/latest/meta-data/public-ipv4
```

- 웹 브라우저에서 새 탭을 열고 해당 주소의 8080 포트로 접속합니다

![Kafka UI](/img/kafka-ui.png)

- 사이드바의 Topics 메뉴 클릭하여 토픽 및 메시지 확인할 수 있습니다

## Hotel 서비스가 다운되었을 때 메시지 흐름 확인하기

### Hotel 서비스를 일시적으로 다운

```bash
aws ecs update-service --service DevOpsDemoDevhotel --cluster DevOpsDemoDev --desired-count 0
aws ecs describe-services --cluster DevOpsDemoDev --service DevOpsDemoDevhotel --query services[0].runningCount
```

### 여행예약 요청

```bash
http post $ALB/v1/trips/ x-auth-token:aws-devops tripId=myTrip2
```

### 생성된 여행 목록 가져오기

- status 가 **Initialized** 상태에서 **Reserved** 상태로 변경되지 않는것 확인합니다

```bash
http get $ALB/v1/trips/ x-auth-token:aws-devops
```

### 데이터베이스 부킹정보 확인

- [Dynamodb 웹 콘솔](https://ap-northeast-2.console.aws.amazon.com/dynamodbv2/home?region=ap-northeast-2#item-explorer?table=trip&maximize=true) 에 접속합니다
- **PK** 필드에 `TRIP#myTrip2` 입력후 **Run** 버튼 클릭합니다
- **myTrip1** 과 달리 hotel, flight 부킹 정보가 생성되지 않았음을 확인할 수 있습니다

![Dynamodb Query Trip2](/img/ddb-query-trip2.png)

### 카프카 메시지 확인

- Kakfa UI 의 토픽 메뉴 에서 메시지 개수를 보면, hotel, flight 메시지 개수가 trip, car 에 비해 1개 작음을 확인할 수 있습니다

![Kafka Topic Messages](/img/kafka-topic-messages.png)

### Tracing 정보 확인

- [X-Ray 의 트레이스 메뉴](https://ap-northeast-2.console.aws.amazon.com/cloudwatch/home?region=ap-northeast-2#xray:traces/)를 통해 트레이스 정보를 확인합니다
- 쿼리창에 `http.method = "POST"` 입력후, **Run Query** 버튼 클릭합니다
- 서비스 맵에서도 hotel, flight 서비스가 보이지 않는 것 확인합니다

## Hotel 서비스가 복구되었을 때 메시지 흐름 확인하기

### Hotel 서비스를 복구

```bash
aws ecs update-service --service DevOpsDemoDevhotel --cluster DevOpsDemoDev --desired-count 1
watch -n 5 aws ecs describe-services --cluster DevOpsDemoDev --service DevOpsDemoDevhotel --query services[0].runningCount
```

### 생성된 여행 목록 가져오기

- 모든 메시지가 **Reserved** 상태로 변경되었는지 확인합니다

```bash
http get $ALB/v1/trips/ x-auth-token:aws-devops
```

## 서비스에 장애가 발생시 확인

### Hotel 서비스에 장애 주입

- 해당 요청으로 장애를 주입하거나, 주입된 장애를 해제할 수 있습니다
  - `:=` 에 주의. httpie 의 raw json input 인디케이터입니다

```bash
http put $ALB/v1/trips/hotels/error/ x-auth-token:aws-devops flag:=true
```

### Hotel 장애주입 여부 확인

```bash
http get $ALB/v1/trips/hotels/error/ x-auth-token:aws-devops
```

### 장애 코드

```golang
/* service/hotel.go */

	// 데이터베이스에 부킹정보를 저장합니다
	booking, err := s.hotelRepository.Book(ctx, evt.Body.TripID)
	if err != nil {
		logger.Errorw("Failed to book hotel", "err", err)
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return err
	}
	span.SetAttributes(
		o11y.AttrString("booking", fmt.Sprintf("%v", booking)),
	)

	// 장애 주입시 에러를 리턴합니다
	if s.ErrorFlag {
		logger.Errorw("Error injection", "err", ErrErrorInjection, "booking", booking)
		span.RecordError(ErrErrorInjection)
		span.SetStatus(o11y.GetStatus(ErrErrorInjection))
		return ErrErrorInjection
	}

	// 카프카에 호텔부킹완료 이벤트를 발행합니다
	if err := s.hotelProducer.PublishHotelBooked(ctx, &booking); err != nil {
		logger.Errorw("Failed to publish HotelBooked", "booking", booking, "err", err)
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return err
	}
```

### 여행예약 요청 및 확인

```bash
http post $ALB/v1/trips/ x-auth-token:aws-devops tripId=myTrip3
http get $ALB/v1/trips/ x-auth-token:aws-devops
```

### 트레이스 확인

- [X-Ray 의 트레이스](https://ap-northeast-1.console.aws.amazon.com/cloudwatch/home?region=ap-northeast-1#xray:service-map/map) 를 통해 장애 발생한 부분 및 상세내용을 확인합니다

### Kafka UI 확인

- 카프카 UI 좌측 사이드바의 Topics 메뉴로 이동합니다
- flight, trip 토픽에 메시지가 1개씩 부족한것 확인할 수 있습니다

### 장애 해제

```bash
http put $ALB/v1/trips/hotels/error/ x-auth-token:aws-devops flag:=false
http get $ALB/v1/trips/hotels/error/ x-auth-token:aws-devops
```

### 요청 재시도

- **myTrip3** 는 **Initialized** 에서 변경되지 않는 것을 확인할 수 있습니다
- **이 경우, 서비스 장애가 해결되어도 myTrip3 요청은 복구할 수 없습니다**

```bash
http post $ALB/v1/trips/ x-auth-token:aws-devops tripId=myTrip3 # 40300 error
http get $ALB/v1/trips/ x-auth-token:aws-devops
```

---

## 리소스 정리

```bash
cd ~/environment/devops-academy-eda-demo/infra
cdk destroy "*"
```
