import { app_config } from "./app_config";

class Fetcher {
  baseUrl: string;

  constructor(url: string) {
    this.baseUrl = url;
  }

  async GET<T>(path: string): Promise<T | string> {
    try {
      const resp = await fetch(app_config.BASE_URL + "/" + path, { method: "GET" });
      const data = (await resp.json()) as T;

      return data;
    } catch (error) {
      return new Promise((rep, _) => rep(`Get failed due ${error}`));
    }
  }
}

export const server = new Fetcher(app_config.BASE_URL);
