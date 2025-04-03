export interface JsonRPCError {
	code: number;
	message: string;
}

interface JsonRPCResponse<T> {
	jsonrpc: "2.0";
	id: string | number;
	result?: T;
	error?: JsonRPCError;
}

interface RpcCommand {
	jsonrpc: "2.0";
	id: string | number;
	method: string;
	params?: any[];
}

export default class RPCClient {
	// return codes for errors
	public static ErrRedirectToLogin = -32302;

	public async execute<T>(method: string, params?: any[]): Promise<T> {
		const id = Date.now();
		const body = {
			jsonrpc: "2.0",
			id,
			method,
			params,
		} as RpcCommand;

		const resp = await fetch("/rpc?" + method, {
			method: "POST",
			body: JSON.stringify(body),
			headers: {
				"Content-Type": "application/json",
				Accept: "application/json",
			},
		});
		if (!resp.ok) {
			throw Error(resp.statusText);
		}
		const jsonRpcResp: JsonRPCResponse<T> = await resp.json();
		if (jsonRpcResp.jsonrpc !== "2.0") {
			throw Error("Not a JSON-RPC response");
		}
		if (jsonRpcResp.id !== id) {
			throw Error("Not a valid JSON-RPC response");
		}
		if (jsonRpcResp.error) {
			return Promise.reject(jsonRpcResp.error);
		}
		return jsonRpcResp.result as T;
	}

}
