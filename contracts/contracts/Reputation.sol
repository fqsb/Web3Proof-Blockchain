// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import "@openzeppelin/contracts/access/AccessControl.sol";

contract Reputation is AccessControl {
    bytes32 public constant UPDATER_ROLE = keccak256("UPDATER_ROLE");

    struct Score {
        uint256 total;
        uint256 projectScore;
        uint256 certScore;
        uint256 activityScore;
        uint64 updatedAt;
    }

    mapping(address => Score) private _scores;

    event ScoreUpdated(address indexed user, uint256 total, uint256 grade);

    constructor(address admin) {
        _grantRole(DEFAULT_ADMIN_ROLE, admin);
        _grantRole(UPDATER_ROLE, admin);
    }

    function updateScore(
        address user,
        uint256 projectScore,
        uint256 certScore,
        uint256 activityScore
    ) external onlyRole(UPDATER_ROLE) {
        uint256 total = projectScore + certScore + activityScore;
        require(total <= 1000, "Score exceeds max");

        _scores[user] = Score({
            total: total,
            projectScore: projectScore,
            certScore: certScore,
            activityScore: activityScore,
            updatedAt: uint64(block.timestamp)
        });

        emit ScoreUpdated(user, total, _grade(total));
    }

    function getScore(address user) external view returns (Score memory) {
        return _scores[user];
    }

    function getGrade(address user) external view returns (uint256) {
        return _grade(_scores[user].total);
    }

    function _grade(uint256 total) private pure returns (uint256) {
        if (total >= 800) return uint256(bytes32("A"));
        if (total >= 600) return uint256(bytes32("B"));
        if (total >= 400) return uint256(bytes32("C"));
        return uint256(bytes32("D"));
    }
}
