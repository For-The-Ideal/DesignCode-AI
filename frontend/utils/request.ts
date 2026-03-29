import type { UseFetchOptions } from "#app";
type Methods = "GET" | "POST" | "DELETE" | "PUT";

const BASE_URL = "http://127.0.0.1:8090";
const requestUrl = "/api/ajaxData"

export interface IResultData<T> {
  code: number;
  data: T;
  message: string;
}

class HttpRequest {
  async request<T>(url: string, method: Methods, data: T) {
    return new Promise(async (resolve, reject) => {
      const newOptions: UseFetchOptions<T> = {
        baseURL: BASE_URL,
        headers: {
            "Content-Type": "application/json",
        },
        method: method,
        body: {
          url,
          // params:{
          //   ...data
          // }
          // ...data,
        } as any,
      };

      if (["get", "GET", "delete", "DELETE"].includes(method)) {
        newOptions.params = data as any;
      }
      if (["post", "POST", "put", "PUT"].includes(method)) {
        newOptions.body = data as any;
      }

      console.log(newOptions, "new");
      let reslut = await useFetch(requestUrl, newOptions)
        .then((res) => {
          resolve(res);
        })
        .catch((error) => {
          reject(error);
        });

      return reslut;
    });
  }

  // 封装常用方法

  get<T>(url: string, params?: T) {
    return this.request(url, "GET", params);
  }

  post<T>(url: string, params: T) {
    return this.request(url, "POST", params);
  }

  Put<T>(url: string, params: T) {
    return this.request(url, "PUT", params);
  }

  Delete<T>(url: string, params: T) {
    return this.request(url, "DELETE", params);
  }
}

const httpRequest = new HttpRequest();

export default httpRequest;
