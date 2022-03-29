// package: omeapi.worker.v1beta1
// file: omeapi/worker/v1beta1/worker.proto

/* tslint:disable */
/* eslint-disable */

import * as grpc from "@grpc/grpc-js";
import * as omeapi_worker_v1beta1_worker_pb from "../../../omeapi/worker/v1beta1/worker_pb";

interface IWorkerServiceService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
    getWorker: IWorkerServiceService_IGetWorker;
    setState: IWorkerServiceService_ISetState;
    queryByState: IWorkerServiceService_IQueryByState;
}

interface IWorkerServiceService_IGetWorker extends grpc.MethodDefinition<omeapi_worker_v1beta1_worker_pb.GetWorkerRequest, omeapi_worker_v1beta1_worker_pb.GetWorkerResponse> {
    path: "/omeapi.worker.v1beta1.WorkerService/GetWorker";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<omeapi_worker_v1beta1_worker_pb.GetWorkerRequest>;
    requestDeserialize: grpc.deserialize<omeapi_worker_v1beta1_worker_pb.GetWorkerRequest>;
    responseSerialize: grpc.serialize<omeapi_worker_v1beta1_worker_pb.GetWorkerResponse>;
    responseDeserialize: grpc.deserialize<omeapi_worker_v1beta1_worker_pb.GetWorkerResponse>;
}
interface IWorkerServiceService_ISetState extends grpc.MethodDefinition<omeapi_worker_v1beta1_worker_pb.SetStateRequest, omeapi_worker_v1beta1_worker_pb.SetStateResponse> {
    path: "/omeapi.worker.v1beta1.WorkerService/SetState";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<omeapi_worker_v1beta1_worker_pb.SetStateRequest>;
    requestDeserialize: grpc.deserialize<omeapi_worker_v1beta1_worker_pb.SetStateRequest>;
    responseSerialize: grpc.serialize<omeapi_worker_v1beta1_worker_pb.SetStateResponse>;
    responseDeserialize: grpc.deserialize<omeapi_worker_v1beta1_worker_pb.SetStateResponse>;
}
interface IWorkerServiceService_IQueryByState extends grpc.MethodDefinition<omeapi_worker_v1beta1_worker_pb.QueryByStateRequest, omeapi_worker_v1beta1_worker_pb.QueryByStateResponse> {
    path: "/omeapi.worker.v1beta1.WorkerService/QueryByState";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<omeapi_worker_v1beta1_worker_pb.QueryByStateRequest>;
    requestDeserialize: grpc.deserialize<omeapi_worker_v1beta1_worker_pb.QueryByStateRequest>;
    responseSerialize: grpc.serialize<omeapi_worker_v1beta1_worker_pb.QueryByStateResponse>;
    responseDeserialize: grpc.deserialize<omeapi_worker_v1beta1_worker_pb.QueryByStateResponse>;
}

export const WorkerServiceService: IWorkerServiceService;

export interface IWorkerServiceServer extends grpc.UntypedServiceImplementation {
    getWorker: grpc.handleUnaryCall<omeapi_worker_v1beta1_worker_pb.GetWorkerRequest, omeapi_worker_v1beta1_worker_pb.GetWorkerResponse>;
    setState: grpc.handleUnaryCall<omeapi_worker_v1beta1_worker_pb.SetStateRequest, omeapi_worker_v1beta1_worker_pb.SetStateResponse>;
    queryByState: grpc.handleUnaryCall<omeapi_worker_v1beta1_worker_pb.QueryByStateRequest, omeapi_worker_v1beta1_worker_pb.QueryByStateResponse>;
}

