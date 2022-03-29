// package: omeapi.location.v1beta1
// file: omeapi/location/v1beta1/location.proto

/* tslint:disable */
/* eslint-disable */

import * as grpc from "@grpc/grpc-js";
import * as omeapi_location_v1beta1_location_pb from "../../../omeapi/location/v1beta1/location_pb";
import * as google_protobuf_timestamp_pb from "google-protobuf/google/protobuf/timestamp_pb";

interface ILocationServiceService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
    updateLocation: ILocationServiceService_IUpdateLocation;
    queryLocation: ILocationServiceService_IQueryLocation;
}

interface ILocationServiceService_IUpdateLocation extends grpc.MethodDefinition<omeapi_location_v1beta1_location_pb.UpdateLocationRequest, omeapi_location_v1beta1_location_pb.UpdateLocationResponse> {
    path: "/omeapi.location.v1beta1.LocationService/UpdateLocation";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<omeapi_location_v1beta1_location_pb.UpdateLocationRequest>;
    requestDeserialize: grpc.deserialize<omeapi_location_v1beta1_location_pb.UpdateLocationRequest>;
    responseSerialize: grpc.serialize<omeapi_location_v1beta1_location_pb.UpdateLocationResponse>;
    responseDeserialize: grpc.deserialize<omeapi_location_v1beta1_location_pb.UpdateLocationResponse>;
}
interface ILocationServiceService_IQueryLocation extends grpc.MethodDefinition<omeapi_location_v1beta1_location_pb.QueryLocationRequest, omeapi_location_v1beta1_location_pb.QueryLocationResponse> {
    path: "/omeapi.location.v1beta1.LocationService/QueryLocation";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<omeapi_location_v1beta1_location_pb.QueryLocationRequest>;
    requestDeserialize: grpc.deserialize<omeapi_location_v1beta1_location_pb.QueryLocationRequest>;
    responseSerialize: grpc.serialize<omeapi_location_v1beta1_location_pb.QueryLocationResponse>;
    responseDeserialize: grpc.deserialize<omeapi_location_v1beta1_location_pb.QueryLocationResponse>;
}

export const LocationServiceService: ILocationServiceService;

export interface ILocationServiceServer extends grpc.UntypedServiceImplementation {
    updateLocation: grpc.handleUnaryCall<omeapi_location_v1beta1_location_pb.UpdateLocationRequest, omeapi_location_v1beta1_location_pb.UpdateLocationResponse>;
    queryLocation: grpc.handleUnaryCall<omeapi_location_v1beta1_location_pb.QueryLocationRequest, omeapi_location_v1beta1_location_pb.QueryLocationResponse>;
}

export interface ILocationServiceClient {
    updateLocation(request: omeapi_location_v1beta1_location_pb.UpdateLocationRequest, callback: (error: grpc.ServiceError | null, response: omeapi_location_v1beta1_location_pb.UpdateLocationResponse) => void): grpc.ClientUnaryCall;
    updateLocation(request: omeapi_location_v1beta1_location_pb.UpdateLocationRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: omeapi_location_v1beta1_location_pb.UpdateLocationResponse) => void): grpc.ClientUnaryCall;
    updateLocation(request: omeapi_location_v1beta1_location_pb.UpdateLocationRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: omeapi_location_v1beta1_location_pb.UpdateLocationResponse) => void): grpc.ClientUnaryCall;
    queryLocation(request: omeapi_location_v1beta1_location_pb.QueryLocationRequest, callback: (error: grpc.ServiceError | null, response: omeapi_location_v1beta1_location_pb.QueryLocationResponse) => void): grpc.ClientUnaryCall;
    queryLocation(request: omeapi_location_v1beta1_location_pb.QueryLocationRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: omeapi_location_v1beta1_location_pb.QueryLocationResponse) => void): grpc.ClientUnaryCall;
    queryLocation(request: omeapi_location_v1beta1_location_pb.QueryLocationRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: omeapi_location_v1beta1_location_pb.QueryLocationResponse) => void): grpc.ClientUnaryCall;
}

export class LocationServiceClient extends grpc.Client implements ILocationServiceClient {
    constructor(address: string, credentials: grpc.ChannelCredentials, options?: Partial<grpc.ClientOptions>);
    public updateLocation(request: omeapi_location_v1beta1_location_pb.UpdateLocationRequest, callback: (error: grpc.ServiceError | null, response: omeapi_location_v1beta1_location_pb.UpdateLocationResponse) => void): grpc.ClientUnaryCall;
    public updateLocation(request: omeapi_location_v1beta1_location_pb.UpdateLocationRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: omeapi_location_v1beta1_location_pb.UpdateLocationResponse) => void): grpc.ClientUnaryCall;
    public updateLocation(request: omeapi_location_v1beta1_location_pb.UpdateLocationRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: omeapi_location_v1beta1_location_pb.UpdateLocationResponse) => void): grpc.ClientUnaryCall;
    public queryLocation(request: omeapi_location_v1beta1_location_pb.QueryLocationRequest, callback: (error: grpc.ServiceError | null, response: omeapi_location_v1beta1_location_pb.QueryLocationResponse) => void): grpc.ClientUnaryCall;
    public queryLocation(request: omeapi_location_v1beta1_location_pb.QueryLocationRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: omeapi_location_v1beta1_location_pb.QueryLocationResponse) => void): grpc.ClientUnaryCall;
    public queryLocation(request: omeapi_location_v1beta1_location_pb.QueryLocationRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: omeapi_location_v1beta1_location_pb.QueryLocationResponse) => void): grpc.ClientUnaryCall;
}
