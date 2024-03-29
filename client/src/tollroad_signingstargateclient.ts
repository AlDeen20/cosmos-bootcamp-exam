import Long from "long"
import { GeneratedType, OfflineSigner, Registry } from "@cosmjs/proto-signing"
import {
    defaultRegistryTypes,
    DeliverTxResponse,
    QueryClient,
    SigningStargateClient,
    SigningStargateClientOptions,
    StdFee,
} from "@cosmjs/stargate"
import { Tendermint34Client } from "@cosmjs/tendermint-rpc"
import {
    setupTollroadExtension,
    TollroadExtension,
} from "./modules/tollroad/queries"
import {
    MsgCreateRoadOperatorEncodeObject,
    MsgCreateUserVaultEncodeObject,
    MsgDeleteRoadOperatorEncodeObject,
    MsgDeleteUserVaultEncodeObject,
    tollroadTypes,
    typeUrlMsgCreateRoadOperator,
    typeUrlMsgCreateUserVault,
    typeUrlMsgDeleteRoadOperator,
    typeUrlMsgDeleteUserVault,
} from "./types/tollroad/messages"

export const tollroadDefaultRegistryTypes: ReadonlyArray<
    [string, GeneratedType]
> = [...defaultRegistryTypes, ...tollroadTypes]

function createDefaultRegistry(): Registry {
    return new Registry(tollroadDefaultRegistryTypes)
}

export class TollroadSigningStargateClient extends SigningStargateClient {
    public readonly tollroadQueryClient: TollroadExtension | undefined

    public static async connectWithSigner(
        endpoint: string,
        signer: OfflineSigner,
        options: SigningStargateClientOptions = {},
    ): Promise<TollroadSigningStargateClient> {
        const tmClient = await Tendermint34Client.connect(endpoint)
        return new TollroadSigningStargateClient(tmClient, signer, {
            registry: createDefaultRegistry(),
            ...options,
        })
    }

    protected constructor(
        tmClient: Tendermint34Client | undefined,
        signer: OfflineSigner,
        options: SigningStargateClientOptions,
    ) {
        super(tmClient, signer, options)
        if (tmClient) {
            this.tollroadQueryClient = QueryClient.withExtensions(
                tmClient,
                setupTollroadExtension,
            )
        }
    }

    public async createRoadOperator(
        creator: string,
        name: string,
        token: string,
        active: boolean,
        fee: StdFee | "auto" | number,
        memo = "",
    ): Promise<DeliverTxResponse> {
        const createMsg: MsgCreateRoadOperatorEncodeObject = {
            typeUrl: typeUrlMsgCreateRoadOperator,
            value: {
                creator: creator,
                name: name,
                token: token,
                active: active,
            },
        }
        return this.signAndBroadcast(creator, [createMsg], fee, memo)
    }

    public async deleteRoadOperator(
        creator: string,
        index: string,
        fee: StdFee | "auto" | number,
        memo = "",
    ): Promise<DeliverTxResponse> {
        const deleteMsg: MsgDeleteRoadOperatorEncodeObject = {
            typeUrl: typeUrlMsgDeleteRoadOperator,
            value: {
                creator: creator,
                index: index,
            },
        }
        return this.signAndBroadcast(creator, [deleteMsg], fee, memo)
    }

    public async createUserVault(
        creator: string,
        roadOperatorIndex: string,
        token: string,
        balance: Long,
        fee: StdFee | "auto" | number,
        memo = "",
    ): Promise<DeliverTxResponse> {
        const createUserVaultMsg: MsgCreateUserVaultEncodeObject = {
            typeUrl: typeUrlMsgCreateUserVault,
            value: {
                creator: creator,
                roadOperatorIndex: roadOperatorIndex,
                token: token,
                balance: balance,
            },
        }
        return this.signAndBroadcast(creator, [createUserVaultMsg], fee, memo)
    }

    public async deleteUserVault(
        creator: string,
        roadOperatorIndex: string,
        token: string,
        fee: StdFee | "auto" | number,
        memo = "",
    ): Promise<DeliverTxResponse> {
        const deleteUserVaultMsg: MsgDeleteUserVaultEncodeObject = {
            typeUrl: typeUrlMsgDeleteUserVault,
            value: {
                creator: creator,
                roadOperatorIndex: roadOperatorIndex,
                token: token,
            },
        }
        return this.signAndBroadcast(creator, [deleteUserVaultMsg], fee, memo)
    }
}