export interface IWorkerServiceClient {
    getWorker(request: omeapi_worker_v1beta1_worker_pb.GetWorkerRequest, callback: (error: grpc.ServiceError | null, response: omeapi_worker_v1beta1_worker_pb.GetWorkerResponse) => void): grpc.ClientUnaryCall;
    getWorker(request: omeapi_worker_v1beta1_worker_pb.GetWorkerRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: omeapi_worker_v1beta1_worker_pb.GetWorkerResponse) => void): grpc.ClientUnaryCall;
    getWorker(request: omeapi_worker_v1beta1_worker_pb.GetWorkerRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: omeapi_worker_v1beta1_worker_pb.GetWorkerResponse) => void): grpc.ClientUnaryCall;
    setState(request: omeapi_worker_v1beta1_worker_pb.SetStateRequest, callback: (error: grpc.ServiceError | null, response: omeapi_worker_v1beta1_worker_pb.SetStateResponse) => void): grpc.ClientUnaryCall;
    setState(request: omeapi_worker_v1beta1_worker_pb.SetStateRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: omeapi_worker_v1beta1_worker_pb.SetStateResponse) => void): grpc.ClientUnaryCall;
    setState(request: omeapi_worker_v1beta1_worker_pb.SetStateRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: omeapi_worker_v1beta1_worker_pb.SetStateResponse) => void): grpc.ClientUnaryCall;
    queryByState(request: omeapi_worker_v1beta1_worker_pb.QueryByStateRequest, callback: (error: grpc.ServiceError | null, response: omeapi_worker_v1beta1_worker_pb.QueryByStateResponse) => void): grpc.ClientUnaryCall;
    queryByState(request: omeapi_worker_v1beta1_worker_pb.QueryByStateRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: omeapi_worker_v1beta1_worker_pb.QueryByStateResponse) => void): grpc.ClientUnaryCall;
    queryByState(request: omeapi_worker_v1beta1_worker_pb.QueryByStateRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: omeapi_worker_v1beta1_worker_pb.QueryByStateResponse) => void): grpc.ClientUnaryCall;
}

export class WorkerServiceClient extends grpc.Client implements IWorkerServiceClient {
    constructor(address: string, credentials: grpc.ChannelCredentials, options?: Partial<grpc.ClientOptions>);
    public getWorker(request: omeapi_worker_v1beta1_worker_pb.GetWorkerRequest, callback: (error: grpc.ServiceError | null, response: omeapi_worker_v1beta1_worker_pb.GetWorkerResponse) => void): grpc.ClientUnaryCall;
    public getWorker(request: omeapi_worker_v1beta1_worker_pb.GetWorkerRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: omeapi_worker_v1beta1_worker_pb.GetWorkerResponse) => void): grpc.ClientUnaryCall;
    public getWorker(request: omeapi_worker_v1beta1_worker_pb.GetWorkerRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: omeapi_worker_v1beta1_worker_pb.GetWorkerResponse) => void): grpc.ClientUnaryCall;
    public setState(request: omeapi_worker_v1beta1_worker_pb.SetStateRequest, callback: (error: grpc.ServiceError | null, response: omeapi_worker_v1beta1_worker_pb.SetStateResponse) => void): grpc.ClientUnaryCall;
    public setState(request: omeapi_worker_v1beta1_worker_pb.SetStateRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: omeapi_worker_v1beta1_worker_pb.SetStateResponse) => void): grpc.ClientUnaryCall;
    public setState(request: omeapi_worker_v1beta1_worker_pb.SetStateRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: omeapi_worker_v1beta1_worker_pb.SetStateResponse) => void): grpc.ClientUnaryCall;
    public queryByState(request: omeapi_worker_v1beta1_worker_pb.QueryByStateRequest, callback: (error: grpc.ServiceError | null, response: omeapi_worker_v1beta1_worker_pb.QueryByStateResponse) => void): grpc.ClientUnaryCall;
    public queryByState(request: omeapi_worker_v1beta1_worker_pb.QueryByStateRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: omeapi_worker_v1beta1_worker_pb.QueryByStateResponse) => void): grpc.ClientUnaryCall;
    public queryByState(request: omeapi_worker_v1beta1_worker_pb.QueryByStateRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: omeapi_worker_v1beta1_worker_pb.QueryByStateResponse) => void): grpc.ClientUnaryCall;
}
