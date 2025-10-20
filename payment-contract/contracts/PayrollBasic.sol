// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

/**
 * PayrollBasic
 * - 교육용 급여 컨트랙트(기본형)
 * - 문법 포인트: mapping, modifier, event, checks-effects-interactions, 배열 순회
 */
contract PayrollBasic is Ownable {
    IERC20 public immutable stableToken;

    /// 직원별 급여 테이블
    mapping(address => uint256) private salaries;
    /// 일괄 지급 순회용 단순 배열(교육용)
    address[] private employees;
    mapping(address => bool) private isEmployee;

    event ProofOfPayment(address indexed employee, uint256 amount);

    constructor(IERC20 stableToken_) Ownable(msg.sender) {
        require(address(stableToken_) != address(0), "INVALID_TOKEN");
        stableToken = stableToken_;
    }

    /** 급여 등록/변경: 관리자만 */
    function setSalary(address employee, uint256 amount) external onlyOwner {
        require(employee != address(0), "INVALID_EMPLOYEE");
        salaries[employee] = amount;
        if (!isEmployee[employee]) {
            isEmployee[employee] = true;
            employees.push(employee);
        }
    }

    /** 개별 수령(클레임): 직원이 스스로 호출 */
    function claim() external {
        uint256 amount = salaries[msg.sender];
        require(amount > 0, "NO_SALARY");

        // Checks-Effects-Interactions: 상태 0화 후 토큰 전송
        salaries[msg.sender] = 0;

        bool ok = stableToken.transfer(msg.sender, amount);
        require(ok, "TRANSFER_FAIL");
        emit ProofOfPayment(msg.sender, amount);
    }

    /** 일괄 지급: 관리자 실행, 소규모 데이터 가정 */
    function payAll() external onlyOwner {
        uint256 len = employees.length;
        for (uint256 i = 0; i < len; i++) {
            address emp = employees[i];
            uint256 amount = salaries[emp];
            if (amount == 0) continue;
            salaries[emp] = 0; // 재입금 방지
            bool ok = stableToken.transfer(emp, amount);
            require(ok, "TRANSFER_FAIL");
            emit ProofOfPayment(emp, amount);
        }
    }

    function getSalary(address employee) external view returns (uint256) {
        return salaries[employee];
    }

    function getEmployees() external view returns (address[] memory) {
        return employees;
    }
}



