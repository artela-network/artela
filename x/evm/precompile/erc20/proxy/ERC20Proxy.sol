// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract ERC20Proxy2 {
    // Address of the precompiled ERC20 contract
    address public constant ERC20_PRECOMPILED_ADDRESS = address(0x0000000000000000000000000000000000001234);

    /**
     * @dev Fallback function that forwards all unknown calls to the precompiled ERC20 contract.
     */
    fallback() external payable {
        _delegate(ERC20_PRECOMPILED_ADDRESS);
    }

    /**
     * @dev Receive function to accept plain Ether transfers.
     */
    receive() external payable {}

    /**
     * @dev Delegates the request to the target address and returns the result to the caller.
     * @param target The address of the target contract (i.e., the precompiled ERC20 contract).
     */
    function _delegate(address target) internal {
        (bool success, bytes memory returndata) = target.delegatecall(msg.data);
        assembly {
            switch success
            case 0 { revert(add(returndata, 32), mload(returndata)) } // Revert on failure
            default { return(add(returndata, 32), mload(returndata)) } // Return data on success
        }
    }

    /**
     * @dev Queries the token balance of an account.
     * @param account The address of the account to query.
     * @return balance The token balance of the account.
     */
    function balanceOf(address account) external returns (uint256) {
        (bool success, bytes memory result) = ERC20_PRECOMPILED_ADDRESS.delegatecall(
            abi.encodeWithSignature("balanceOf(address)", account)
        );
        require(success, "ERC20: delegate call failed");
        return abi.decode(result, (uint256));
    }

    /**
     * @dev Transfers tokens to a specified address.
     * @param to The address to receive the tokens.
     * @param amount The amount of tokens to transfer.
     * @return success Whether the transfer was successful.
     */
    function transfer(address to, uint256 amount) external returns (bool) {
        (bool success, bytes memory result) = ERC20_PRECOMPILED_ADDRESS.delegatecall(
            abi.encodeWithSignature("transfer(address,uint256)", to, amount)
        );
        require(success, "ERC20: transfer failed");
        return abi.decode(result, (bool));
    }

    /**
     * @dev Approves a spender to use a specified amount of tokens.
     * @param spender The address authorized to spend the tokens.
     * @param amount The amount of tokens to approve.
     * @return success Whether the approval was successful.
     */
    function approve(address spender, uint256 amount) external returns (bool) {
        (bool success, bytes memory result) = ERC20_PRECOMPILED_ADDRESS.delegatecall(
            abi.encodeWithSignature("approve(address,uint256)", spender, amount)
        );
        require(success, "ERC20: approve failed");
        return abi.decode(result, (bool));
    }

    /**
     * @dev Queries the allowance of a spender for a specific owner.
     * @param owner The address of the token holder.
     * @param spender The address authorized to spend the tokens.
     * @return remaining The remaining amount of tokens the spender can use.
     */
    function allowance(address owner, address spender) external returns (uint256) {
        (bool success, bytes memory result) = ERC20_PRECOMPILED_ADDRESS.delegatecall(
            abi.encodeWithSignature("allowance(address,address)", owner, spender)
        );
        require(success, "ERC20: allowance query failed");
        return abi.decode(result, (uint256));
    }
}
