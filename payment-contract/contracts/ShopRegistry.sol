// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

/**
 * ShopRegistry
 * - 교육용 결제 컨트랙트
 * - 문법 포인트: pragma, import, contract 상속(Ownable), struct/mapping, event, require, transferFrom
 */
contract ShopRegistry is Ownable {
    /// 결제에 사용할 스테이블 토큰 주소
    IERC20 public immutable stableToken;

    /// 상점 정보 구조체: 월 한도(limit), 누적 사용액(spent), 수령 지갑(merchant)
    struct Shop {
        uint256 limit;
        uint256 spent;
        address merchant;
    }

    /// shopId(bytes32) → Shop 매핑: 상태 저장의 핵심 자료구조
    mapping(bytes32 => Shop) private shops;

    /// 결제 영수증 이벤트: 프론트엔드 토스트 표시에 사용
    event Receipt(address indexed payer, bytes32 indexed shopId, uint256 amount);

    /**
     * 생성자(Constructor)
     * - immutable: 배포 후 변경 불가, 가스 최적화에 도움
     */
    constructor(IERC20 stableToken_) Ownable(msg.sender) {
        require(address(stableToken_) != address(0), "INVALID_TOKEN");
        stableToken = stableToken_;
    }

    /** 상점 등록: 관리자만 가능 */
    function registerShop(bytes32 shopId, uint256 limit, address merchant) external onlyOwner {
        require(merchant != address(0), "INVALID_MERCHANT");
        Shop storage s = shops[shopId];
        require(s.merchant == address(0), "ALREADY_EXISTS");
        s.limit = limit;
        s.merchant = merchant;
        // spent 초기값은 0
    }

    /** 한도 변경: 관리자만 가능 */
    function setLimit(bytes32 shopId, uint256 newLimit) external onlyOwner {
        Shop storage s = _requireShop(shopId);
        s.limit = newLimit;
    }

    /** 결제: payer → merchant로 transferFrom 수행 */
    function pay(bytes32 shopId, uint256 amount) external {
        require(amount > 0, "INVALID_AMOUNT");
        Shop storage s = _requireShop(shopId);

        uint256 newSpent = s.spent + amount;
        require(newSpent <= s.limit, "LIMIT_EXCEEDED");

        s.spent = newSpent;
        // 결제 토큰 전송: 사용자는 사전에 approve(컨트랙트 주소, 금액) 필요
        bool ok = stableToken.transferFrom(msg.sender, s.merchant, amount);
        require(ok, "TRANSFER_FAIL");

        emit Receipt(msg.sender, shopId, amount);
    }

    /** 상점 정보 조회: view 함수 */
    function getShop(bytes32 shopId) external view returns (uint256 limit, uint256 spent, address merchant) {
        Shop storage s = _requireShop(shopId);
        return (s.limit, s.spent, s.merchant);
    }

    /** 내부 유틸: 상점 존재 보장 */
    function _requireShop(bytes32 shopId) private view returns (Shop storage) {
        Shop storage s = shops[shopId];
        require(s.merchant != address(0), "SHOP_NOT_FOUND");
        return s;
    }
}


