import httpRequest from "~/utils/request";

export const userApi = {
  async userInfo(params = {}, server = true) {
    let options = {
      url: `/api/v1/user/info`,
      method: "get",
      params,
      server,
    };
    let result = await httpRequest.get(options);
    return result;
  },
  
  async updateUserInfo(params = {}, server = true) {
    let options = {
      url: `/api/v1/user/update`,
      method: "post",
      params,
      server,
    };
    let result = await httpRequest.post(options);
    return result;
  },

  async updatePassword(params = {}, server = true) {
    let options = {
      url: `/api/v1/user/updatePassword`,
      method: "post",
      params,
      server,
    };
    let result = await httpRequest.post(options);
    return result;
  },
};
