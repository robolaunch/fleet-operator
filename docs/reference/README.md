# API Reference


## fleet.roboscale.io/v1alpha1

Package v1alpha1 contains API Schema definitions for the fleet v1alpha1 API group.

### Resource Types
- [Fleet](#fleet)



#### Fleet



Fleet manages lifecycle and configuration of multiple robots and robot's connectivity layer that contains DDS Discovery Server and ROS bridge services.



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `fleet.roboscale.io/v1alpha1`
| `kind` _string_ | `Fleet`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[FleetSpec](#fleetspec)_ | Specification of the desired behavior of the Fleet. |
| `status` _[FleetStatus](#fleetstatus)_ | Most recently observed status of the Fleet. |


#### FleetSpec





_Appears in:_
- [Fleet](#fleet)

| Field | Description |
| --- | --- |
| `discoveryServerTemplate` _[DiscoveryServerSpec](#discoveryserverspec)_ | Discovery server configuration of fleet. For detailed information, refer the document for the API group `robot.roboscale.io`. |
| `hybrid` _boolean_ | Determines if the fleet should be federated across clusters or not. |
| `instances` _string array_ | If `.spec.hybrid` is true, this field includes Kubernetes cluster names which the fleet will be federated to. |


#### FleetStatus





_Appears in:_
- [Fleet](#fleet)

| Field | Description |
| --- | --- |
| `phase` _FleetPhase_ | Fleet phase. |
| `namespaceStatus` _[OwnedNamespaceStatus](#ownednamespacestatus)_ | Namespace status. Fleet creates namespace if the `.spec.hybrid` is set to `true`. It creates `FederatedNamespace` if `false`. |
| `discoveryServerStatus` _[OwnedResourceStatus](#ownedresourcestatus)_ | Discovery server instance status. For detailed information, refer the document for the API group `robot.roboscale.io`. |
| `attachedRobots` _[AttachedRobot](#attachedrobot) array_ | Attached launch object information. |


#### AttachedRobot





_Appears in:_
- [FleetStatus](#fleetstatus)

| Field | Description |
| --- | --- |
| `reference` _[ObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#objectreference-v1-core)_ | Resource reference for attached robot. |
| `phase` _RobotPhase_ | Attached robot phase. For detailed information, refer the document for the API group `robot.roboscale.io`. |
| `fleetCompatibility` _[FleetCompatibilityStatus](#fleetcompatibilitystatus)_ | Compatibility status of attached robot with the fleet. |


#### FleetCompatibilityStatus





_Appears in:_
- [AttachedRobot](#attachedrobot)

| Field | Description |
| --- | --- |
| `isCompatible` _boolean_ | Indicates the robot's compatibility with fleet. |
| `reason` _string_ | Indicates the possible incompatibility reason of an attached robot. |


#### OwnedNamespaceStatus





_Appears in:_
- [FleetStatus](#fleetstatus)

| Field | Description |
| --- | --- |
| `resource` _[OwnedResourceStatus](#ownedresourcestatus)_ | Generic structure of the most recent status of an owned object. For detailed information, refer the document for the API group `robot.roboscale.io`. |
| `federated` _boolean_ | Sets to `true` if the owned namespace is federated. |
| `ready` _boolean_ | Sets to `true` if the namespace is ready for the resources to be deployed such as robot. |


