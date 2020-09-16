pragma solidity ^0.7.0;

import "./Context.sol";
import "./ERC20.sol";
import "./ERC20Detailed.sol";
import "./SafeMath.sol";

contract HiToolCoin is Context, ERC20, ERC20Detailed {
    using SafeMath for uint256;
    address _owner;                                     // 合约管理账号
    uint256 public _tokenExchangeRate;                  // xxx 代币 兑换 1 ETH
    address payable public _ethFundDeposit;             // ETH存放地址

    //@func                            构造函数                             
    //@name                            代币名称
    //@symbol                          代币符号
    //@decimals                        代币小数点(建议 18 和eth小数点一致 可以保证eth购买代币无误差问题)
    //@tokenExchangeRate               代币汇率 xxx 代币 兑换 1 ETH (建议 1 保持与eth汇率一直)
    //@totalSupply                     代币总发行量
    //@ethFundDeposit                  代币获取ETH存放地址
    constructor (string memory name, string memory symbol, uint8 decimal,uint256 tokenExchangeRate,uint256 totalSupply,address payable ethFundDeposit) public ERC20Detailed(name, symbol, decimal) {
        _owner = _msgSender();
        _tokenExchangeRate = 1; //1 :代币汇率 1ETH -> 1HTC
        _ethFundDeposit = _msgSender(); //合约ETH 存储地址
        //初始化币，并把所有的币都给部署智能合约的ETH钱包地址
        //100000：代币的总数量
        _mint(_msgSender(), 1000000 * (10 ** uint256(decimal)));
    }

    //@func                             访问修饰符 合约管理账号
    modifier onlyOwner {
        require(
            msg.sender == _owner,
            "Only owner can call this function."
        );
        _;
    }

    //@func                             设置新的所有者地址
    function changeOwner(address payable newFundDeposit) onlyOwner external {
        require (newFundDeposit != address(0x0));
        _ethFundDeposit = newFundDeposit;
    }
 
    //@func                             设置token汇率
    function setTokenExchangeRate(uint256 tokenExchangeRate) onlyOwner external {
        require (tokenExchangeRate != 0);
        require (tokenExchangeRate != _tokenExchangeRate);
        _tokenExchangeRate = tokenExchangeRate;
    }
 
    //@func                             转账ETH 到ETH存放地址
    function transferETH() onlyOwner external {
        require (address(this).balance != 0);
        _ethFundDeposit.transfer(address(this).balance);
    }
 
    //@func                             购买代币
    fallback () payable external {
        require(msg.value != 0,"The purchase amount cannot be 0");
        require(_msgSender() != _owner,"The owner cannot perform this operation");
        uint256 tokens = msg.value.mul(_tokenExchangeRate);
        _transfer(_owner,_msgSender(),tokens);
    }
}