## Hardhat 세팅·배포 가이드 (교육용)

본 문서는 교육용으로 Hardhat 설치부터 로컬/테스트넷 배포, 간단 상호작용 예시까지 단계별로 안내합니다. 이 저장소의 `payment-contract` 디렉터리를 기준으로 설명합니다.

---

### 1) 사전 준비
- **Node.js 20 권장** (v18 이상 권장)
- **패키지 관리자**: npm(권장) 또는 pnpm/yarn

---

### 2) 빠른 시작(이 저장소 기준)
- 경로: `/Users/jeonmingyu/dev/personal/edu/stable-payment-edu/payment-contract`

1. 의존성 설치
```bash
cd /Users/jeonmingyu/dev/personal/edu/stable-payment-edu/payment-contract
npm i
```

2. 컴파일
```bash
npx hardhat compile
```

3. 로컬(에페메럴) 네트워크로 원샷 배포
```bash
npm run deploy:local
# 내부적으로: hardhat run scripts/deploy.js --network hardhat
```
- 콘솔 출력에 `StableTokenMock`, `ShopRegistry`, `PayrollBasic` 주소가 순서대로 표시됩니다.
- 배포 스크립트(`scripts/deploy.js`)는 테스트용 민트/상점등록/급여등록 예시를 수행합니다.

---

### 3) 지속형 로컬 노드 사용(주소 재사용·여러 번 상호작용)
1. 노드 실행(터미널 A)
```bash
npx hardhat node
```

2. 배포(터미널 B)
```bash
npx hardhat run scripts/deploy.js --network localhost
```

3. 콘솔 상호작용(터미널 C)
```bash
npx hardhat console --network localhost
```

콘솔 예시(배포 로그의 주소를 사용하세요):
```javascript
const [owner] = await ethers.getSigners();

// 배포된 주소 입력
const stable  = await ethers.getContractAt("StableTokenMock", "0xStable...");
const shop    = await ethers.getContractAt("ShopRegistry",    "0xShop...");
const payroll = await ethers.getContractAt("PayrollBasic",    "0xPayroll...");

// 토큰 민트(교육용)
await (await stable.mint(owner.address, ethers.parseUnits("1000", 18))).wait();

// 결제 흐름: approve → pay
const shopId = ethers.id("coffee-shop"); // deploy.js와 동일 예시 ID
await (await stable.approve(await shop.getAddress(), ethers.parseUnits("100", 18))).wait();
await (await shop.pay(shopId, ethers.parseUnits("100", 18))).wait();

// 급여 흐름: setSalary → claim
await (await payroll.setSalary(owner.address, ethers.parseUnits("50", 18))).wait();
await (await payroll.claim()).wait();
```

참고:
- 본 프로젝트는 Ethers v6 API를 사용합니다(`ethers.parseUnits`, `waitForDeployment()` 등).

---

### 4) 테스트넷 배포(Sepolia 예시)
1. 환경 변수 준비
```bash
npm i dotenv
```
프로젝트 루트(`payment-contract`)에 `.env` 생성:
```env
SEPOLIA_RPC_URL=https://...
PRIVATE_KEY=0x배포자개인키
```

2. `hardhat.config.js`에 네트워크 추가
```javascript
require("@nomicfoundation/hardhat-toolbox");
require("dotenv").config();

module.exports = {
  solidity: { version: "0.8.24", settings: { optimizer: { enabled: true, runs: 200 } } },
  networks: {
    hardhat: {},
    sepolia: {
      url: process.env.SEPOLIA_RPC_URL,
      accounts: [process.env.PRIVATE_KEY],
    },
  },
};
```

3. 테스트넷 배포 실행
```bash
npx hardhat run scripts/deploy.js --network sepolia
```

주의:
- 테스트넷 지갑에 충분한 테스트 ETH가 필요합니다(가스비).

---

### 5) 처음부터 새 Hardhat 프로젝트 만들기(선택)
```bash
mkdir my-hardhat && cd my-hardhat
npm init -y
npm i -D hardhat @nomicfoundation/hardhat-toolbox
npx hardhat   # Create a JavaScript project 선택
```
이후 `contracts/`에 솔리디티 파일, `scripts/deploy.js` 작성 →
```bash
npx hardhat compile
npx hardhat run scripts/deploy.js --network hardhat

# 또는 지속형 노드
npx hardhat node
npx hardhat run scripts/deploy.js --network localhost
```

---

### 6) 자주 겪는 문제(트러블슈팅)
- `HH1: You are not inside a Hardhat project` → `payment-contract` 폴더에서 실행했는지 확인
- Node 버전 문제 → Node 20 권장, LTS 권장
- Ethers v6/BigInt 관련 에러 → v6 API 문법(예: `parseUnits`, `toBigInt`) 사용 확인

---

### 7) 참고 파일
- `hardhat.config.js`: 컴파일러/네트워크 설정
- `scripts/deploy.js`: 배포 및 샘플 셋업(민트/상점등록/급여등록)
- `contracts/*.sol`: 컨트랙트 구현(ShopRegistry, PayrollBasic, StableTokenMock)

---

필요 시 테스트넷 네트워크를 더 추가하거나, 배포된 주소를 `event-scanner`와 연동하는 방법도 확장 안내할 수 있습니다.


