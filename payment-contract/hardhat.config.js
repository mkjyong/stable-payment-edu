require("@nomicfoundation/hardhat-toolbox");

/**
 * Hardhat 기본 설정
 * - Solidity 0.8.24
 * - 로컬 hardhat 네트워크 사용
 */
module.exports = {
  solidity: {
    version: "0.8.24",
    settings: {
      optimizer: { enabled: true, runs: 200 }
    }
  },
  networks: {
    hardhat: {}
  }
};



