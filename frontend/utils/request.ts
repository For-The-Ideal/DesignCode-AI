import { P } from "vue-router/dist/index-BzEKChPW.js";

function request<T>(url: string, method: string, params?: T) {
    let option = {
        url,
        method,
    }
  // return fetch(url, {
  //     method,
  //     body: body ? JSON.stringify(body) : undefined,
  // })
}

export const httpPequest = {
  get(url: string) {
    return request(url, "GET");
  },

  post<T>(url: string, params: T) {
    return request(url, "POST", params);
  },

  delete(url: string) {
   return request(url, "DELETE");
  },
};
