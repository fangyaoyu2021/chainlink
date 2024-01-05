// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {OptimismCrossDomainForwarder} from "../../../dev/optimism/OptimismCrossDomainForwarder.sol";
import {MockOVMCrossDomainMessenger} from "../../mocks/optimism/MockOVMCrossDomainMessenger.sol";
import {Greeter} from "../../../../../v0.8/tests/Greeter.sol";
import {L2EPTest} from "../L2EPTest.t.sol";

// Use this command from the /contracts directory to run this test file:
//
//  FOUNDRY_PROFILE=l2ep forge test -vvv --match-path ./src/v0.8/l2ep/test/v1_0_0/optimism/OptimismCrossDomainForwarder.t.sol
//
contract OptimismCrossDomainForwarderTest is L2EPTest {
  /// Helper variables
  address internal s_strangerAddr = vm.addr(0x1);
  address internal s_l1OwnerAddr = vm.addr(0x2);

  /// Contracts
  MockOVMCrossDomainMessenger internal s_mockOptimismCrossDomainMessenger;
  OptimismCrossDomainForwarder internal s_optimismCrossDomainForwarder;
  Greeter internal s_greeter;

  /// Events
  event L1OwnershipTransferRequested(address indexed from, address indexed to);
  event L1OwnershipTransferred(address indexed from, address indexed to);

  /// Setup
  function setUp() public {
    // Deploys contracts
    vm.startPrank(s_l1OwnerAddr, s_l1OwnerAddr);
    s_mockOptimismCrossDomainMessenger = new MockOVMCrossDomainMessenger(s_l1OwnerAddr);
    s_optimismCrossDomainForwarder = new OptimismCrossDomainForwarder(
      s_mockOptimismCrossDomainMessenger,
      s_l1OwnerAddr
    );
    s_greeter = new Greeter(address(s_optimismCrossDomainForwarder));
    vm.stopPrank();
  }

  /// @param message - the new greeting message, which will be passed as an argument to Greeter#setGreeting
  /// @return a 2-layer encoding such that decoding the first layer provides the CrossDomainForwarder#forward
  ///         function selector and the corresponding arguments to the forward function, and decoding the
  ///         second layer provides the Greeter#setGreeting function selector and the corresponding
  ///         arguments to the set greeting function (which in this case is the input message)
  function encodeCrossDomainForwardMessage(string memory message) public view returns (bytes memory) {
    return
      abi.encodeWithSelector(
        s_optimismCrossDomainForwarder.forward.selector,
        address(s_greeter),
        abi.encodeWithSelector(s_greeter.setGreeting.selector, message)
      );
  }
}

contract Constructor is OptimismCrossDomainForwarderTest {
  /// @notice it should set the owner correctly
  function test_Owner() public {
    assertEq(s_optimismCrossDomainForwarder.owner(), s_l1OwnerAddr);
  }

  /// @notice it should set the l1Owner correctly
  function test_L1Owner() public {
    assertEq(s_optimismCrossDomainForwarder.l1Owner(), s_l1OwnerAddr);
  }

  /// @notice it should set the crossdomain messenger correctly
  function test_CrossDomainMessenger() public {
    assertEq(s_optimismCrossDomainForwarder.crossDomainMessenger(), address(s_mockOptimismCrossDomainMessenger));
  }

  /// @notice it should set the typeAndVersion correctly
  function test_TypeAndVersion() public {
    assertEq(s_optimismCrossDomainForwarder.typeAndVersion(), "OptimismCrossDomainForwarder 1.0.0");
  }
}

contract Forward is OptimismCrossDomainForwarderTest {
  /// @notice it should not be callable by unknown address
  function test_NotCallableByUnknownAddress() public {
    vm.startPrank(s_strangerAddr, s_strangerAddr);
    vm.expectRevert("Sender is not the L2 messenger");
    s_optimismCrossDomainForwarder.forward(address(s_greeter), abi.encode(""));
    vm.stopPrank();
  }

  /// @notice it should be callable by crossdomain messenger address / L1 owner
  function test_Forward() public {
    // Sets msg.sender and tx.origin
    vm.startPrank(s_strangerAddr, s_strangerAddr);

    // Defines the cross domain message to send
    string memory greeting = "hello";

    // Sends the message
    s_mockOptimismCrossDomainMessenger.sendMessage(
      address(s_optimismCrossDomainForwarder), // target
      encodeCrossDomainForwardMessage(greeting), // message
      0 // gas limit
    );

    // Checks that the greeter got the message
    assertEq(s_greeter.greeting(), greeting);

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }

  /// @notice it should revert when contract call reverts
  function test_ForwardRevert() public {
    // Sets msg.sender and tx.origin
    vm.startPrank(s_strangerAddr, s_strangerAddr);

    // Sends an invalid message
    vm.expectRevert("Invalid greeting length");
    s_mockOptimismCrossDomainMessenger.sendMessage(
      address(s_optimismCrossDomainForwarder), // target
      encodeCrossDomainForwardMessage(""), // message
      0 // gas limit
    );

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }
}

