// package: omeapi.worker.v1beta1
// file: omeapi/worker/v1beta1/worker.proto

/* tslint:disable */
/* eslint-disable */

import * as jspb from "google-protobuf";

export class GetWorkerRequest extends jspb.Message { 
    getWorkerId(): string;
    setWorkerId(value: string): GetWorkerRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GetWorkerRequest.AsObject;
    static toObject(includeInstance: boolean, msg: GetWorkerRequest): GetWorkerRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: GetWorkerRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GetWorkerRequest;
    static deserializeBinaryFromReader(message: GetWorkerRequest, reader: jspb.BinaryReader): GetWorkerRequest;
}

export namespace GetWorkerRequest {
    export type AsObject = {
        workerId: string,
    }
}

export class GetWorkerResponse extends jspb.Message { 

    hasWorker(): boolean;
    clearWorker(): void;
    getWorker(): Worker | undefined;
    setWorker(value?: Worker): GetWorkerResponse;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GetWorkerResponse.AsObject;
    static toObject(includeInstance: boolean, msg: GetWorkerResponse): GetWorkerResponse.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: GetWorkerResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GetWorkerResponse;
    static deserializeBinaryFromReader(message: GetWorkerResponse, reader: jspb.BinaryReader): GetWorkerResponse;
}

export namespace GetWorkerResponse {
    export type AsObject = {
        worker?: Worker.AsObject,
    }
}

export class SetStateRequest extends jspb.Message { 
    getWorkerId(): string;
    setWorkerId(value: string): SetStateRequest;
    getState(): WorkerState;
    setState(value: WorkerState): SetStateRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SetStateRequest.AsObject;
    static toObject(includeInstance: boolean, msg: SetStateRequest): SetStateRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: SetStateRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SetStateRequest;
    static deserializeBinaryFromReader(message: SetStateRequest, reader: jspb.BinaryReader): SetStateRequest;
}

export namespace SetStateRequest {
    export type AsObject = {
        workerId: string,
        state: WorkerState,
    }
}

export class SetStateResponse extends jspb.Message { 

    hasWorker(): boolean;
    clearWorker(): void;
    getWorker(): Worker | undefined;
    setWorker(value?: Worker): SetStateResponse;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SetStateResponse.AsObject;
    static toObject(includeInstance: boolean, msg: SetStateResponse): SetStateResponse.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: SetStateResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SetStateResponse;
    static deserializeBinaryFromReader(message: SetStateResponse, reader: jspb.BinaryReader): SetStateResponse;
}

export namespace SetStateResponse {
    export type AsObject = {
        worker?: Worker.AsObject,
    }
}

export class QueryByStateRequest extends jspb.Message { 
    getState(): WorkerState;
    setState(value: WorkerState): QueryByStateRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): QueryByStateRequest.AsObject;
    static toObject(includeInstance: boolean, msg: QueryByStateRequest): QueryByStateRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: QueryByStateRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): QueryByStateRequest;
    static deserializeBinaryFromReader(message: QueryByStateRequest, reader: jspb.BinaryReader): QueryByStateRequest;
}

export namespace QueryByStateRequest {
    export type AsObject = {
        state: WorkerState,
    }
}

export class Worker extends jspb.Message { 
    getWorkerId(): string;
    setWorkerId(value: string): Worker;
    getState(): WorkerState;
    setState(value: WorkerState): Worker;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Worker.AsObject;
    static toObject(includeInstance: boolean, msg: Worker): Worker.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Worker, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Worker;
    static deserializeBinaryFromReader(message: Worker, reader: jspb.BinaryReader): Worker;
}

export namespace Worker {
    export type AsObject = {
        workerId: string,
        state: WorkerState,
    }
}

export class QueryByStateResponse extends jspb.Message { 
    clearWorkersList(): void;
    getWorkersList(): Array<Worker>;
    setWorkersList(value: Array<Worker>): QueryByStateResponse;
    addWorkers(value?: Worker, index?: number): Worker;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): QueryByStateResponse.AsObject;
    static toObject(includeInstance: boolean, msg: QueryByStateResponse): QueryByStateResponse.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: QueryByStateResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): QueryByStateResponse;
    static deserializeBinaryFromReader(message: QueryByStateResponse, reader: jspb.BinaryReader): QueryByStateResponse;
}

export namespace QueryByStateResponse {
    export type AsObject = {
        workersList: Array<Worker.AsObject>,
    }
}

export enum WorkerState {
    WORKER_STATE_UNSPECIFIED = 0,
    WORKER_STATE_ONLINE = 1,
    WORKER_STATE_OFFLINE = 2,
    WORKER_STATE_ON_JOB = 3,
}
