import httpRequest from "~/utils/request";

export const loginApi = {
  async captcha(server = true) {
    let options = {
      url: `/api/v1/auth/captcha`,
      method: "get",
      server,
    };
    let result = await httpRequest.get(options);
    return result;
  },
  async login(params = {}, server = true) {
    let options = {
      url: `/api/v1/auth/login`,
      method: "post",
      params,
      server,
    };
    let result = await httpRequest.post(options);
    return result;
  },
  async logout(params = {}, server = true) {
    let options = {
      url: `/api/v1/auth/logout`,
      method: "post",
      params,
      server,
    };
    let result = await httpRequest.post(options);
    return result;
  },
  async register(params = {}, server = true) {
    let options = {
      url: `/api/v1/auth/register`,
      method: "post",
      params,
      server,
    };
    let result = await httpRequest.post(options);
    return result;
  },
};
