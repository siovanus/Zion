pragma solidity >=0.7.0 <0.9.0;

/**
 * @dev Interface of the SideChainManager contract
 */

interface ISideChainManager {

    struct SideChain {
        uint64 chainID;
        uint64 router;
        string name;
        bytes CCMCAddress;
        bytes extraInfo;
    }

    function getSideChain(uint64 chainID) external view returns(SideChain memory sidechain);

    function updateFee(uint64 chainID, uint64 viewNum, int fee, bytes calldata signature) external returns (bool success);

    function registerAsset(uint64 chainID, uint64[] calldata AssetMapKey, bytes[] calldata AssetMapValue, uint64[] calldata LockProxyMapKey, bytes[] calldata LockProxyMapValue) external returns (bool success);

    function getFee(uint64 chainID) external view returns (bytes memory);
}
