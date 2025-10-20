# 교육용 결제·급여 컨트랙트 (payment-contract)

본 디렉터리는 1일 실습용 최소 예제입니다. 목표는 스테이블 코인 기반의 간단한 결제(ShopRegistry)와 급여(PayrollBasic) 흐름을 구현·이해하는 것입니다. 하드햇 환경에서 바로 배포·호출·이벤트 관찰까지 진행합니다.

---

## 1) 학습 목표
- 스테이블 코인 결제 흐름 이해: approve → transferFrom → 이벤트
- 온체인 권한 제어 이해: onlyOwner(modifier) 패턴
- 상태 저장·조회: mapping, struct, event, require
- 급여 기본형: set → claim → payAll 루틴, 재입금 방지(상태 0화)

---

## 2) 구성(2~3 컨트랙트)
- ShopRegistry: 상점 등록/한도/결제/영수증 이벤트
- PayrollBasic: 급여 등록/개별 클레임/일괄 지급/증빙 이벤트
- StableTokenMock: 테스트용 ERC20(민트 기능 포함) — 실습 편의 제공용

비고: 오라클·Merkle·Proxy 등 고급 요소는 제외(교육 단순화).

---

## 3) 사용자 스토리(요구사항)
1. 결제자(Payer)는 `ShopRegistry`가 보유한 스테이블 토큰에 대해 `approve` 후 `pay(shopId, amount)`를 호출하면 상점 지갑으로 토큰이 전송된다.
2. 운영자(Owner)는 상점을 등록하고(`registerShop`) 월 한도를 관리한다(`setLimit`).
3. HR(Owner)은 월 급여를 등록한다(`setSalary`). 직원은 스스로 급여를 수령한다(`claim`). 필요 시 한 번에 지급한다(`payAll`).
4. 모든 결제/급여 지급은 이벤트를 발행하여 프론트엔드가 Toast로 표시한다.

---

## 4) 기능 명세

### 4.1 ShopRegistry
- 상태
  - `stableToken`: 결제에 사용할 ERC20 주소
  - `shops[bytes32 => Shop]`: 상점 정보 저장
  - `Shop { limit, spent, merchant }`: 월 한도, 누적 결제액, 상점 수령 지갑
- 함수
  - `registerShop(id, limit, merchant) onlyOwner`
  - `setLimit(id, newLimit) onlyOwner`
  - `pay(id, amount)` — `transferFrom(payer → merchant)` 수행, `spent` 증가, `Receipt` 이벤트
  - `getShop(id) view returns(limit, spent, merchant)`
- 이벤트
  - `Receipt(payer, shopId, amount)`
- 제약
  - 상점 존재 여부 확인, 한도 초과 방지, 금액>0, merchant != 0x0

### 4.2 PayrollBasic
- 상태
  - `stableToken`: 급여 지급에 사용할 ERC20 주소
  - `salaries[address => uint256]`: 직원별 급여
  - `employees[]`, `isEmployee[address => bool]`: 일괄지급 순회용(간단 구현)
- 함수
  - `setSalary(employee, amount) onlyOwner`
  - `claim()` — 본인 급여 수령, 재입금 방지를 위해 0으로 세팅 후 전송
  - `payAll() onlyOwner` — 모든 직원에 대해 미지급 급여 일괄 전송
  - `getSalary(employee) view returns(uint256)`
- 이벤트
  - `ProofOfPayment(employee, amount)`
- 제약
  - 금액>0, 잔액 부족 시 전송 실패, 일괄 지급은 가스 비용 고려(교육용 소규모 데이터 가정)

---

## 5) 설계·문법 포인트
- Pragma: `^0.8.24` — 최신 안전 산술 내장, 외부 SafeMath 불필요
- Modifier: `onlyOwner` — 관리자 전용 함수 보호
- Mapping/Struct: 상태 관리에 핵심 자료구조로 사용
- Event: 프론트에서 결제/급여 결과를 실시간 수신해 UI에 노출
- transferFrom: 결제는 `approve → pay` 2단계로 구현(보안·UX 트레이드오프 체험)
- 재입금 방지(Idempotency): `claim()`에서 상태를 0으로 만든 뒤 전송(Checks-Effects-Interactions 순서)

---

## 6) 실행 방법
1. 설치
   - Node 20 권장
   - 저장소 루트에서:
     - `cd payment-contract`
     - `pnpm i` 또는 `npm i`
2. 컴파일: `npx hardhat compile`
3. 로컬 배포: `npx hardhat run scripts/deploy.js`
4. 상호작용(예):
   - StableTokenMock 민트 → `approve` → `ShopRegistry.pay`
   - `PayrollBasic.setSalary` → `claim` 또는 `payAll`

---

## 7) 테스트 아이디어(선택)
- 결제 한도 초과 시 revert
- 미등록 상점 결제 시 revert
- 급여 0원일 때 `claim` revert
- `payAll` 후 모든 급여 0으로 재설정 확인



