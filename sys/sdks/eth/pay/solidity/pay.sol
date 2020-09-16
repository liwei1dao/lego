pragma solidity ^0.7.0;

contract Account {
    address payable public reciever;
    event Recharge (uint256 value);
    event Flush(address to, uint256 value);

    constructor(address payable _reciever) public {
        reciever = _reciever;
    }

    fallback() external payable{
        emit Recharge(msg.value);
    }

    function flush() public {
        uint256 balance = address(this).balance;
        if (balance == 0){
            return;
        }
        reciever.transfer(balance);
        emit Flush(reciever, balance);
    }
}

contract Wallet {
    address payable public admin;
    mapping(address => bool) public accounts;

    event Create(address);

    constructor() public {
        admin = msg.sender;
    }

    modifier OnlyAdmin {
        require(msg.sender == admin, "403");
        _;
    }

    // 在这里创建新的 CREATE2 账户，保证 CREATE2 的地址参数始终是当前合约
    function create(address payable _to, bytes32 _salt) public OnlyAdmin {
        Account a = new Account{salt: _salt}(_to);
        emit Create(address(a));
    }
}