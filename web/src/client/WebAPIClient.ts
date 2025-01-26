import { SourceModel } from "./models";
import RPCClient from "./RPCClient";

export class WebAPIClient {
  private client: RPCClient;

  public constructor() {
    this.client = new RPCClient();
  }

  public listSources(): Promise<SourceModel[]> {
    return this.client.execute<SourceModel[]>("ListSources");
  }

  public saveLabel(
    source: string,
    fromLabel: string,
    toLabel: string
  ): Promise<SourceModel> {
    return this.client.execute<SourceModel>("ReplaceLabel", [
      source,
      fromLabel,
      toLabel,
    ]);
  }

  public removeSource(
    source: string,
    label: string,
  ) {
    return this.client.execute<string>("RemoveSource", [
      source,
      label,
    ])
  }

}
