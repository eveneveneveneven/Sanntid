# Sanntid

This is a project for running multiple elevators on the same network, including fault tolerance and an approach to optimal service.

## Functionality

This elevatornetwork can:

- include any number of elevators in one network
- no orders lost, assuming one elevator is always on the network
- different levels of fault tolerance, which can handle crashes and undefined states
- soft resolving of internal orders in case of a killed program
- fully functional backup system

## Design

The way this network is designed is a master/slave approach. There will always be a Master on the network, which will be the sole truth of what the network status is. The Master also has the responsibility to update all slaves with the current network status.

In case of a missing Master, the slave with the id=1 will take over as a Master, and resume maintaining the network status. Since every slave has a copy of the network status no orders will be lost. 

In case of multiple Masters on the network, the Master with no connections prior to the dispute will resolve the current orders on the local elevator, and reset the program. This will usually only happen if an elevator has lost internet connection (by pulling the cord etc.), and is reconnected. 