contract TransferL1Ownership is OptimismCrossDomainForwarderTest {
  /// @notice it should not be callable by non-owners
  function test_NotCallableByNonOwners() public {
    vm.startPrank(s_strangerAddr, s_strangerAddr);
    vm.expectRevert("Sender is not the L2 messenger");
    s_optimismCrossDomainForwarder.transferL1Ownership(s_strangerAddr);
    vm.stopPrank();
  }

  /// @notice it should not be callable by L2 owner
  function test_NotCallableByL2Owner() public {
    vm.startPrank(s_l1OwnerAddr, s_l1OwnerAddr);
    assertEq(s_optimismCrossDomainForwarder.owner(), s_l1OwnerAddr);
    vm.expectRevert("Sender is not the L2 messenger");
    s_optimismCrossDomainForwarder.transferL1Ownership(s_strangerAddr);
    vm.stopPrank();
  }

  /// @notice it should be callable by current L1 owner
  function test_CallableByL1Owner() public {
    // Sets msg.sender and tx.origin
    vm.startPrank(s_strangerAddr, s_strangerAddr);

    // Defines the cross domain message to send
    vm.expectEmit(false, false, false, true);
    emit L1OwnershipTransferRequested(s_optimismCrossDomainForwarder.l1Owner(), s_strangerAddr);

    // Sends the message
    s_mockOptimismCrossDomainMessenger.sendMessage(
      address(s_optimismCrossDomainForwarder), // target
      abi.encodeWithSelector(s_optimismCrossDomainForwarder.transferL1Ownership.selector, s_strangerAddr), // message
      0 // gas limit
    );

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }

  /// @notice it should be callable by current L1 owner to zero address
  function test_CallableByL1OwnerOrZeroAddress() public {
    // Sets msg.sender and tx.origin
    vm.startPrank(s_strangerAddr, s_strangerAddr);

    // Defines the cross domain message to send
    vm.expectEmit(false, false, false, true);
    emit L1OwnershipTransferRequested(s_optimismCrossDomainForwarder.l1Owner(), address(0));

    // Sends the message
    s_mockOptimismCrossDomainMessenger.sendMessage(
      address(s_optimismCrossDomainForwarder), // target
      abi.encodeWithSelector(s_optimismCrossDomainForwarder.transferL1Ownership.selector, address(0)), // message
      0 // gas limit
    );

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }
}

contract AcceptL1Ownership is OptimismCrossDomainForwarderTest {
  /// @notice it should not be callable by non pending-owners
  function test_NotCallableByNonPendingOwners() public {
    // Sets msg.sender and tx.origin
    vm.startPrank(s_strangerAddr, s_strangerAddr);

    // Sends the message
    vm.expectRevert("Must be proposed L1 owner");
    s_mockOptimismCrossDomainMessenger.sendMessage(
      address(s_optimismCrossDomainForwarder), // target
      abi.encodeWithSelector(s_optimismCrossDomainForwarder.acceptL1Ownership.selector), // message
      0 // gas limit
    );

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }

  /// @notice it should be callable by pending L1 owner
  function test_CallableByPendingL1Owner() public {
    // Sets msg.sender and tx.origin
    vm.startPrank(s_strangerAddr, s_strangerAddr);

    // Request ownership transfer
    s_mockOptimismCrossDomainMessenger.sendMessage(
      address(s_optimismCrossDomainForwarder), // target
      abi.encodeWithSelector(s_optimismCrossDomainForwarder.transferL1Ownership.selector, s_strangerAddr), // message
      0 // gas limit
    );

    // Sets a mock message sender
    s_mockOptimismCrossDomainMessenger._setMockMessageSender(s_strangerAddr);

    // Prepares expected event payload
    vm.expectEmit(false, false, false, true);
    emit L1OwnershipTransferred(s_l1OwnerAddr, s_strangerAddr);

    // Accepts ownership transfer request
    s_mockOptimismCrossDomainMessenger.sendMessage(
      address(s_optimismCrossDomainForwarder), // target
      abi.encodeWithSelector(s_optimismCrossDomainForwarder.acceptL1Ownership.selector, s_strangerAddr), // message
      0 // gas limit
    );

    // Asserts that the ownership was actually transferred
    assertEq(s_optimismCrossDomainForwarder.l1Owner(), s_strangerAddr);

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }
}