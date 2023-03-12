/* eslint-disable */
import {
  makeGenericClientConstructor,
  ChannelCredentials,
  ChannelOptions,
  UntypedServiceImplementation,
  handleUnaryCall,
  Client,
  ClientUnaryCall,
  Metadata,
  CallOptions,
  ServiceError,
} from "@grpc/grpc-js";
import Long from "long";
import _m0 from "protobufjs/minimal";

export const protobufPackage = "transaction_temporal_workflow.api";

export interface CreateTransactionRequest {
  transactionId: string;
  amount: number;
  productCode: string;
  userId: string;
  idempotencyKey: string;
}

export interface CreateTransactionResponse {
  message: string;
}

export interface ProcessTransactionRequest {
  transactionId: string;
  idempotencyKey: string;
}

export interface ProcessTransactionResponse {
  message: string;
}

function createBaseCreateTransactionRequest(): CreateTransactionRequest {
  return {
    transactionId: "",
    amount: 0,
    productCode: "",
    userId: "",
    idempotencyKey: "",
  };
}

export const CreateTransactionRequest = {
  encode(
    message: CreateTransactionRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.transactionId !== "") {
      writer.uint32(10).string(message.transactionId);
    }
    if (message.amount !== 0) {
      writer.uint32(16).int64(message.amount);
    }
    if (message.productCode !== "") {
      writer.uint32(26).string(message.productCode);
    }
    if (message.userId !== "") {
      writer.uint32(34).string(message.userId);
    }
    if (message.idempotencyKey !== "") {
      writer.uint32(42).string(message.idempotencyKey);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): CreateTransactionRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateTransactionRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.transactionId = reader.string();
          break;
        case 2:
          message.amount = longToNumber(reader.int64() as Long);
          break;
        case 3:
          message.productCode = reader.string();
          break;
        case 4:
          message.userId = reader.string();
          break;
        case 5:
          message.idempotencyKey = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CreateTransactionRequest {
    return {
      transactionId: isSet(object.transactionId)
        ? String(object.transactionId)
        : "",
      amount: isSet(object.amount) ? Number(object.amount) : 0,
      productCode: isSet(object.productCode) ? String(object.productCode) : "",
      userId: isSet(object.userId) ? String(object.userId) : "",
      idempotencyKey: isSet(object.idempotencyKey)
        ? String(object.idempotencyKey)
        : "",
    };
  },

  toJSON(message: CreateTransactionRequest): unknown {
    const obj: any = {};
    message.transactionId !== undefined &&
      (obj.transactionId = message.transactionId);
    message.amount !== undefined && (obj.amount = Math.round(message.amount));
    message.productCode !== undefined &&
      (obj.productCode = message.productCode);
    message.userId !== undefined && (obj.userId = message.userId);
    message.idempotencyKey !== undefined &&
      (obj.idempotencyKey = message.idempotencyKey);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CreateTransactionRequest>, I>>(
    object: I
  ): CreateTransactionRequest {
    const message = createBaseCreateTransactionRequest();
    message.transactionId = object.transactionId ?? "";
    message.amount = object.amount ?? 0;
    message.productCode = object.productCode ?? "";
    message.userId = object.userId ?? "";
    message.idempotencyKey = object.idempotencyKey ?? "";
    return message;
  },
};

function createBaseCreateTransactionResponse(): CreateTransactionResponse {
  return { message: "" };
}

export const CreateTransactionResponse = {
  encode(
    message: CreateTransactionResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.message !== "") {
      writer.uint32(10).string(message.message);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): CreateTransactionResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateTransactionResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.message = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CreateTransactionResponse {
    return {
      message: isSet(object.message) ? String(object.message) : "",
    };
  },

  toJSON(message: CreateTransactionResponse): unknown {
    const obj: any = {};
    message.message !== undefined && (obj.message = message.message);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CreateTransactionResponse>, I>>(
    object: I
  ): CreateTransactionResponse {
    const message = createBaseCreateTransactionResponse();
    message.message = object.message ?? "";
    return message;
  },
};

function createBaseProcessTransactionRequest(): ProcessTransactionRequest {
  return { transactionId: "", idempotencyKey: "" };
}

export const ProcessTransactionRequest = {
  encode(
    message: ProcessTransactionRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.transactionId !== "") {
      writer.uint32(10).string(message.transactionId);
    }
    if (message.idempotencyKey !== "") {
      writer.uint32(18).string(message.idempotencyKey);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): ProcessTransactionRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseProcessTransactionRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.transactionId = reader.string();
          break;
        case 2:
          message.idempotencyKey = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ProcessTransactionRequest {
    return {
      transactionId: isSet(object.transactionId)
        ? String(object.transactionId)
        : "",
      idempotencyKey: isSet(object.idempotencyKey)
        ? String(object.idempotencyKey)
        : "",
    };
  },

  toJSON(message: ProcessTransactionRequest): unknown {
    const obj: any = {};
    message.transactionId !== undefined &&
      (obj.transactionId = message.transactionId);
    message.idempotencyKey !== undefined &&
      (obj.idempotencyKey = message.idempotencyKey);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ProcessTransactionRequest>, I>>(
    object: I
  ): ProcessTransactionRequest {
    const message = createBaseProcessTransactionRequest();
    message.transactionId = object.transactionId ?? "";
    message.idempotencyKey = object.idempotencyKey ?? "";
    return message;
  },
};

function createBaseProcessTransactionResponse(): ProcessTransactionResponse {
  return { message: "" };
}

export const ProcessTransactionResponse = {
  encode(
    message: ProcessTransactionResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.message !== "") {
      writer.uint32(10).string(message.message);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): ProcessTransactionResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseProcessTransactionResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.message = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ProcessTransactionResponse {
    return {
      message: isSet(object.message) ? String(object.message) : "",
    };
  },

  toJSON(message: ProcessTransactionResponse): unknown {
    const obj: any = {};
    message.message !== undefined && (obj.message = message.message);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ProcessTransactionResponse>, I>>(
    object: I
  ): ProcessTransactionResponse {
    const message = createBaseProcessTransactionResponse();
    message.message = object.message ?? "";
    return message;
  },
};

export type TransactionService = typeof TransactionService;
export const TransactionService = {
  createTransaction: {
    path: "/transaction_temporal_workflow.api.Transaction/CreateTransaction",
    requestStream: false,
    responseStream: false,
    requestSerialize: (value: CreateTransactionRequest) =>
      Buffer.from(CreateTransactionRequest.encode(value).finish()),
    requestDeserialize: (value: Buffer) =>
      CreateTransactionRequest.decode(value),
    responseSerialize: (value: CreateTransactionResponse) =>
      Buffer.from(CreateTransactionResponse.encode(value).finish()),
    responseDeserialize: (value: Buffer) =>
      CreateTransactionResponse.decode(value),
  },
  processTransaction: {
    path: "/transaction_temporal_workflow.api.Transaction/ProcessTransaction",
    requestStream: false,
    responseStream: false,
    requestSerialize: (value: ProcessTransactionRequest) =>
      Buffer.from(ProcessTransactionRequest.encode(value).finish()),
    requestDeserialize: (value: Buffer) =>
      ProcessTransactionRequest.decode(value),
    responseSerialize: (value: ProcessTransactionResponse) =>
      Buffer.from(ProcessTransactionResponse.encode(value).finish()),
    responseDeserialize: (value: Buffer) =>
      ProcessTransactionResponse.decode(value),
  },
} as const;

export interface TransactionServer extends UntypedServiceImplementation {
  createTransaction: handleUnaryCall<
    CreateTransactionRequest,
    CreateTransactionResponse
  >;
  processTransaction: handleUnaryCall<
    ProcessTransactionRequest,
    ProcessTransactionResponse
  >;
}

export interface TransactionClient extends Client {
  createTransaction(
    request: CreateTransactionRequest,
    callback: (
      error: ServiceError | null,
      response: CreateTransactionResponse
    ) => void
  ): ClientUnaryCall;
  createTransaction(
    request: CreateTransactionRequest,
    metadata: Metadata,
    callback: (
      error: ServiceError | null,
      response: CreateTransactionResponse
    ) => void
  ): ClientUnaryCall;
  createTransaction(
    request: CreateTransactionRequest,
    metadata: Metadata,
    options: Partial<CallOptions>,
    callback: (
      error: ServiceError | null,
      response: CreateTransactionResponse
    ) => void
  ): ClientUnaryCall;
  processTransaction(
    request: ProcessTransactionRequest,
    callback: (
      error: ServiceError | null,
      response: ProcessTransactionResponse
    ) => void
  ): ClientUnaryCall;
  processTransaction(
    request: ProcessTransactionRequest,
    metadata: Metadata,
    callback: (
      error: ServiceError | null,
      response: ProcessTransactionResponse
    ) => void
  ): ClientUnaryCall;
  processTransaction(
    request: ProcessTransactionRequest,
    metadata: Metadata,
    options: Partial<CallOptions>,
    callback: (
      error: ServiceError | null,
      response: ProcessTransactionResponse
    ) => void
  ): ClientUnaryCall;
}

export const TransactionClient = makeGenericClientConstructor(
  TransactionService,
  "transaction_temporal_workflow.api.Transaction"
) as unknown as {
  new (
    address: string,
    credentials: ChannelCredentials,
    options?: Partial<ChannelOptions>
  ): TransactionClient;
  service: typeof TransactionService;
};

declare var self: any | undefined;
declare var window: any | undefined;
declare var global: any | undefined;
var globalThis: any = (() => {
  if (typeof globalThis !== "undefined") return globalThis;
  if (typeof self !== "undefined") return self;
  if (typeof window !== "undefined") return window;
  if (typeof global !== "undefined") return global;
  throw "Unable to locate global object";
})();

type Builtin =
  | Date
  | Function
  | Uint8Array
  | string
  | number
  | boolean
  | undefined;

export type DeepPartial<T> = T extends Builtin
  ? T
  : T extends Array<infer U>
  ? Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U>
  ? ReadonlyArray<DeepPartial<U>>
  : T extends {}
  ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin
  ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & {
      [K in Exclude<keyof I, KeysOfUnion<P>>]: never;
    };

function longToNumber(long: Long): number {
  if (long.gt(Number.MAX_SAFE_INTEGER)) {
    throw new globalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
  }
  return long.toNumber();
}

if (_m0.util.Long !== Long) {
  _m0.util.Long = Long as any;
  _m0.configure();
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
