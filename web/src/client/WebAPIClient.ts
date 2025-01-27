import { SourceModel } from "./models";
import RPCClient from "./RPCClient";

export default class WebAPIClient {
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

  public connectSource(
    url: string,
    label: string,
  ) {
    return this.client.execute<string>("Connect", [
      url,
      label,
    ])
  }

  public updateSource(
    source: string,
    label: string,
  ) {
    return this.client.execute<string>("Update", [
      source,
      label,
    ])
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
