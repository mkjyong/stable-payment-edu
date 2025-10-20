// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";

/**
 * StableTokenMock
 * - 테스트 스테이블 코인(ERC20)
 * - 문법 포인트: 상속, constructor, mint 함수(교육 편의를 위해 anyone mint)
 */
contract StableTokenMock is ERC20 {
    constructor() ERC20("Stable Mock USD", "sUSD") {}

    function mint(address to, uint256 amount) external {
        _mint(to, amount);
    }
}


