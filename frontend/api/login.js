import httpRequest from "~/utils/request";

export const loginApi = {
  async login(params = {}, server = true) {
    let options = {
      url: `/api/login`,
      method: "post",
      params,
      server,
    };
    let result = await httpRequest.post(options);
    return result;
  },
  async logout(params = {}, server = true) {
    let options = {
      url: `/api/logout`,
      method: "post",
      params,
      server,
    };
    let result = await httpRequest.post(options);
    return result;
  },
  async register(params = {}, server = true) {
    let options = {
      url: `/api/register`,
      method: "post",
      params,
      server,
    };
    let result = await httpRequest.post(options);
    return result;
  },
};
