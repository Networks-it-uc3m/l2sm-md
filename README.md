# l2sm-md

L2S-M Multi Domain Client (L2S-M MD) is a gRPC and cli client that can be used for managing Multi Domain L2S-M resources, more specifically, manage Inter-Cluster Overlays, and Inter-Cluster L2Networks from one single point. 

If you want to learn more about L2S-M and it's core concepts, please refer to the [main repository](!https://www.github.com/Networks-it-uc3m/L2S-M). 

This tool is meant to be used from a Control Plane cluster or a docker image that has access to the other clusters. Additional setup must be done from the worker clusters to have a user with cluster role bindings that enable the managing of L2S-M Custom Resources (CRs), these are the Overlays, NetworkEdgeDevices and L2Networks. 

L2S-M MD then will create the necesary resources in the worker clusters and enable connectivity. The main use enables you to create the overlays and l2networks. In order to do this, you need an IDCO Provider (Inter Domain Connectivity Orchestrator) deployed and accesible from all clusters. Please refer to the architecture guide if you are interested in understanding the multi domain core concepts. 

## 1. Creation of Inter Cluster Overlays

In order to have an overlay interconnecting multiple clusters with L2S-M, Network Edge Devices must be deployed inside each managed cluster, in order to do so with L2S-M MD, you can do a gRPC request [(the proto file can be found here)](!./api/v1/l2smmd.proto):
```sh
message CreateOverlayRequest {
    OverlayTopology overlay = 1;
}
message OverlayTopology {
    Provider provider = 1;
    repeated Cluster clusters = 2;
    repeated Link links = 3;
}
message Provider {
    string name = 1;
    string domain = 2;
}
```
Specifying an OverlayTopology, which consists of a Provider, with the name of the IDCO Provider and Domain, accesible from every worker cluster, and a sequence of clusters and its links, which consist of the name of the clusters, with their configuration for managing them, and the links specifying which clusters connect to which. If links is left blank, a topology will be auto generated using the Net-Overlay-Builder.

## 2. 