import httpRequest from "~/utils/request";

export const configApi = {

  // 获取生成配置字典（框架/平台/选项/高级/组件库）
  async getGenConfig(server = true) {
    let options = {
      url: `/api/v1/config/options`,
      method: "get",
      server,
    };
    let result = await httpRequest.post(options);
    return result;
  },
};
