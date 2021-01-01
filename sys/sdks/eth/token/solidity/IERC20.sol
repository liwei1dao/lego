pragma solidity ^0.7.0;

interface IERC20 {
    /**
     * @dev 返回代币总量
     */
    function totalSupply() external view returns (uint256);

    /**
     * @dev 查询用户的代币数
     */
    function balanceOf(address account) external view returns (uint256);

    /**
     * @dev 从调用者钱包向另一个地址发送余额
     */
    function transfer(address recipient, uint256 amount) external returns (bool);

    /**
     * @dev 允许_spender从你的账户转出_value的余额，调用多次会覆盖可用量。某些DEX功能需要此功能
     */
    function allowance(address owner, address spender) external view returns (uint256);

    /**
     * @dev 允许_spender从你的账户转出_value余额
     */
    function approve(address spender, uint256 amount) external returns (bool);

    /**
     * @dev 从一个地址向另一个地址发送余额
     */
    function transferFrom(address sender, address recipient, uint256 amount) external returns (bool);

    /**
     * @dev 交易事件
     */
    event Transfer(address indexed from, address indexed to, uint256 value);

    /**
     * @dev approve 事件
     */
    event Approval(address indexed owner, address indexed spender, uint256 value);
}