const { ethers } = require("hardhat");

/**
 * 배포 스크립트
 * - StableTokenMock → ShopRegistry → PayrollBasic 순으로 배포
 * - 배포 후 샘플 상점/직원 데이터 셋업 예시 로그 출력
 */
async function main() {
  const [deployer] = await ethers.getSigners();
  console.log("Deployer:", deployer.address);

  const StableTokenMock = await ethers.getContractFactory("StableTokenMock");
  const stable = await StableTokenMock.deploy();
  await stable.waitForDeployment();
  console.log("StableTokenMock:", await stable.getAddress());

  const ShopRegistry = await ethers.getContractFactory("ShopRegistry");
  const shop = await ShopRegistry.deploy(await stable.getAddress());
  await shop.waitForDeployment();
  console.log("ShopRegistry:", await shop.getAddress());

  const PayrollBasic = await ethers.getContractFactory("PayrollBasic");
  const payroll = await PayrollBasic.deploy(await stable.getAddress());
  await payroll.waitForDeployment();
  console.log("PayrollBasic:", await payroll.getAddress());

  // 샘플 데이터(로컬 테스트): 토큰 민트 → approve → 상점 등록 → 급여 등록
  const mintAmt = ethers.parseUnits("10000", 18);
  await (await stable.mint(deployer.address, mintAmt)).wait();

  // 상점 등록
  const shopId = ethers.id("coffee-shop");
  await (await shop.registerShop(shopId, ethers.parseUnits("5000", 18), deployer.address)).wait();

  // 급여 등록(자기 자신 예시)
  await (await payroll.setSalary(deployer.address, ethers.parseUnits("1000", 18))).wait();

  console.log("Setup complete. You can now approve & pay, or claim/payAll.");
}

main().catch((e) => {
  console.error(e);
  process.exit(1);
});



