import httpRequest from "~/utils/request";

export const commonApi = {
  // 获取模板
  async getTemplate(params = {}, server = true) {
    let options = {
      url: `/api/v1/template?id=${params.id}`,
      method: "get",
      params,
      server,
    };
    let result = await httpRequest.get(options);
    return result;
  },

  // 创建AI生成代码任务
  async generateUi(params = {}, server = true) {
    let options = {
      url: `/api/v1/generate-ui`,
      method: "post",
      params,
      server,
    };
    let result = await httpRequest.post(options);
    return result;
  },
};
