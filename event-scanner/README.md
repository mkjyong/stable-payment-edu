### Event Scanner (PayrollBasic.ProofOfPayment)

Go 기반 이더리움 이벤트 스캐너로, `PayrollBasic` 컨트랙트의 `ProofOfPayment(address indexed employee, uint256 amount)` 이벤트를 수집합니다.

### 준비물
- **Go 1.22+**
- **RPC 노드**
  - 과거 로그 수집: HTTP RPC (`https://...`)
  - 실시간 구독: WebSocket RPC (`wss://...`)
- **컨트랙트 주소**

### 설치
```bash
cd @event-scanner
go mod tidy
```

### 실행 - 과거 로그 스캔 (HTTP FilterLogs)
```bash
export RPC=https://your-http-endpoint
export CONTRACT_ADDRESS=0xYourContract

# 전체 블록 범위 스캔
go run . -mode http

# 특정 시작 블록부터 스캔
go run . -mode http -from 10000000

# 특정 직원 주소로 필터링
go run . -mode http -employee 0xEmployee
```

### 실행 - 실시간 구독 (WebSocket SubscribeFilterLogs)
```bash
export RPC=wss://your-ws-endpoint
export CONTRACT_ADDRESS=0xYourContract

go run . -mode ws

# 특정 직원 주소로 필터링
go run . -mode ws -employee 0xEmployee
```

### 출력 포맷
각 이벤트는 JSON 한 줄로 출력됩니다.
```json
{"block":123,"txHash":"0x...","logIndex":0,"employee":"0x...","amount":"1000000000000000000"}
```

### 참고
- ABI는 소스 내에 하드코딩되어 있습니다. 필요 시 하드햇 아티팩트 JSON에서 ABI를 읽어와 교체할 수 있습니다.
- HTTP 모드에선 `-from`으로 시작 블록을 지정할 수 있습니다. 끝 블록은 최신 블록입니다.
- WS 모드는 구독 시작 시점 이후의 이벤트만 수신합니다.


