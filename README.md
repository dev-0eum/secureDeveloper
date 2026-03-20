# Go Secure Coding Practice

보안 코딩 연습을 위한 시작 프로젝트입니다.

처음부터 구조를 예쁘게 나누기보다, `cmd/server/main.go` 하나에 코드를 모아 둔 상태에서
먼저 흐름을 이해하고 직접 분리 기준을 고민할 수 있게 만드는 것이 목적입니다.

지난 과제와 수업에서 설명했던 내용, 그리고 전달된 가이드 문서를 떠올리면서
어떤 기능부터 구현하고 어떤 기준으로 나눌지 스스로 판단해 보세요.

## 프로젝트 목적

- 로그인 흐름을 먼저 이해하기
- 더미 API를 실제 동작 코드로 바꿔 보기
- 게시판과 금융 기능을 단계적으로 채워 보기
- 입력 검증, 인증 확인, 권한 검사, 응답 설계를 직접 고민해 보기
- 구현 후 어떤 코드들을 묶어 리팩터링할지 판단해 보기

## 현재 상태

- 서버 코드는 현재 `cmd/server/main.go` 하나에 들어 있습니다.
- 로그인만 SQLite 조회를 사용합니다.
- 로그인 성공 시 `authorization` 쿠키와 `Authorization` 헤더용 토큰을 함께 사용할 수 있습니다.
- 세션은 DB가 아니라 메모리 맵으로 관리합니다.
- 게시판 API는 고정 더미 응답을 반환합니다.
- 금융 API는 더미 응답을 반환합니다.
- 정적 화면은 SPA 형태로 준비되어 있습니다.

## 기본 계정

- `alice / alice1234`
- `bob / bob1234`
- `charlie / charlie1234`

## 주요 API

인증
- `POST /api/auth/register`
- `POST /api/auth/login`
- `POST /api/auth/logout`
- `POST /api/auth/withdraw`

사용자
- `GET /api/me`

게시판
- `GET /api/posts`
- `GET /api/posts/:id`
- `POST /api/posts`
- `PUT /api/posts/:id`
- `DELETE /api/posts/:id`

금융
- `POST /api/banking/deposit`
- `POST /api/banking/withdraw`
- `POST /api/banking/transfer`

## 참고 파일

- `schema.sql`
- `seed.sql`
- `query_examples.sql`

## 먼저 해볼 작업

1. 회원가입 더미 핸들러를 실제 `INSERT`로 바꿔 보기
2. 게시글 작성, 조회, 수정, 삭제를 실제 SQL 또는 원하는 저장 방식으로 바꿔 보기
3. 입금, 출금, 송금 더미 핸들러를 실제 로직으로 바꿔 보기
4. 실패 조건과 입력 검증 규칙을 정리해 보기
5. 반복되는 인증 확인 코드를 어떻게 줄일지 생각해 보기

## 작업하면서 점검할 질문

- 이 코드는 요청 처리인지, 비즈니스 규칙인지, DB 접근인지?
- 같은 인증 확인 코드가 반복되고 있지 않은가?
- 게시글 수정/삭제 권한 검사는 어디에 두는 것이 자연스러운가?
- 출금과 송금에서 어떤 실패 조건을 먼저 막아야 하는가?
- 어떤 응답은 그대로 내려도 되고, 어떤 응답은 가공이 필요한가?
- 언제부터 `handler`, `service`, `store` 같은 구조로 나누는 것이 좋은가?

## 실행 방법

프로젝트 루트에서 실행합니다.

```powershell
go run ./cmd/server
```
처음 상태로 다시 시작하고 싶으면 `app.db`를 지운 뒤 다시 실행하면 됩니다.

### 작업한 내용의 특징에 대해서 작성해주세요
1. 윈도우 방화벽 이슈로 인한 실행 포트 변경
> 기존 8080 세팅에서 8081로 변경

2. 의존성 라이브러리 설치 및 활용
./requirements에 작성

3. Signup Handler
> main.go:138
> Prepared Statement Query

4. Log Middleware 설정
> ./cmd/server/middleware 에서 관리
> main.go:116 Init (Log Rotation)
> main.go:126 JSONLogger

5. /posts GET, POST 구현
> main.go:378
> main.go:438

6. /posts/:id GET PUT DELETE 구현
> main.go:472
> main.go:528
> main.go:565 // 여기서는 post 작성자인지 확인해주는 작업이 있어야했음.

7. /withdraw 구현
> main.go:226

8. transaction을 이용한 예금 구현
> main.go:281

> deposit을 추가하는 것만 국한해서 볼 때는 굳이 transaction을 할 필요없다. 실패하면 err를 반환하면 되기 때문이다.
> 사실상 transfer 단계에서 A의 balance를 차감하고 B의 balance를 증가시키는 두 개의 행위가 하나의 묶음이 되어 처리하도록 하는 기술이 transaction이지만, 시간적 한계로 인해 이 단계에서 구현하였다.

9. 비슷한 query가 많았어서 DB package를 구성하고자 했다
> ./ext/DB/sqlite 아래에 Select, Insert, Update, Delete별로 func을 구현하여 main에서는 CRUD만 구현하고 QUERY는 호출하기만 하도록 구현하고자 했다.

10. 기능 별 package
> ./svc (service) 아래에 해당 페이지의 기능(로그인, 회원가입, 게시글, 뱅킹 등)별로 분할하여 개발하고자 했다.

11. validate package
> 각 RESTful API를 보면 json 타입 검증(ShouldBindJSON), 토큰 검증, 세션 생성 함수를 중복적으로 사용한다.
> 이를 decoration과 같이 만들어놓으면 각 api에서 처음에 호출만 하면 되도록 설계하고자 했다.
ex> ./validate/bindjson/bind.go