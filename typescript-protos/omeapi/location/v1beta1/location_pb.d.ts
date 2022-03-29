// package: omeapi.location.v1beta1
// file: omeapi/location/v1beta1/location.proto

/* tslint:disable */
/* eslint-disable */

import * as jspb from "google-protobuf";
import * as google_protobuf_timestamp_pb from "google-protobuf/google/protobuf/timestamp_pb";

export class UpdateLocationRequest extends jspb.Message { 
    getWorkerId(): string;
    setWorkerId(value: string): UpdateLocationRequest;
    getLongitude(): number;
    setLongitude(value: number): UpdateLocationRequest;
    getLatitude(): number;
    setLatitude(value: number): UpdateLocationRequest;

    hasUpdateTime(): boolean;
    clearUpdateTime(): void;
    getUpdateTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setUpdateTime(value?: google_protobuf_timestamp_pb.Timestamp): UpdateLocationRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): UpdateLocationRequest.AsObject;
    static toObject(includeInstance: boolean, msg: UpdateLocationRequest): UpdateLocationRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: UpdateLocationRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): UpdateLocationRequest;
    static deserializeBinaryFromReader(message: UpdateLocationRequest, reader: jspb.BinaryReader): UpdateLocationRequest;
}

export namespace UpdateLocationRequest {
    export type AsObject = {
        workerId: string,
        longitude: number,
        latitude: number,
        updateTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    }
}

export class TollgateCrossing extends jspb.Message { 
    getTollgateId(): string;
    setTollgateId(value: string): TollgateCrossing;
    getLongitude(): number;
    setLongitude(value: number): TollgateCrossing;
    getLatitude(): number;
    setLatitude(value: number): TollgateCrossing;
    getDirection(): string;
    setDirection(value: string): TollgateCrossing;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): TollgateCrossing.AsObject;
    static toObject(includeInstance: boolean, msg: TollgateCrossing): TollgateCrossing.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: TollgateCrossing, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): TollgateCrossing;
    static deserializeBinaryFromReader(message: TollgateCrossing, reader: jspb.BinaryReader): TollgateCrossing;
}

export namespace TollgateCrossing {
    export type AsObject = {
        tollgateId: string,
        longitude: number,
        latitude: number,
        direction: string,
    }
}

export class UpdateLocationResponse extends jspb.Message { 
    getWorkerId(): string;
    setWorkerId(value: string): UpdateLocationResponse;

    hasTollgateCrossing(): boolean;
    clearTollgateCrossing(): void;
    getTollgateCrossing(): TollgateCrossing | undefined;
    setTollgateCrossing(value?: TollgateCrossing): UpdateLocationResponse;

    hasUpdateTime(): boolean;
    clearUpdateTime(): void;
    getUpdateTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setUpdateTime(value?: google_protobuf_timestamp_pb.Timestamp): UpdateLocationResponse;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): UpdateLocationResponse.AsObject;
    static toObject(includeInstance: boolean, msg: UpdateLocationResponse): UpdateLocationResponse.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: UpdateLocationResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): UpdateLocationResponse;
    static deserializeBinaryFromReader(message: UpdateLocationResponse, reader: jspb.BinaryReader): UpdateLocationResponse;
}

export namespace UpdateLocationResponse {
    export type AsObject = {
        workerId: string,
        tollgateCrossing?: TollgateCrossing.AsObject,
        updateTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    }
}

export class QueryLocationRequest extends jspb.Message { 
    getWorkerId(): string;
    setWorkerId(value: string): QueryLocationRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): QueryLocationRequest.AsObject;
    static toObject(includeInstance: boolean, msg: QueryLocationRequest): QueryLocationRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: QueryLocationRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): QueryLocationRequest;
    static deserializeBinaryFromReader(message: QueryLocationRequest, reader: jspb.BinaryReader): QueryLocationRequest;
}

export namespace QueryLocationRequest {
    export type AsObject = {
        workerId: string,
    }
}

export class QueryLocationResponse extends jspb.Message { 
    getWorkerId(): string;
    setWorkerId(value: string): QueryLocationResponse;
    getLongitude(): number;
    setLongitude(value: number): QueryLocationResponse;
    getLatitude(): number;
    setLatitude(value: number): QueryLocationResponse;

    hasLastSeenTime(): boolean;
    clearLastSeenTime(): void;
    getLastSeenTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setLastSeenTime(value?: google_protobuf_timestamp_pb.Timestamp): QueryLocationResponse;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): QueryLocationResponse.AsObject;
    static toObject(includeInstance: boolean, msg: QueryLocationResponse): QueryLocationResponse.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: QueryLocationResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): QueryLocationResponse;
    static deserializeBinaryFromReader(message: QueryLocationResponse, reader: jspb.BinaryReader): QueryLocationResponse;
}

export namespace QueryLocationResponse {
    export type AsObject = {
        workerId: string,
        longitude: number,
        latitude: number,
        lastSeenTime?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    }
}
